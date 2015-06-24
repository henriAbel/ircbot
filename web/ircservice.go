package web

import (
	irc "../irc"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type LinkService struct {
	sync.RWMutex
	database irc.IrcDatabase
}

type Link struct {
	Key         int64
	Link        string
	Sender_id   int64
	Date        int64
	Link_type   string
	Sender_name string
}

func NewLinkService() LinkService {
	return LinkService{
		database: irc.IrcDatabase{},
	}
}

func transform(dblinks *[]irc.DBLink) *[]Link {
	if len(*dblinks) == 0 {
		return &[]Link{}
	}
	var links = make([]Link, len(*dblinks))
	for i := 0; i < len(*dblinks); i++ {
		dbLink := (*dblinks)[i]
		links[i] = Link{
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

func (l *LinkService) GetAll(w rest.ResponseWriter, r *rest.Request) {
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

func (l *LinkService) Raw(w rest.ResponseWriter, r *rest.Request) {
	resourceType := r.PathParam("type")

	resourceId, err := strconv.Atoi(r.PathParam("id"))
	if err != nil {
		rest.Error(w, "Can't parse resource id", 400)
		return
	}

	l.Lock()
	rawData := l.database.GetRaw(resourceId, resourceType)
	l.Unlock()
	if len(rawData.Data) == 0 {
		link := l.database.GetLinkById(resourceId)
		w.Header().Set("Location", link.Link.String)
		w.WriteHeader(303)
	} else {
		w.Header().Set("Content-Type", rawData.Mime_type.String)
		w.(http.ResponseWriter).Write(rawData.Data)
	}
}
