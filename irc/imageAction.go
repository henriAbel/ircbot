package irc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

func newImageAction() *imageAction {
	cs1 := make(chan *DBLink, 1)
	cs2 := make(chan *DBLink, 1)
	cs3 := make(chan *DBLink, 1)

	action := imageAction{
		database: IrcDatabase{},
		imageCh:  cs1,
		gifCh:    cs2,
		webmCh:   cs3,
		queue:    make([]*DBLink, 0),
	}

	go action.Image(cs1)
	go action.Gif(cs2)
	go action.WebM(cs3)

	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for _ = range ticker.C {
			if len(action.queue) > 0 {
				var l = &DBLink{}
				l, action.queue = action.queue[len(action.queue)-1], action.queue[:len(action.queue)-1]
				action.checkImage(l)
			}
		}
	}()
	return &action
}

var ImageAction = newImageAction()

type imageAction struct {
	database IrcDatabase
	imageCh  chan *DBLink
	gifCh    chan *DBLink
	webmCh   chan *DBLink
	queue    []*DBLink
}

func (i *imageAction) CheckImage(link *DBLink) {
	if !contains(i.queue, link) {
		i.queue = append(i.queue, link)
	}
}

func (i *imageAction) checkImage(link *DBLink) {
	switch link.Link_type.String {
	case Image:
		if _, err := i.database.GetRaw(link.Key.Int64, RawImage); err != nil {
			i.imageCh <- link
		} else if _, err := i.database.GetRaw(link.Key.Int64, RawImageThumbnail); err != nil {
			i.imageCh <- link
		} else if notExists(link.Key.Int64, "image") {
			i.imageCh <- link
		} else if notExists(link.Key.Int64, "thumb") {
			i.imageCh <- link
		}

	case WebM:
		if _, err := i.database.GetRaw(link.Key.Int64, RawWebm); err != nil {
			i.webmCh <- link
		}
		if _, err := i.database.GetRaw(link.Key.Int64, RawWebmFrame); err != nil {
			i.webmCh <- link
		}
	}
}

func (i *imageAction) Image(cs chan *DBLink) {
	for link := range cs {
		contentType, data, err := i.Download(link.Link.String)
		if err != nil {
			i.database.RemoveLink(link)
			fmt.Println(fmt.Sprintf("Removing dead link %s", link.Link.String))
			continue
		}
		if contentType == "image/png" || contentType == "image/jpeg" {
			if _, err := i.database.GetRaw(link.Key.Int64, RawImage); err != nil {
				i.database.AddRaw(RawImage, link, contentType)
			}
			if _, err := i.database.GetRaw(link.Key.Int64, RawImageThumbnail); err != nil {
				i.database.AddRaw(RawImageThumbnail, link, contentType)
			}

			imgType := "jpeg"
			if contentType == "image/png" {
				imgType = "png"
			}
			filePath := path.Join(GetConfig().DataPath, "image", strconv.FormatInt(link.Key.Int64, 10))
			thumbPath := path.Join(GetConfig().DataPath, "thumb") + "/"
			thumbFilePath := path.Join(thumbPath, strconv.FormatInt(link.Key.Int64, 10)) + "." + imgType
			if notExists(link.Key.Int64, "image") {
				ioutil.WriteFile(filePath, data, os.ModePerm)
				exec.Command("sh", "-c", fmt.Sprintf("vipsthumbnail -s 128 -o %s%%s.%s %s && mv %s %s`basename %s .%s`",
					thumbPath, imgType, filePath, thumbFilePath, thumbPath, thumbFilePath, imgType)).Output()
				if notExists(link.Key.Int64, "thumb") {
					i.database.RemoveLink(link)
					os.Remove(filePath)
					fmt.Println(fmt.Sprintf("Cannot make thumbnail %s", link.Link.String))
				}
			}
		} else {
			fmt.Println(fmt.Sprintf("Can't make thumbnail, unknown image format %s %s", contentType, link.Link.String))
			continue
		}
	}
}

func (i *imageAction) Gif(cs chan *DBLink) {
	for link := range cs {
		_, data, err := i.Download(link.Link.String)
		if err != nil {
			i.database.RemoveLink(link)
			fmt.Println(fmt.Sprintf("Removing dead link %s", link.Link))
			continue
		}
		if _, err := i.database.GetRaw(link.Key.Int64, RawWebm); err != nil {
			i.database.AddRaw(RawWebm, link, "video/webm")
		}
		if _, err := i.database.GetRaw(link.Key.Int64, RawWebmFrame); err != nil {
			i.database.AddRaw(RawWebmFrame, link, "image/jpeg")
		}

		filePath := path.Join(GetConfig().DataPath, "/gif/", strconv.FormatInt(link.Key.Int64, 10))
		outPath := path.Join(GetConfig().DataPath, "/webm/", strconv.FormatInt(link.Key.Int64, 10))
		framePath := path.Join(GetConfig().DataPath, "/thumb/", strconv.FormatInt(link.Key.Int64, 10))
		if notExists(link.Key.Int64, "webm") {
			ioutil.WriteFile(filePath, data, os.ModePerm)
			exec.Command("sh", "-c", fmt.Sprintf("ffmpeg -i %s -c:v libvpx -crf 12 -b:v 500k -f webm %s && rm %s", filePath, outPath, filePath)).Output()
			exec.Command("sh", "-c", fmt.Sprintf("ffmpeg -i %s -f image2 -ss 00 -vframes 1 %s", outPath, framePath)).Output()
			if notExists(link.Key.Int64, "thumb") {
				i.database.RemoveLink(link)
				os.Remove(filePath)
				os.Remove(outPath)
				fmt.Println(fmt.Sprintf("Cannot convert GIF -> WEBM %s", link.Link.String))
			}
		}
	}
}

func (i *imageAction) WebM(cs chan *DBLink) {
	for link := range cs {
		contentType, data, err := i.Download(link.Link.String)
		if err != nil {
			i.database.RemoveLink(link)
			fmt.Println(fmt.Sprintf("Removing dead link %s", link.Link))
			continue
		}
		if contentType == "image/gif" {
			i.gifCh <- link
			continue
		}
		if _, err := i.database.GetRaw(link.Key.Int64, RawWebm); err != nil {
			i.database.AddRaw(RawWebm, link, "video/webm")
		}
		if _, err := i.database.GetRaw(link.Key.Int64, RawWebmFrame); err != nil {
			i.database.AddRaw(RawWebmFrame, link, "image/jpeg")
		}

		filePath := path.Join(GetConfig().DataPath, "/webm/", strconv.FormatInt(link.Key.Int64, 10))
		framePath := path.Join(GetConfig().DataPath, "/thumb/", strconv.FormatInt(link.Key.Int64, 10))
		if notExists(link.Key.Int64, "webm") {
			ioutil.WriteFile(filePath, data, os.ModePerm)
			exec.Command("sh", "-c", fmt.Sprintf("ffmpeg -i %s -f image2 -ss 00 -vframes 1 %s", filePath, framePath)).Output()
			if notExists(link.Key.Int64, "thumb") {
				i.database.RemoveLink(link)
				os.Remove(filePath)
				fmt.Println(fmt.Sprintf("Invalid WebM %s", link.Link.String))
			}
		}
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
		fmt.Println(fmt.Sprintf("Image returned with wrong status code %s", response.StatusCode))
		return "", []byte{}, fmt.Errorf("StatusCode %s", response.StatusCode)
	}
	data, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	return response.Header.Get("Content-Type"), data, err
}

func notExists(id int64, fileType string) bool {
	if _, err := os.Stat(path.Join(GetConfig().DataPath, fileType, strconv.FormatInt(id, 10))); os.IsNotExist(err) {
		return true
	}
	return false
}

func contains(list []*DBLink, element *DBLink) bool {
	for _, a := range list {
		if a.Key.Int64 == element.Key.Int64 {
			return true
		}
	}
	return false
}
