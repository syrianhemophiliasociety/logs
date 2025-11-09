package main

import (
	"net/http"
	"regexp"
	"shs/actions"
	"shs/app"
	"shs/config"
	"shs/handlers/apis"
	"shs/handlers/middlewares/auth"
	"shs/handlers/middlewares/contenttype"
	"shs/handlers/middlewares/logger"
	"shs/jwt"
	"shs/log"
	"shs/mariadb"
	"shs/redis"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"
)

var (
	minifyer       *minify.M
	usecases       *actions.Actions
	authMiddleware *auth.Middleware
)

func init() {
	mariadbRepo, err := mariadb.New()
	if err != nil {
		log.Fatalln(err)
	}

	cache := redis.New()

	app := app.New(mariadbRepo, cache)
	jwtUtil := jwt.New[actions.TokenPayload]()

	usecases = actions.New(
		app,
		cache,
		jwtUtil,
	)

	authMiddleware = auth.New(usecases)

	minifyer = minify.New()
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
}

func main() {
	emailLoginApi := apis.NewUsernameLoginApi(usecases)
	meApi := apis.NewMeApi(usecases)
	accountApi := apis.NewAccountApi(usecases)
	bloodTestApi := apis.NewBloodTestApi(usecases)
	medicineApi := apis.NewMedicineApi(usecases)
	virusApi := apis.NewVirusApi(usecases)
	addressApi := apis.NewAddressApi(usecases)
	patientApi := apis.NewPatientApi(usecases)

	v1ApisHandler := http.NewServeMux()
	v1ApisHandler.HandleFunc("POST /login/username", emailLoginApi.HandleUsernameLogin)

	v1ApisHandler.HandleFunc("GET /me/auth", authMiddleware.AuthApi(meApi.HandleAuthCheck))
	v1ApisHandler.HandleFunc("GET /me/logout", authMiddleware.AuthApi(meApi.HandleLogout))

	v1ApisHandler.HandleFunc("POST /account/admin", authMiddleware.AuthApi(accountApi.HandleCreateAdminAccount))
	v1ApisHandler.HandleFunc("POST /account/secritary", authMiddleware.AuthApi(accountApi.HandleCreateSecritaryAccount))

	v1ApisHandler.HandleFunc("POST /bloodtest", authMiddleware.AuthApi(bloodTestApi.HandleCreateBloodTest))
	v1ApisHandler.HandleFunc("GET /bloodtest/{id}", authMiddleware.AuthApi(bloodTestApi.HandleGetBloodTest))
	v1ApisHandler.HandleFunc("GET /bloodtest/all", authMiddleware.AuthApi(bloodTestApi.HandleListBloodTests))
	v1ApisHandler.HandleFunc("DELETE /bloodtest/{id}", authMiddleware.AuthApi(bloodTestApi.HandleDeleteBloodTest))

	v1ApisHandler.HandleFunc("POST /virus", authMiddleware.AuthApi(virusApi.HandleCreateVirus))
	v1ApisHandler.HandleFunc("GET /virus/all", authMiddleware.AuthApi(virusApi.HandleListViri))
	v1ApisHandler.HandleFunc("DELETE /virus/{id}", authMiddleware.AuthApi(virusApi.HandleDeleteVirus))

	v1ApisHandler.HandleFunc("POST /medicine", authMiddleware.AuthApi(medicineApi.HandleCreateMedicine))
	v1ApisHandler.HandleFunc("GET /medicine/all", authMiddleware.AuthApi(medicineApi.HandleListMedicines))
	v1ApisHandler.HandleFunc("DELETE /medicine/{id}", authMiddleware.AuthApi(medicineApi.HandleDeleteMedicine))

	v1ApisHandler.HandleFunc(
		"GET /address/goveronate/{goveronate}/suburb/{suburb}/street/{street}",
		authMiddleware.AuthApi(addressApi.HandleFindAddress))

	v1ApisHandler.HandleFunc("POST /patient", authMiddleware.AuthApi(patientApi.HandleCreatePatient))

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/v1/", http.StripPrefix("/v1", contenttype.Json(v1ApisHandler)))

	log.Info("Starting http server at port " + config.Env().Port)
	switch config.Env().GoEnv {
	case config.GoEnvBeta, config.GoEnvDev, config.GoEnvTest:
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, logger.Handler(applicationHandler)))
	case config.GoEnvProd:
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, minifyer.Middleware(applicationHandler)))
	}
}
