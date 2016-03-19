package irc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

func newImageAction() *imageAction {
	action := imageAction{
		database: IrcDatabase{},
		queue:    make([]*DBLink, 0),
	}

	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for _ = range ticker.C {
			if len(action.queue) > 0 {
				var l = &DBLink{}
				l, action.queue = action.queue[len(action.queue)-1], action.queue[:len(action.queue)-1]
				action.CheckLink(l)
			}
		}
	}()
	return &action
}

// ImageAction - Instance of imageAction struct
var ImageAction = newImageAction()

type imageAction struct {
	database IrcDatabase
	queue    []*DBLink
}

// Adds link to queue, non blocking
func (i *imageAction) AppendCheckLink(link *DBLink) {
	b, _ := contains(i.queue, link)
	if !b {
		log.Debugf("Append link to queue %d", link.Key.Int64)
		i.queue = append(i.queue, link)
	}
}

// Blocks until link is handled
func (i *imageAction) CheckLink(link *DBLink) {
	switch link.Link_type.String {
	case Image:
		if _, err := i.database.GetRaw(link.Key.Int64, RawImage); err != nil {
			i.Image(link)
		} else if _, err := i.database.GetRaw(link.Key.Int64, RawImageThumbnail); err != nil {
			i.Image(link)
		} else if notExists(link.Key.Int64, "image") {
			i.Image(link)
		} else if notExists(link.Key.Int64, "thumb") {
			i.Image(link)
		}

	case WebM:
		if _, err := i.database.GetRaw(link.Key.Int64, RawWebm); err != nil {
			i.WebM(link)
		} else if _, err := i.database.GetRaw(link.Key.Int64, RawWebmFrame); err != nil {
			i.WebM(link)
		}
	case Gif:
		if _, err := i.database.GetRaw(link.Key.Int64, RawGif); err != nil {
			i.Gif(link)
		} else if notExists(link.Key.Int64, "webm") {
			i.Gif(link)
		} else if notExists(link.Key.Int64, "thumb") {
			i.Gif(link)
		}
	}
	b, index := contains(i.queue, link)
	if b {
		log.Debugf("Removing duplicate link from queue %d", link.Key.Int64)
		i.queue = append(i.queue[:index], i.queue[index+1:]...)
	}
}

func (i *imageAction) Image(link *DBLink) {
	log.WithFields(log.Fields{
		"id":   link.Key.Int64,
		"link": link.Link.String,
	}).Debug("Downloading image")
	contentType, data, err := i.Download(link.Link.String)
	log.WithFields(log.Fields{
		"id":          link.Key.Int64,
		"ContentType": contentType,
		"DataLen":     len(data),
		"Err":         err,
	}).Debug("Download complete")
	if err != nil {
		i.database.RemoveLink(link)
		log.WithField("id", link.Key.Int64).Info("Removing dead link")
		return
	}
	if contentType == "image/png" || contentType == "image/jpeg" {
		if _, err := i.database.GetRaw(link.Key.Int64, RawImage); err != nil {
			log.WithField("id", link.Key.Int64).Debugf("Adding %s to raw ", RawImage)
			i.database.AddRaw(RawImage, link, contentType)
		}
		if _, err := i.database.GetRaw(link.Key.Int64, RawImageThumbnail); err != nil {
			log.WithField("id", link.Key.Int64).Debugf("Adding %s to raw ", RawImageThumbnail)
			i.database.AddRaw(RawImageThumbnail, link, contentType)
		}

		imgType := "jpeg"
		if contentType == "image/png" {
			imgType = "png"
		}

		filePath := filepath.Join(GetConfig().DataPath, "image", strconv.FormatInt(link.Key.Int64, 10))
		thumbPath := filepath.Join(GetConfig().DataPath, "thumb") + "/"
		thumbFilePath := filepath.Join(thumbPath, strconv.FormatInt(link.Key.Int64, 10)) + "." + imgType
		ioutil.WriteFile(filePath, data, os.ModePerm)
		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/c", fmt.Sprintf("vipsthumbnail -s 128 -c --interpolator nearest -o %s%%s.%s %s & rename %s *. ",
				thumbPath, imgType, filePath, thumbFilePath)).Output()
		} else {
			exec.Command("sh", "-c", fmt.Sprintf("vipsthumbnail -s 128 -c --interpolator nearest -o %s%%s.%s %s && mv %s %s`basename %s .%s`",
				thumbPath, imgType, filePath, thumbFilePath, thumbPath, thumbFilePath, imgType)).Output()
		}
		if notExists(link.Key.Int64, "thumb") {
			i.database.RemoveLink(link)
			os.Remove(filePath)
			log.WithField("id", link.Key.Int64).Errorf("Can't make thumbnail %s, vips error", link.Link.String)
		} else {
			log.WithField("id", link.Key.Int64).Debug("Image download/convert successful")
		}
	} else {
		log.Errorf("Can't make thumbnail, unknown image format %s %s", contentType, link.Link.String)
		i.database.RemoveLink(link)
	}
}

func (i *imageAction) Gif(link *DBLink) {
	_, data, err := i.Download(link.Link.String)
	if err != nil {
		i.database.RemoveLink(link)
		log.Infof("Removing dead link %s", link.Link.String)
		return
	}
	if _, err := i.database.GetRaw(link.Key.Int64, RawWebm); err != nil {
		i.database.AddRaw(RawWebm, link, "video/webm")
	}
	if _, err := i.database.GetRaw(link.Key.Int64, RawWebmFrame); err != nil {
		i.database.AddRaw(RawWebmFrame, link, "image/jpeg")
	}

	filePath := filepath.Join(GetConfig().DataPath, "/gif/", strconv.FormatInt(link.Key.Int64, 10)) + ".gif"
	outPath := filepath.Join(GetConfig().DataPath, "/webm/", strconv.FormatInt(link.Key.Int64, 10))
	framePath := filepath.Join(GetConfig().DataPath, "/thumb/", strconv.FormatInt(link.Key.Int64, 10))
	ioutil.WriteFile(filePath, data, os.ModePerm)
	gifToWebM(filePath, outPath, framePath)
	if notExists(link.Key.Int64, "thumb") {
		i.database.RemoveLink(link)
		os.Remove(filePath)
		os.Remove(outPath)
		fmt.Println(fmt.Sprintf("Cannot convert GIF -> WEBM %s", link.Link.String))
	}
}

func (i *imageAction) WebM(link *DBLink) {
	log.WithFields(log.Fields{
		"id":   link.Key.Int64,
		"link": link.Link.String,
	}).Debug("Downloading webm")
	contentType, data, err := i.Download(link.Link.String)
	log.WithFields(log.Fields{
		"id":          link.Key.Int64,
		"ContentType": contentType,
		"DataLen":     len(data),
		"Err":         err,
	}).Debug("Download complete")
	if err != nil {
		i.database.RemoveLink(link)
		log.WithField("id", link.Key.Int64).Info("Removing dead link")
		return
	}
	if _, err := i.database.GetRaw(link.Key.Int64, RawWebm); err != nil {
		i.database.AddRaw(RawWebm, link, "video/webm")
	}
	if _, err := i.database.GetRaw(link.Key.Int64, RawWebmFrame); err != nil {
		i.database.AddRaw(RawWebmFrame, link, "image/jpeg")
	}

	filePath := filepath.Join(GetConfig().DataPath, "/webm/", strconv.FormatInt(link.Key.Int64, 10))
	framePath := filepath.Join(GetConfig().DataPath, "/thumb/", strconv.FormatInt(link.Key.Int64, 10))
	ioutil.WriteFile(filePath, data, os.ModePerm)
	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/c", fmt.Sprintf("ffmpeg -y -i %s -f image2 -ss 00 -vframes 1 %s", filePath, framePath)).Output()
	} else {
		exec.Command("sh", "-c", fmt.Sprintf("ffmpeg -y -i %s -f image2 -ss 00 -vframes 1 %s", filePath, framePath)).Output()
	}
	if notExists(link.Key.Int64, "thumb") {
		i.database.RemoveLink(link)
		os.Remove(filePath)
		log.WithField("id", link.Key.Int64).Error("WebM convert error")
	}
}

func (i *imageAction) Download(url string) (string, []byte, error) {
	response, err := http.Get(url)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		fmt.Println(fmt.Sprintf("Can't download Image %s err %s", url, err))
		return "", []byte{}, err
	}
	if response.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("Image returned with wrong status code %d", response.StatusCode))
		return "", []byte{}, fmt.Errorf("StatusCode %d", response.StatusCode)
	}
	data, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	return response.Header.Get("Content-Type"), data, err
}

// Checks if there are unconverted gifs, missing thumbnails etc
func (i *imageAction) StartupCheck() {
	// Check unconverted gifs
	arr, _ := ioutil.ReadDir(filepath.Join(GetConfig().DataPath, "gif"))
	for _, t := range arr {
		baseName := strings.TrimSuffix(t.Name(), filepath.Ext(t.Name()))
		filePath := filepath.Join(GetConfig().DataPath, "/gif/", t.Name())
		outPath := filepath.Join(GetConfig().DataPath, "/webm/", baseName)
		framePath := filepath.Join(GetConfig().DataPath, "/thumb/", baseName)
		gifToWebM(filePath, outPath, framePath)
	}
	// TODO check missing thumbnails
}

func gifToWebM(filePath, outPath, framePath string) {
	exec.Command("sh", "-c", fmt.Sprintf("ffmpeg -y -i %s -c:v libvpx -crf 12 -b:v 500k -f webm %s && rm %s", filePath, outPath, filePath)).Output()
	exec.Command("sh", "-c", fmt.Sprintf("ffmpeg -y -i %s -f image2 -ss 00 -vframes 1 %s", outPath, framePath)).Output()
}

func notExists(id int64, fileType string) bool {
	if _, err := os.Stat(filepath.Join(GetConfig().DataPath, fileType, strconv.FormatInt(id, 10))); os.IsNotExist(err) {
		return true
	}
	return false
}

// Check element exists in list, if true returns elemnt index on list
func contains(list []*DBLink, element *DBLink) (bool, int) {
	for i, a := range list {
		if a.Key.Int64 == element.Key.Int64 {
			return true, i
		}
	}
	return false, 0
}
