package main

import (
	"net/http"
	"os"
	"regexp"
	"shs/actions"
	"shs/app"
	"shs/config"
	"shs/handlers/apis"
	"shs/handlers/middlewares/auth"
	"shs/handlers/middlewares/contenttype"
	"shs/handlers/middlewares/ismobile"
	"shs/handlers/middlewares/logger"
	"shs/handlers/middlewares/version"
	"shs/handlers/middlewares/webauth"
	"shs/handlers/middlewares/webi18n"
	"shs/handlers/middlewares/webtheme"
	webapis "shs/handlers/web/apis"
	webhtmx "shs/handlers/web/htmx"
	"shs/handlers/web/pages"
	"shs/handlers/web/static"
	"shs/jwt"
	"shs/log"
	"shs/mariadb"
	"shs/redis"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var appVersion = os.Getenv("VERSION")

func main() {
	repo, err := mariadb.New()
	if err != nil {
		log.Fatalln(err)
	}
	cache := redis.New()
	app := app.New(repo, cache)
	jwtUtil := jwt.New[actions.TokenPayload]()
	usecases := actions.New(
		app,
		cache,
		jwtUtil,
	)
	authMiddleware := auth.New(usecases)
	webAuthMiddleware := webauth.New(usecases)
	minifyer := minify.New()
	minifyer.AddFunc("text/css", css.Minify)
	minifyer.AddFunc("text/html", html.Minify)
	minifyer.AddFunc("image/svg+xml", svg.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	minifyer.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	///
	/// PAGES
	///

	pagesHandler := http.NewServeMux()

	pagesHandler.HandleFunc("/robots.txt", static.HandleRobots)
	pagesHandler.HandleFunc("/sitemap.xml", static.HandleSitemap)
	pagesHandler.HandleFunc("/favicon.ico", static.HandleFavicon)

	pages := pages.New(usecases)
	pagesHandler.HandleFunc("/", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleHomePage)))
	pagesHandler.HandleFunc("GET /about", contenttype.Html(webAuthMiddleware.OptionalAuthPage(pages.HandleAboutPage)))
	pagesHandler.HandleFunc("GET /privacy", contenttype.Html(webAuthMiddleware.OptionalAuthPage(pages.HandlePrivacyPage)))
	pagesHandler.HandleFunc("GET /login", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleLoginPage)))
	pagesHandler.HandleFunc("GET /viruses", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleVirusesPage)))
	pagesHandler.HandleFunc("GET /medicines", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleMedicinesPage)))
	pagesHandler.HandleFunc("GET /medicine/{id}", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleMedicinePage)))
	pagesHandler.HandleFunc("GET /blood-tests", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleBloodTestsPage)))
	pagesHandler.HandleFunc("GET /blood-test/{id}", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleBloodTestPage)))
	pagesHandler.HandleFunc("GET /management", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleManagementPage)))
	pagesHandler.HandleFunc("GET /management/account/{id}", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleAccountManagementPage)))
	pagesHandler.HandleFunc("GET /patients", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandlePatientsPage)))
	pagesHandler.HandleFunc("GET /patient/{id}", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandlePatientPage)))
	pagesHandler.HandleFunc("GET /patient/{id}/blood-test-result/{btr_id}", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandlePatientBloodTestResultPage)))
	pagesHandler.HandleFunc("GET /patient/{id}/visit/{visit_id}", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandlePatientVisitPage)))
	pagesHandler.HandleFunc("GET /diagnoses", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleDiagnosesPage)))
	pagesHandler.HandleFunc("GET /statistics", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandleStatisticsPage)))

	pagesHandler.HandleFunc("GET /patient/medications", contenttype.Html(webAuthMiddleware.AuthPage(pages.HandlePatientMedicationsPage)))

	///
	/// REST APIS
	///

	emailLoginApi := apis.NewUsernameLoginApi(usecases)
	meApi := apis.NewMeApi(usecases)
	accountApi := apis.NewAccountApi(usecases)
	bloodTestApi := apis.NewBloodTestApi(usecases)
	medicineApi := apis.NewMedicineApi(usecases)
	virusApi := apis.NewVirusApi(usecases)
	addressApi := apis.NewAddressApi(usecases)
	patientApi := apis.NewPatientApi(usecases)
	diagnosisApi := apis.NewDiagnosisApi(usecases)

	v1ApisHandler := http.NewServeMux()
	v1ApisHandler.HandleFunc("POST /login/username", emailLoginApi.HandleUsernameLogin)

	v1ApisHandler.HandleFunc("GET /me/auth", authMiddleware.AuthApi(meApi.HandleAuthCheck))
	v1ApisHandler.HandleFunc("GET /me/logout", authMiddleware.AuthApi(meApi.HandleLogout))

	v1ApisHandler.HandleFunc("GET /accounts/{id}", authMiddleware.AuthApi(accountApi.HandleGetAccount))
	v1ApisHandler.HandleFunc("DELETE /accounts/{id}", authMiddleware.AuthApi(accountApi.HandleDeleteAccount))
	v1ApisHandler.HandleFunc("PUT /accounts/{id}", authMiddleware.AuthApi(accountApi.HandleUpdateAccount))
	v1ApisHandler.HandleFunc("POST /accounts/admin", authMiddleware.AuthApi(accountApi.HandleCreateAdminAccount))
	v1ApisHandler.HandleFunc("POST /accounts/secritary", authMiddleware.AuthApi(accountApi.HandleCreateSecritaryAccount))
	v1ApisHandler.HandleFunc("POST /accounts/jointlogist", authMiddleware.AuthApi(accountApi.HandleCreateJointlogistAccount))
	v1ApisHandler.HandleFunc("GET /accounts", authMiddleware.AuthApi(accountApi.HandleListAllAccounts))

	v1ApisHandler.HandleFunc("POST /bloodtests", authMiddleware.AuthApi(bloodTestApi.HandleCreateBloodTest))
	v1ApisHandler.HandleFunc("GET /bloodtests/{id}", authMiddleware.AuthApi(bloodTestApi.HandleGetBloodTest))
	v1ApisHandler.HandleFunc("GET /bloodtests", authMiddleware.AuthApi(bloodTestApi.HandleListBloodTests))
	v1ApisHandler.HandleFunc("DELETE /bloodtests/{id}", authMiddleware.AuthApi(bloodTestApi.HandleDeleteBloodTest))

	v1ApisHandler.HandleFunc("POST /diagnoses", authMiddleware.AuthApi(diagnosisApi.HandleCreateDiagnosis))
	v1ApisHandler.HandleFunc("GET /diagnoses", authMiddleware.AuthApi(diagnosisApi.HandleListDiagnosiss))
	v1ApisHandler.HandleFunc("DELETE /diagnoses/{id}", authMiddleware.AuthApi(diagnosisApi.HandleDeleteDiagnosis))

	v1ApisHandler.HandleFunc("POST /viruses", authMiddleware.AuthApi(virusApi.HandleCreateVirus))
	v1ApisHandler.HandleFunc("GET /viruses", authMiddleware.AuthApi(virusApi.HandleListViruses))
	v1ApisHandler.HandleFunc("DELETE /viruses/{id}", authMiddleware.AuthApi(virusApi.HandleDeleteVirus))

	v1ApisHandler.HandleFunc("POST /medicines", authMiddleware.AuthApi(medicineApi.HandleCreateMedicine))
	v1ApisHandler.HandleFunc("GET /medicines", authMiddleware.AuthApi(medicineApi.HandleListMedicines))
	v1ApisHandler.HandleFunc("GET /medicines/{id}", authMiddleware.AuthApi(medicineApi.HandleGetMedicine))
	v1ApisHandler.HandleFunc("PUT /medicines/{id}/amount", authMiddleware.AuthApi(medicineApi.HandleUpdateMedicineAmount))
	v1ApisHandler.HandleFunc("DELETE /medicines/{id}", authMiddleware.AuthApi(medicineApi.HandleDeleteMedicine))

	v1ApisHandler.HandleFunc(
		"GET /addresses/goveronate/{goveronate}/suburb/{suburb}/street/{street}",
		authMiddleware.AuthApi(addressApi.HandleFindAddress))

	v1ApisHandler.HandleFunc("POST /patients", authMiddleware.AuthApi(patientApi.HandleCreatePatient))
	v1ApisHandler.HandleFunc("GET /patients/{id}/card", authMiddleware.AuthApi(patientApi.HandleGenerateCard))
	v1ApisHandler.HandleFunc("DELETE /patients/{id}", authMiddleware.AuthApi(patientApi.HandleDeletePatient))
	v1ApisHandler.HandleFunc("GET /patients/{id}", authMiddleware.AuthApi(patientApi.HandleGetPatient))
	v1ApisHandler.HandleFunc("GET /patients/last", authMiddleware.AuthApi(patientApi.HandleListLastPatients))
	v1ApisHandler.HandleFunc(
		"GET /patients/public-id/{public_id}/first-name/{first_name}/last-name/{last_name}/father-name/{father_name}/mother-name/{mother_name}/national-id/{national_id}/phone-number/{phone_number}",
		authMiddleware.AuthApi(patientApi.HandleFindPatients))
	v1ApisHandler.HandleFunc("POST /patients/import/csv", authMiddleware.AuthApi(patientApi.HandleImportPatientsFromCsv))

	v1ApisHandler.HandleFunc("POST /patients/bloodtest", authMiddleware.AuthApi(patientApi.HandleCreatePatientBloodTestResult))
	v1ApisHandler.HandleFunc("PUT /patients/{id}/bloodtest/{btr_id}/pending", authMiddleware.AuthApi(patientApi.HandleUpdatePendingBloodTestResult))
	v1ApisHandler.HandleFunc("POST /patients/{id}/checkup", authMiddleware.AuthApi(patientApi.HandleCheckUp))
	v1ApisHandler.HandleFunc("POST /patients/diagnosis", authMiddleware.AuthApi(patientApi.HandleCreatePatientDiagnosisResult))
	v1ApisHandler.HandleFunc("POST /patients/{id}/joints-evaluation", authMiddleware.AuthApi(patientApi.HandleCreatePatientJointsEvaluation))
	v1ApisHandler.HandleFunc("GET /patients/{id}/joints-evaluations", authMiddleware.AuthApi(patientApi.HandleListPatientJointsEvaluations))
	v1ApisHandler.HandleFunc("GET /patients/{id}/visits", authMiddleware.AuthApi(patientApi.HandleListPatientVisits))

	// TODO: separate this from admin patient endpoints
	v1ApisHandler.HandleFunc("POST /patients/visit/{visit_id}/medicine/{med_id}", authMiddleware.AuthApi(patientApi.HandleUsePrescribedMedicineForVisit))

	v1ApisHandler.HandleFunc("GET /me/patient/last-visit", authMiddleware.AuthApi(patientApi.HandleGetPatientLastVisit))

	if config.Env().GoEnv == config.GoEnvTest || config.Env().GoEnv == config.GoEnvDev {
		v1ApisHandler.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"message": "yeeehaww"}`))
		})

		v1ApisHandler.HandleFunc("POST /tests/reset/db", func(w http.ResponseWriter, r *http.Request) {
			err := repo.DeleteAll()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"message": "resetting DB failed"}`))
				return
			}

			_ = repo.CreateSuperAdmin()

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"message": "yeeehaww"}`))
		})

		v1ApisHandler.HandleFunc("POST /tests/reset/cache", func(w http.ResponseWriter, r *http.Request) {
			err := cache.FlushAll()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"message": "flushing cache failed"}`))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"message": "yeeehaww"}`))
		})
	}

	///
	/// WEB APIS
	///

	usernameLoginWebApi := webapis.NewUsernameLoginApi(usecases)
	logoutWebApi := webapis.NewLogoutApi(usecases)
	virusWebApi := webapis.NewVirusApi(usecases)
	medicineWebApi := webapis.NewMedicineApi(usecases)
	bloodTestWebApi := webapis.NewBloodTestApi(usecases)
	accountWebApi := webapis.NewAccountApi(usecases)
	patientWebApi := webapis.NewPatientApi(usecases)
	diagnosisWebApi := webapis.NewDiagnosisApi(usecases)

	webApisHandler := http.NewServeMux()
	webApisHandler.HandleFunc("POST /login/username", usernameLoginWebApi.HandleUsernameLogin)
	webApisHandler.HandleFunc("GET /logout", webAuthMiddleware.AuthApi(logoutWebApi.HandleLogout))

	webApisHandler.HandleFunc("POST /virus", webAuthMiddleware.AuthApi(virusWebApi.HandleCreateVirus))
	webApisHandler.HandleFunc("DELETE /virus/{id}", webAuthMiddleware.AuthApi(virusWebApi.HandleDeleteVirus))

	webApisHandler.HandleFunc("POST /medicine", webAuthMiddleware.AuthApi(medicineWebApi.HandleCreateMedicine))
	webApisHandler.HandleFunc("DELETE /medicine/{id}", webAuthMiddleware.AuthApi(medicineWebApi.HandleDeleteMedicine))
	webApisHandler.HandleFunc("PUT /medicine/{id}", webAuthMiddleware.AuthApi(medicineWebApi.HandleUpdateMedicine))

	webApisHandler.HandleFunc("POST /blood-test", webAuthMiddleware.AuthApi(bloodTestWebApi.HandleCreateBloodTest))
	webApisHandler.HandleFunc("DELETE /blood-test/{id}", webAuthMiddleware.AuthApi(bloodTestWebApi.HandleDeleteBloodTest))

	webApisHandler.HandleFunc("POST /diagnosis", webAuthMiddleware.AuthApi(diagnosisWebApi.HandleCreateDiagnosis))
	webApisHandler.HandleFunc("DELETE /diagnosis/{id}", webAuthMiddleware.AuthApi(diagnosisWebApi.HandleDeleteDiagnosis))

	webApisHandler.HandleFunc("POST /account", webAuthMiddleware.AuthApi(accountWebApi.HandleCreateAccount))
	webApisHandler.HandleFunc("PUT /account/{id}", webAuthMiddleware.AuthApi(accountWebApi.HandleUpdateAccount))
	webApisHandler.HandleFunc("DELETE /account/{id}", webAuthMiddleware.AuthApi(accountWebApi.HandleDeleteAccount))

	webApisHandler.HandleFunc("POST /patient", webAuthMiddleware.AuthApi(patientWebApi.HandleCreatePatient))
	webApisHandler.HandleFunc("POST /patient/{id}/blood-test", webAuthMiddleware.AuthApi(patientWebApi.HandleCreatePatientBloodTestResult))
	webApisHandler.HandleFunc("POST /patient/{id}/diagnosis", webAuthMiddleware.AuthApi(patientWebApi.HandleCreatePatientDiagnosisResult))
	webApisHandler.HandleFunc("POST /patient/{id}/checkup", webAuthMiddleware.AuthApi(patientWebApi.HandleCreatePatientCheckUp))
	webApisHandler.HandleFunc("GET /patient/{id}/card", webAuthMiddleware.AuthApi(patientWebApi.HandleGenerateCard))
	webApisHandler.HandleFunc("PUT /patient/{id}/blood-test-result/{btr_id}/pending", webAuthMiddleware.AuthApi(patientWebApi.HandleUpdatePatientPendingBloodTestResult))
	webApisHandler.HandleFunc("PUT /patient/{id}", webAuthMiddleware.AuthApi(patientWebApi.HandleUpdatePatient))
	webApisHandler.HandleFunc("POST /patient/{id}/joints-evaluation", webAuthMiddleware.AuthApi(patientWebApi.HandleCreatePatientJointsEvaluation))
	webApisHandler.HandleFunc("POST /patient/{id}/prophylaxes", webAuthMiddleware.AuthApi(patientWebApi.HandleCreatePatientProphylaxis))
	webApisHandler.HandleFunc("POST /patient/visit/{visit_id}/medicine/{med_id}", webAuthMiddleware.AuthApi(patientWebApi.HandlePatientUseMedicine))
	webApisHandler.HandleFunc("DELETE /patient/{id}", webAuthMiddleware.AuthApi(patientWebApi.HandleDeletePatient))
	webApisHandler.HandleFunc("POST /patients/import/csv", webAuthMiddleware.AuthApi(patientWebApi.HandleUploadImportPatientsFromCsv))

	///
	/// HTMX APIS
	///

	patientHtmx := webhtmx.NewPatientHtmx(usecases)

	htmxHandler := http.NewServeMux()
	htmxHandler.HandleFunc("POST /patient/find", webAuthMiddleware.AuthApi(patientHtmx.HandleFindPatients))
	htmxHandler.HandleFunc("GET /patient/{id}/view", webAuthMiddleware.AuthApi(patientHtmx.HandlePatientDetailsView))
	htmxHandler.HandleFunc("GET /patient/{id}/update", webAuthMiddleware.AuthApi(patientHtmx.HandlePatientUpdateView))

	applicationHandler := http.NewServeMux()
	applicationHandler.Handle("/", version.Handler(appVersion, webi18n.Handler(ismobile.Handler(webtheme.Handler(pagesHandler)))))
	applicationHandler.Handle("/assets/", http.StripPrefix("/assets", static.AssetsHandler(minifyer)))
	applicationHandler.Handle("/api/json/", http.StripPrefix("/api/json", contenttype.Json(v1ApisHandler)))
	applicationHandler.Handle("/api/web/", webi18n.Handler(ismobile.Handler(webtheme.Handler(http.StripPrefix("/api/web", webApisHandler)))))
	applicationHandler.Handle("/htmx/", webi18n.Handler(ismobile.Handler(webtheme.Handler(http.StripPrefix("/htmx", htmxHandler)))))

	log.Info("Starting http server at port " + config.Env().Port)
	switch config.Env().GoEnv {
	case config.GoEnvBeta, config.GoEnvDev, config.GoEnvTest:
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, logger.Handler(applicationHandler)))
	case config.GoEnvProd:
		log.Fatalln(http.ListenAndServe(":"+config.Env().Port, minifyer.Middleware(applicationHandler)))
	}
}
