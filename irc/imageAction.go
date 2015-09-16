package irc

import (
	"bufio"
	"bytes"
	"fmt"
	"gopkg.in/h2non/bimg.v0"
	"image"
	gif "image/gif"
	jpeg "image/jpeg"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func newImageAction() *imageAction {
	cs1 := make(chan *DBLink, 10)
	cs2 := make(chan *DBLink, 10)

	action := imageAction{
		database: IrcDatabase{},
		imageCh:  cs1,
		gifCh:    cs2,
		queue:    make([]*DBLink, 0),
	}

	go action.Image(cs1)
	go action.Gif(cs2)

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
	queue    []*DBLink
}

func (i *imageAction) CheckImage(link *DBLink) {
	i.queue = append(i.queue, link)
}

func (i *imageAction) checkImage(link *DBLink) {
	switch link.Link_type.String {
	case Image:
		if len(i.database.GetRaw(link.Key.Int64, RawImage).Data) == 0 {
			i.imageCh <- link
		} else if len(i.database.GetRaw(link.Key.Int64, RawImageThumbnail).Data) == 0 {
			i.imageCh <- link
		}
	case Gif:
		if len(i.database.GetRaw(link.Key.Int64, RawGifFrame).Data) == 0 {
			i.gifCh <- link
		}
	}
	runtime.GC()
}

func (i *imageAction) Image(cs chan *DBLink) {
	for link := range cs {
		contentType, data, err := i.Download(link.Link.String)
		if err != nil {
			i.database.RemoveLink(link)
			fmt.Println(fmt.Sprintf("Removing dead link %s", link.Link))
			continue
		}
		if contentType == "image/png" || contentType == "image/jpeg" {
			newImage, err := bimg.NewImage(data).Thumbnail(128)
			if err != nil {
				fmt.Println(fmt.Sprintf("Can't parse image %s", link.Link.String))
				i.database.RemoveLink(link)
				continue
			}
			i.database.AddRaw(RawImageThumbnail, link, "image/jpeg", newImage)
			i.database.AddRaw(RawImage, link, contentType, data)
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
		img, err := gif.Decode(bytes.NewReader(data))
		if err != nil {
			fmt.Println(fmt.Sprintf("Can't parse GIF %s", err))
			i.database.RemoveLink(link)
			continue
		}

		b := jpegEncode(img)
		i.database.AddRaw(RawGifFrame, link, "image/jpeg", b.Bytes())
		i.database.AddRaw(RawGif, link, "image/gif", data)
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

func jpegEncode(img image.Image) *bytes.Buffer {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err := jpeg.Encode(w, img, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("Can't convert to jpeg %s", err))
		panic(err)
	}
	return &b
}
