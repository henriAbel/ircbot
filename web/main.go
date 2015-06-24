package web

import (
	irc "../irc"
	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"time"
)

func StartWeb() {
	irc.GetConfig()
	jwt_middleware := &jwt.JWTMiddleware{
		Key:        []byte("VerySecrestKey"),
		Realm:      "jwt auth",
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			return userId == "admin" && password == "admin"
		}}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	/*api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/login"
		},
		IfTrue: jwt_middleware,
	})*/

	service := NewLinkService()

	api_router, _ := rest.MakeRouter(
		rest.Post("/login", jwt_middleware.LoginHandler),
		rest.Get("/refresh_token", jwt_middleware.RefreshHandler),
		rest.Get("/links", service.GetAll),
		rest.Get("/raw/:id/:type", service.Raw),
	)
	api.SetApp(api_router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
	http.Handle("/bower_components/", http.FileServer(http.Dir("./web/static/")))
	http.Handle("/", http.FileServer(http.Dir("./web/static/app/")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
