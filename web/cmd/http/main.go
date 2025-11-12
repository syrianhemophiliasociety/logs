package main

import (
	"net/http"
	"os"
	"regexp"
	"shs-web/actions"
	"shs-web/config"
	"shs-web/handlers/apis"
	"shs-web/handlers/middlewares/auth"
	"shs-web/handlers/middlewares/clienthash"
	"shs-web/handlers/middlewares/contenttype"
	"shs-web/handlers/middlewares/ismobile"
	"shs-web/handlers/middlewares/locale"
	"shs-web/handlers/middlewares/logger"
	"shs-web/handlers/middlewares/theme"
	"shs-web/handlers/middlewares/version"
	"shs-web/handlers/pages"
	"shs-web/handlers/static"
	"shs-web/log"
	"shs-web/redis"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var (
	minifyer *minify.M

	usecases       *actions.Actions
	authMiddleware *auth.Middleware

	appVersion = "git-latest"
)

func init() {
	minifyer = minify.New()
	minifyer.AddFunc("text/css", css.Minify)
	minifyer.AddFunc("text/html", html.Minify)
	minifyer.AddFunc("image/svg+xml", svg.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	cache := redis.New()
	usecases = actions.New(cache)
	authMiddleware = auth.New(usecases)

	v := os.Getenv("VERSION")
	if v != "" {
		appVersion = v
	}
}

func main() {
	pagesHandler := http.NewServeMux()

	pagesHandler.HandleFunc("/robots.txt", static.HandleRobots)
	pagesHandler.HandleFunc("/sitemap.xml", static.HandleSitemap)
	pagesHandler.HandleFunc("/favicon.ico", static.HandleFavicon)

	pages := pages.New(usecases)
	pagesHandler.HandleFunc("/", contenttype.Html(authMiddleware.AuthPage(pages.HandleHomePage)))
	pagesHandler.HandleFunc("GET /about", contenttype.Html(pages.HandleAboutPage))
	pagesHandler.HandleFunc("GET /privacy", contenttype.Html(pages.HandlePrivacyPage))
	pagesHandler.HandleFunc("GET /login", contenttype.Html(authMiddleware.AuthPage(pages.HandleLoginPage)))
	pagesHandler.HandleFunc("GET /viruses", contenttype.Html(authMiddleware.AuthPage(pages.HandleVirusesPage)))
	pagesHandler.HandleFunc("GET /medicines", contenttype.Html(authMiddleware.AuthPage(pages.HandleMedicinesPage)))
	pagesHandler.HandleFunc("GET /blood-tests", contenttype.Html(authMiddleware.AuthPage(pages.HandleBloodTestsPage)))

	usernameLoginApi := apis.NewUsernameLoginApi(usecases)
	logoutApi := apis.NewLogoutApi(usecases)
	virusApi := apis.NewVirusApi(usecases)
	medicineApi := apis.NewMedicineApi(usecases)
	bloodTestApi := apis.NewBloodTestApi(usecases)

	apisHandler := http.NewServeMux()
	apisHandler.HandleFunc("POST /login/username", usernameLoginApi.HandleUsernameLogin)
	apisHandler.HandleFunc("GET /logout", authMiddleware.AuthApi(logoutApi.HandleLogout))

	apisHandler.HandleFunc("POST /virus", authMiddleware.AuthApi(virusApi.HandleCreateVirus))
	apisHandler.HandleFunc("DELETE /virus/{id}", authMiddleware.AuthApi(virusApi.HandleDeleteVirus))

	apisHandler.HandleFunc("POST /medicine", authMiddleware.AuthApi(medicineApi.HandleCreateMedicine))
	apisHandler.HandleFunc("DELETE /medicine/{id}", authMiddleware.AuthApi(medicineApi.HandleDeleteMedicine))

	apisHandler.HandleFunc("POST /blood-test", authMiddleware.AuthApi(bloodTestApi.HandleCreateBloodTest))
	apisHandler.HandleFunc("DELETE /blood-test/{id}", authMiddleware.AuthApi(bloodTestApi.HandleDeleteBloodTest))

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", locale.Handler(ismobile.Handler(theme.Handler(pagesHandler))))
	applicationHandler.Handle("/api/", ismobile.Handler(theme.Handler(http.StripPrefix("/api", apisHandler))))
	applicationHandler.Handle("/assets/", http.StripPrefix("/assets", static.AssetsHandler(minifyer)))

	log.Info("Starting http server at port " + config.Env().Port)
	if config.Env().GoEnv == "dev" || config.Env().GoEnv == "beta" {
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, version.Handler(appVersion, clienthash.Handler(logger.Handler(ismobile.Handler(theme.Handler(applicationHandler)))))))
	}
	log.Fatalln(http.ListenAndServe(":"+config.Env().Port, version.Handler(appVersion, clienthash.Handler(ismobile.Handler(theme.Handler(minifyer.Middleware(applicationHandler)))))))
}
