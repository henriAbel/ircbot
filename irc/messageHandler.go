package irc

import (
	"bufio"
	"bytes"
	"fmt"
	simpleJson "github.com/bitly/go-simplejson"
	client "github.com/fluffle/goirc/client"
	gif "image/gif"
	jpeg "image/jpeg"
	"io/ioutil"
	"net/http"
	urlLib "net/url"
	"regexp"
	"strings"
	"time"
)

var Images = []string{".jpg", ".png", ".jpeg"}

const (
	Image   = "image"
	Gif     = "gif"
	Link    = "link"
	Youtube = "youtube"
)

type Handler struct {
	conn         *client.Conn
	urlRegex     *regexp.Regexp
	commandRegex *regexp.Regexp
	channel      string
	database     IrcDatabase
}

func NewHandler(conn *client.Conn, channel string) Handler {
	r, _ := regexp.Compile("http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*(),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+")
	r2, _ := regexp.Compile("^!([^\\s]+)")
	h := Handler{
		conn:         conn,
		urlRegex:     r,
		commandRegex: r2,
		channel:      channel,
		database:     IrcDatabase{},
	}
	// Just for error checking
	db := h.database.Open()
	db.Close()
	return h
}

func (h *Handler) Recv(line string, sender string) {
	url := h.urlRegex.FindString(line)
	if len(url) > 0 {
		cs := make(chan *DBLink)
		var linkType = Link
		if strings.Index(url, "youtube.com/watch") > 0 {
			parsedUrl, _ := urlLib.Parse(url)
			query, _ := urlLib.ParseQuery(parsedUrl.RawQuery)
			h.Youtube(query["v"][0], sender)
			linkType = Youtube

		} else if strings.Index(url, "://youtu.be/") > 0 {
			parsedUrl, _ := urlLib.Parse(url)
			h.Youtube(parsedUrl.Path[1:], sender)
			linkType = Youtube
		} else {
			urlSuffix := url[strings.LastIndex(url, "."):]
			if urlSuffix == ".gif" {
				linkType = Gif
				go h.Gif(cs)
			} else if inArray(urlSuffix, Images) {
				linkType = Image
				go h.Image(cs)
			}
		}

		link := h.database.GetLink(url)
		if (*link == DBLink{}) {
			l := h.database.AddLink(url, linkType, sender)
			select {
			case cs <- l:
			default:
			}
		} else {
			loc, _ := time.LoadLocation("Europe/Tallinn")
			h.conn.Privmsg(h.channel, fmt.Sprintf("OLD! Selle on varem saatnud juba %s %s!\n", link.Sender_name.String, link.Date.In(loc).Format("2006-01-02 15:04")))
			h.database.LogDuplicate(link, sender)
		}
	}
	if strings.HasPrefix(line, "!") {
		command := h.commandRegex.FindAllString(line, 1)[0]
		if command == "!google" {
			h.GoogleSearch(strings.TrimSpace(strings.Replace(line, command, "", -1)))
		}
	}
}

func (h *Handler) Image(cs chan *DBLink) {
	link := <-cs
	response, err := http.Get(link.Link.String)
	if err != nil || response.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("Can't download Image %s err %s", link.Link.String, err))
		return
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	h.database.AddRaw(RawImage, link, response.Header.Get("Content-Type"), data)
}

func (h *Handler) Gif(cs <-chan *DBLink) {
	fmt.Println("method gif")
	link := <-cs
	fmt.Println("channel gif")
	url := link.Link.String
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("Can't download GIF %s err %s", url, err))
		return
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("Can't open http response %s", err))
	}
	img, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		fmt.Println(fmt.Sprintf("Cant parse GIF %s", err))
		return
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err = jpeg.Encode(w, img.Image[0], nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("Can't convert to jpeg %s", err))
	}
	h.database.AddRaw(RawGifFrame, link, "image/jpeg", b.Bytes())
	h.database.AddRaw(RawGif, link, "image/gif", data)
}

func (h *Handler) Youtube(videoId string, sender string) {
	var apiUrl = fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet&id=%s&key=%s", videoId, GetConfig().YoutubeApiKey)
	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	json, _ := simpleJson.NewFromReader(res.Body)
	title, _ := json.Get("items").GetIndex(0).Get("snippet").Get("title").String()
	h.conn.Privmsg(h.channel, title)
}

func (h *Handler) GoogleSearch(query string) {
	var queryUrl = fmt.Sprintf("http://ajax.googleapis.com/ajax/services/search/web?v=1.0&q=%s", urlLib.QueryEscape(query))
	res, err := http.Get(queryUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	json, _ := simpleJson.NewFromReader(res.Body)
	result, _ := json.Get("responseData").Get("results").GetIndex(0).Get("url").String()
	h.conn.Privmsg(h.channel, result)
}

func inArray(a string, list []string) bool {
	for _, b := range list {
		if strings.HasPrefix(a, b) {
			return true
		}
	}
	return false
}
