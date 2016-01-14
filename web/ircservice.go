package web

import (
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"

	irc "../irc"

	"github.com/ant0ine/go-json-rest/rest"
)

type linkService struct {
	sync.RWMutex
	database irc.IrcDatabase
}

type link struct {
	Key         int64
	Link        string
	Sender_id   int64
	Date        int64
	Link_type   string
	Sender_name string
}

type statGroupLink struct {
	Count int64
	Type  string
}

type statGroupUser struct {
	Count     int64
	User_name string
	Sender_id int64
}

type stats struct {
	GroupLink *[]irc.DBStatGroupLink
	GroupUser *[]irc.DBStatGroupUser
}

type count struct {
	Count int
}

func NewLinkService() linkService {
	return linkService{
		database: irc.IrcDatabase{},
	}
}

func transform(dblinks *[]irc.DBLink) *[]link {
	if len(*dblinks) == 0 {
		return &[]link{}
	}
	var links = make([]link, len(*dblinks))
	for i := 0; i < len(*dblinks); i++ {
		dbLink := (*dblinks)[i]
		links[i] = link{
			Key:         dbLink.Key.Int64,
			Link:        dbLink.Link.String,
			Sender_id:   dbLink.Sender_id.Int64,
			Date:        dbLink.Date.Unix(),
			Link_type:   dbLink.Link_type.String,
			Sender_name: dbLink.Sender_name.String,
		}
	}
	return &links
}

func (l *linkService) GetAll(w rest.ResponseWriter, r *rest.Request) {
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	filterParam := r.URL.Query().Get("filter")
	filter := []string{}
	if len(filterParam) > 0 {
		filter = strings.Split(filterParam, ",")
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 30
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		offset = 0
	}

	l.Lock()
	dblinks := transform(l.database.GetAll(limit, offset, filter))
	l.Unlock()
	w.WriteJson(dblinks)
}

func (l *linkService) GetCount(w rest.ResponseWriter, r *rest.Request) {
	filterParam := r.URL.Query().Get("filter")
	filter := []string{}
	if len(filterParam) > 0 {
		filter = strings.Split(filterParam, ",")
	}

	l.Lock()
	dbCount := l.database.GetCount(filter)
	l.Unlock()

	w.WriteJson(count{dbCount})
}

func (l *linkService) GetAllStats(w rest.ResponseWriter, r *rest.Request) {
	groupStats := l.database.GetLinkGroupStat()
	userStats := l.database.GetUserStat()
	w.WriteJson(stats{groupStats, userStats})
}

func (l *linkService) Raw(w rest.ResponseWriter, r *rest.Request) {
	resourceType := r.PathParam("type")
	resourceId, err := strconv.ParseInt(r.PathParam("id"), 10, 64)
	if err != nil {
		rest.Error(w, "Can't parse resource id", 400)
		return
	}
	data := []byte{}

	switch resourceType {
	case irc.RawImage:
		data, err = ioutil.ReadFile(path.Join(irc.GetConfig().DataPath, "image", r.PathParam("id")))
	case irc.RawImageThumbnail:
		data, err = ioutil.ReadFile(path.Join(irc.GetConfig().DataPath, "thumb", r.PathParam("id")))
	case irc.RawWebm:
		data, err = ioutil.ReadFile(path.Join(irc.GetConfig().DataPath, "webm", r.PathParam("id")))
	case irc.RawWebmFrame:
		data, err = ioutil.ReadFile(path.Join(irc.GetConfig().DataPath, "thumb", r.PathParam("id")))
	}

	if len(data) == 0 {
		link := l.database.GetLinkById(resourceId)
		w.WriteHeader(503)
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		irc.ImageAction.AppendCheckLink(link)
	} else {
		dbRaw, _ := l.database.GetRaw(resourceId, resourceType)
		w.Header().Set("Content-Type", dbRaw.Mime_type.String)
		w.(http.ResponseWriter).Write(data)
	}
}
