package web

import (
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	irc "../irc"
	//log "github.com/Sirupsen/logrus"

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

type stats struct {
	GroupLink  *[]irc.DBStatGroupLink
	GroupUser  *[]irc.DBStatGroupUser
	Duplicates *[]irc.DBStatDuplicate
}

type count struct {
	Count int
}

type config struct {
	ExternalLibraries bool
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

func (l *linkService) GetConfig(w rest.ResponseWriter, r *rest.Request) {
	c := config{irc.ImageAction.ExternalLibraries}
	w.WriteJson(c)
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
	duplicates := l.database.GetDuplicateStat()
	w.WriteJson(stats{groupStats, userStats, duplicates})
}

func (l *linkService) Raw(w rest.ResponseWriter, r *rest.Request) {
	resourceType := r.PathParam("type")
	resourceId, err := strconv.ParseInt(r.PathParam("id"), 10, 64)
	if err != nil {
		rest.Error(w, "Can't parse resource id", 400)
		return
	}

	var file *os.File

	switch resourceType {
	case irc.RawImage:
		file, err = os.Open(path.Join(irc.GetDataPath(), "image", r.PathParam("id")))
	case irc.RawImageThumbnail:
		file, err = os.Open(path.Join(irc.GetDataPath(), "thumb", r.PathParam("id")))
	case irc.RawWebm:
		file, err = os.Open(path.Join(irc.GetDataPath(), "webm", r.PathParam("id")))
	case irc.RawWebmFrame:
		file, err = os.Open(path.Join(irc.GetDataPath(), "thumb", r.PathParam("id")))
	case irc.RawGif:
		file, err = os.Open(path.Join(irc.GetDataPath(), "gif", r.PathParam("id")+".gif"))
	}

	if err != nil || file == nil {
		link := l.database.GetLinkById(resourceId)
		w.WriteHeader(503)
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		irc.ImageAction.AppendCheckLink(link)
	} else {
		dbRaw, _ := l.database.GetRaw(resourceId, resourceType)
		fi, _ := file.Stat()
		w.Header().Set("Content-Type", dbRaw.Mime_type.String)
		w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
		http.ServeContent(w.(http.ResponseWriter), r.Request, "", time.Now(), file)
	}
}
