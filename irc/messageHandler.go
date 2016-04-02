package irc

import (
	"fmt"
	"net/http"
	urlLib "net/url"
	"regexp"
	"strings"
	"time"

	simpleJson "github.com/bitly/go-simplejson"
	client "github.com/fluffle/goirc/client"
)

var Images = []string{".jpg", ".png", ".jpeg"}

const (
	Image   = "image"
	Link    = "link"
	Youtube = "youtube"
	WebM    = "webm"
	Gif     = "gif"
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
		var linkType = Link
		parsedUrl, _ := urlLib.Parse(url)
		// Modify dropbox urls
		// TODO Is this good enough??
		if strings.Contains(url, "dropbox.com/s/") {
			values := parsedUrl.Query()
			values.Del("dl")
			values.Add("raw", "1")
			url = strings.Replace(url, url[strings.LastIndex(url, "?")+1:], values.Encode(), -1)
		}
		if strings.Index(url, "youtube.com/watch") > 0 {
			query, _ := urlLib.ParseQuery(parsedUrl.RawQuery)
			h.Youtube(query["v"][0], sender)
			linkType = Youtube

		} else if strings.Index(url, "://youtu.be/") > 0 {
			h.Youtube(parsedUrl.Path[1:], sender)
			linkType = Youtube
		} else {
			urlSuffix := url[strings.LastIndex(url, "."):]
			if urlSuffix == ".gif" {
				linkType = Gif
			} else if StartsWith(urlSuffix, Images) {
				linkType = Image
			} else if urlSuffix == ".webm" {
				linkType = WebM
			} else if urlSuffix == ".gifv" {
				linkType = WebM
				url = strings.Replace(url, ".gifv", ".webm", 1)
			}
		}

		link := h.database.GetLink(url)
		if (*link == DBLink{}) {
			l := h.database.AddLink(url, linkType, sender)
			go ImageAction.CheckLink(l)
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

func (h *Handler) Youtube(videoId string, sender string) {
	var apiUrl = fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet&id=%s&key=%s", videoId, GetYoutubeApiKey())
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
	unescaped, _ := urlLib.QueryUnescape(result)
	h.conn.Privmsg(h.channel, unescaped)
}

func StartsWith(a string, list []string) bool {
	for _, b := range list {
		if strings.HasPrefix(a, b) {
			return true
		}
	}
	return false
}
