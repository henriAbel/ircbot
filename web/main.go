package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	irc "../irc"
	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
)

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
				return request.URL.Path != "/login"
			}
			return false
		},
		IfTrue: jwtMiddleware,
	})

	service := NewLinkService()

	apiRouter, _ := rest.MakeRouter(
		rest.Post("/login", jwtMiddleware.LoginHandler),
		rest.Get("/refresh_token", jwtMiddleware.RefreshHandler),
		rest.Get("/links", service.GetAll),
		rest.Get("/links/count", service.GetCount),
		rest.Get("/raw/:id/:type", service.Raw),
		rest.Get("/stat/all", service.GetAllStats),
	)
	api.SetApp(apiRouter)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/bower_components/", http.FileServer(http.Dir("./web/static/")))
	http.Handle("/", http.FileServer(http.Dir("./web/static/app/")))
	listeningPort := fmt.Sprintf(":%d", irc.GetWebPort())
	log.Fatal(http.ListenAndServe(listeningPort, nil))
}
