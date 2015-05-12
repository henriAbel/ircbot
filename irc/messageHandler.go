package irc

import (
	//"encoding/json"
	"fmt"
	simpleJson "github.com/bitly/go-simplejson"
	client "github.com/fluffle/goirc/client"
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
	conn     *client.Conn
	urlRegex *regexp.Regexp
	channel  string
	database IrcDatabase
}

func NewHandler(conn *client.Conn, channel string) Handler {
	r, _ := regexp.Compile("http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*(),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+")
	h := Handler{
		conn:     conn,
		urlRegex: r,
		channel:  channel,
		database: IrcDatabase{},
	}
	// Just for error checking
	db := h.database.Open()
	db.Close()
	return h
}

func (h *Handler) Recv(line string, sender string) {
	url := h.urlRegex.FindString(line)
	if len(url) > 0 {
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
			} else if inArray(urlSuffix, Images) {
				linkType = Image
			}
		}

		link := h.database.GetLink(url)
		if (*link == DBLink{}) {
			h.database.AddLink(url, linkType, sender)
		} else {
			loc, _ := time.LoadLocation("Europe/Tallinn")
			h.conn.Privmsg(h.channel, fmt.Sprintf("OLD! Selle on varem saatnud juba %s %s!\n", link.sender_name.String, link.date.In(loc).Format("2006-01-02 15:04")))
			h.database.LogDuplicate(link, sender)
		}
	}
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

func inArray(a string, list []string) bool {
	for _, b := range list {
		if strings.HasPrefix(a, b) {
			return true
		}
	}
	return false
}
