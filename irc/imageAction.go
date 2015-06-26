package irc

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	gif "image/gif"
	jpeg "image/jpeg"
	png "image/png"
	"io/ioutil"
	"net/http"
)

var ImageAction = imageAction{IrcDatabase{}}

type imageAction struct {
	database IrcDatabase
}

func (i *imageAction) CheckImage(link *DBLink) {
	cs := make(chan *DBLink)
	switch link.Link_type.String {
	case Image:
		if len(i.database.GetRaw(link.Key.Int64, RawImage).Data) == 0 {
			go i.Image(cs)
			cs <- link
		} else if len(i.database.GetRaw(link.Key.Int64, RawImageThumbnail).Data) == 0 {
			go i.Image(cs)
			cs <- link
		}
	case Gif:
		if len(i.database.GetRaw(link.Key.Int64, RawGifFrame).Data) == 0 {
			go i.Gif(cs)
			cs <- link
		}
	}
}

func (i *imageAction) Image(cs chan *DBLink) {
	link := <-cs
	contentType, data, err := download(link.Link.String)
	if err != nil {
		return
	}
	var image image.Image
	if contentType == "image/png" {
		image, err = png.Decode(bytes.NewReader(data))
		if err != nil {
			fmt.Println(fmt.Sprintf("Can't parse image %s", link.Link.String))
			return
		}
	} else if contentType == "image/jpeg" {
		image, err = jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			fmt.Println(fmt.Sprintf("Can't parse image %s", link.Link.String))
			return
		}
	} else {
		fmt.Println(fmt.Sprintf("Can't make thumbnail, unknown image format %s", link.Link.String))
		return
	}
	resizedImage := resize.Thumbnail(128, 128, image, resize.NearestNeighbor)
	b := jpegEncode(resizedImage)

	i.database.AddRaw(RawImageThumbnail, link, "image/jpeg", b.Bytes())
	i.database.AddRaw(RawImage, link, contentType, data)
}

func (i *imageAction) Gif(cs chan *DBLink) {
	link := <-cs
	_, data, err := download(link.Link.String)

	if err != nil {
		return
	}
	img, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		fmt.Println(fmt.Sprintf("Can't parse GIF %s", err))
		return
	}

	b := jpegEncode(img.Image[0])

	i.database.AddRaw(RawGifFrame, link, "image/jpeg", b.Bytes())
	i.database.AddRaw(RawGif, link, "image/gif", data)
}

func download(url string) (string, []byte, error) {
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("Can't download Image %s err %s", url, err))
		return "", []byte{}, err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
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
