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

func StartWeb(config *irc.Config) {
	jwt_middleware := &jwt.JWTMiddleware{
		Key:        []byte("VerySecrestKey"),
		Realm:      "jwt auth",
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			return password == config.WebPassword
		}}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			if len(config.WebPassword) > 1 {
				token := request.Request.URL.Query()["authorization"]
				if len(token) > 0 {
					request.Header.Add("Authorization", "Bearer "+token[0])
				}
				return request.URL.Path != "/login"
			}
			return false
		},
		IfTrue: jwt_middleware,
	})

	service := NewLinkService()

	api_router, _ := rest.MakeRouter(
		rest.Post("/login", jwt_middleware.LoginHandler),
		rest.Get("/refresh_token", jwt_middleware.RefreshHandler),
		rest.Get("/links", service.GetAll),
		rest.Get("/links/count", service.GetCount),
		rest.Get("/raw/:id/:type", service.Raw),
	)
	api.SetApp(api_router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/bower_components/", http.FileServer(http.Dir("./web/static/")))
	http.Handle("/", http.FileServer(http.Dir("./web/static/app/")))
	listeningPort := fmt.Sprintf(":%d", irc.GetConfig().WebPort)
	log.Fatal(http.ListenAndServe(listeningPort, nil))
}
