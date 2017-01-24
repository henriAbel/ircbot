package web

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	irc "../irc"
	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/soheilhy/cmux"
)

func serveHTTP1(l net.Listener, m *http.ServeMux) {
	s := &http.Server{
		Handler: m,
	}
	if err := s.Serve(l); err != cmux.ErrListenerClosed {
		panic(err)
	}
}

func serveHTTPS(l net.Listener, m *http.ServeMux) {
	certificate, err := tls.LoadX509KeyPair(irc.GetCertFile(), irc.GetKeyFile())
	if err != nil {
		log.Panic(err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		Rand:         rand.Reader,
	}
	tlsl := tls.NewListener(l, config)
	serveHTTP1(tlsl, m)
}

func StartWeb() {
	var DefaultDevStack = []rest.Middleware{
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.RecoverMiddleware{
			EnableResponseStackTrace: true,
		},
		&rest.JsonIndentMiddleware{},
		&rest.ContentTypeCheckerMiddleware{},
	}

	jwtMiddleware := &jwt.JWTMiddleware{
		Key:        []byte("VerySecrestKey"),
		Realm:      "jwt auth",
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			return password == irc.GetWebPassword()
		}}

	api := rest.NewApi()
	api.Use(DefaultDevStack...)
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			if len(irc.GetWebPassword()) > 1 {
				token := request.Request.URL.Query()["authorization"]
				if len(token) > 0 {
					request.Header.Add("Authorization", "Bearer "+token[0])
				}
				return request.URL.Path != "/login" && request.URL.Path != "/config"
			}
			return false
		},
		IfTrue: jwtMiddleware,
	})

	service := NewLinkService()

	apiRouter, _ := rest.MakeRouter(
		rest.Post("/login", jwtMiddleware.LoginHandler),
		rest.Get("/refresh_token", jwtMiddleware.RefreshHandler),
		rest.Get("/config", service.GetConfig),
		rest.Get("/links", service.GetAll),
		rest.Get("/links/count", service.GetCount),
		rest.Get("/raw/:id/:type", service.Raw),
		rest.Get("/stat/all", service.GetAllStats),
	)
	api.SetApp(apiRouter)

	listeningPort := fmt.Sprintf(":%d", irc.GetWebPort())
	enableTls := len(irc.GetCertFile()) > 0 && len(irc.GetKeyFile()) > 0
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	mux.Handle("/bower_components/", http.FileServer(http.Dir("./web/static/")))
	mux.Handle("/", http.FileServer(http.Dir("./web/static/dist/")))

	l, err := net.Listen("tcp", listeningPort)
	if err != nil {
		log.Fatal(err)
	}
	m := cmux.New(l)

	httpl := m.Match(cmux.HTTP1Fast())
	if enableTls {
		tlsl := m.Match(cmux.Any())
		redirectMux := http.NewServeMux()
		redirectMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			u := r.URL
			u.Host = r.Host
			u.Scheme = "https"
			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
		})
		go serveHTTPS(tlsl, mux)
		go serveHTTP1(httpl, redirectMux)
	} else {
		go serveHTTP1(httpl, mux)
	}

	if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}

}
