package pages

import (
	"fmt"
	"net/http"
	"shs-web/actions"
	"shs-web/config"
	"shs-web/handlers/middlewares/contenttype"
	"shs-web/i18n"
	"shs-web/views/components"
	"shs-web/views/layouts"
	"shs-web/views/pages"
	"slices"
	"strconv"

	_ "github.com/a-h/templ"
)

type pagesHandler struct {
	usecases *actions.Actions
}

func New(usecases *actions.Actions) *pagesHandler {
	return &pagesHandler{
		usecases: usecases,
	}
}

func (p *pagesHandler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavHome)
		w.Header().Set("HX-Push-Url", "/")
		pages.Index().Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavHome,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Index()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavAbout)
		w.Header().Set("HX-Push-Url", "/about")
		pages.About().Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavAbout,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.About()).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePrivacyPage(w http.ResponseWriter, r *http.Request) {
	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPrivacy)
		w.Header().Set("HX-Push-Url", "/privacy")
		pages.Privacy().Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPrivacy,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Privacy()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavLogin)
		w.Header().Set("HX-Push-Url", "/login")
		pages.Login().Render(r.Context(), w)
		return
	}

	layouts.Raw(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavLogin,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Login()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleVirusesPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	viruses, err := p.usecases.ListAllViruses(actions.ListAllVirusesParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTests, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavViruses)
		w.Header().Set("HX-Push-Url", "/viruses")
		pages.Viruses(viruses, bloodTests).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavViruses,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Viruses(viruses, bloodTests)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleMedicinesPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	medicines, err := p.usecases.ListAllMedicines(actions.ListAllMedicinesParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavMedicine)
		w.Header().Set("HX-Push-Url", "/medicines")
		pages.Medicines(medicines).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavMedicine,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Medicines(medicines)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleMedicinePage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	id := r.PathValue("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	medicine, err := p.usecases.GetMedicine(actions.GetMedicineParams{
		RequestContext: ctx,
		MedicineId:     uint(intId),
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", "/medicine/"+id)
		pages.Medicine(medicine).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Medicine(medicine)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleBloodTestsPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	bloodTests, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavBloodTests)
		w.Header().Set("HX-Push-Url", "/blood-tests")
		pages.BloodTests(bloodTests).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavBloodTests,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.BloodTests(bloodTests)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleBloodTestPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	bloodTests, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTestIndex := slices.IndexFunc(bloodTests, func(bt actions.BloodTest) bool {
		return bt.Id == uint(id)
	})
	if bloodTestIndex < 0 {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavBloodTests)
		w.Header().Set("HX-Push-Url", "/blood-test/"+strconv.Itoa(id))
		pages.BloodTest(bloodTests[bloodTestIndex]).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavBloodTests,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.BloodTest(bloodTests[bloodTestIndex])).Render(r.Context(), w)
}

func (p *pagesHandler) HandleManagementPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	accounts, err := p.usecases.ListAllAccounts(actions.ListAllAccountsParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavManagement)
		w.Header().Set("HX-Push-Url", "/management")
		pages.Management(accounts).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavManagement,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Management(accounts)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleAccountManagementPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	account, err := p.usecases.GetAccount(actions.GetAccountParams{
		RequestContext: ctx,
		AccountId:      uint(id),
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavManagement)
		w.Header().Set("HX-Push-Url", "/management/account/"+strconv.Itoa(int(account.Id)))
		pages.Account(account).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavManagement,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Account(account)).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePatientsPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	bloodTests, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	viruses, err := p.usecases.ListAllViruses(actions.ListAllVirusesParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	lastPatients, err := p.usecases.ListLastPatients(actions.ListLastPatientsParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatients)
		w.Header().Set("HX-Push-Url", "/patients")
		pages.Patients(bloodTests, viruses, lastPatients).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatients,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Patients(bloodTests, viruses, lastPatients)).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePatientPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	id := r.PathValue("id")

	patient, err := p.usecases.GetPatient(actions.GetPatientParams{
		RequestContext: ctx,
		PatientId:      id,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTests, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	viruses, err := p.usecases.ListAllViruses(actions.ListAllVirusesParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	allMedicine, err := p.usecases.ListAllMedicines(actions.ListAllMedicinesParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	visits, err := p.usecases.ListPatientVisits(actions.ListPatientVisitsParams{
		RequestContext: ctx,
		PatientId:      patient.PublicId,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	diagnoses, err := p.usecases.ListAllDiagnoses(actions.ListAllDiagnosesParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", "/patient/"+id)
		pages.Patient(patient, bloodTests, viruses, allMedicine, visits, diagnoses).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Patient(patient, bloodTests, viruses, allMedicine, visits, diagnoses)).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePatientBloodTestResultPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	id := r.PathValue("id")
	btrId := r.PathValue("btr_id")

	patient, err := p.usecases.GetPatient(actions.GetPatientParams{
		RequestContext: ctx,
		PatientId:      id,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTestResultIndex := slices.IndexFunc(patient.BloodTests, func(btr actions.BloodTestResult) bool {
		return strconv.Itoa(int(btr.Id)) == btrId
	})
	if bloodTestResultIndex < 0 {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTests, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTestIndex := slices.IndexFunc(bloodTests, func(bt actions.BloodTest) bool {
		return bt.Id == patient.BloodTests[bloodTestResultIndex].BloodTestId
	})
	if bloodTestIndex < 0 {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", fmt.Sprintf("/patient/%s/blood-test-result/%s", patient.PublicId, btrId))
		pages.PatientBloodTestResult(patient, patient.BloodTests[bloodTestResultIndex], bloodTests[bloodTestIndex]).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.PatientBloodTestResult(patient, patient.BloodTests[bloodTestResultIndex], bloodTests[bloodTestIndex])).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePatientVisitPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	id := r.PathValue("id")
	visitId, err := strconv.Atoi(r.PathValue("visit_id"))
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	patient, err := p.usecases.GetPatient(actions.GetPatientParams{
		RequestContext: ctx,
		PatientId:      id,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	visits, err := p.usecases.ListPatientVisits(actions.ListPatientVisitsParams{
		RequestContext: ctx,
		PatientId:      id,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	visitIndex := slices.IndexFunc(visits, func(v actions.Visit) bool {
		return v.Id == uint(visitId)
	})
	if visitIndex < 0 {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", fmt.Sprintf("/patient/%s/visit/%d", patient.PublicId, visitId))
		pages.PatientVisit(patient, visits[visitIndex]).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.PatientVisit(patient, visits[visitIndex])).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePatientMedicationsPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	payload, err := p.usecases.GetPatientLastVisit(actions.GetPatientLastVisitParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", "/patient/medications")
		pages.PatientMedicine(payload).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.PatientMedicine(payload)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleDiagnosesPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	diagnoses, err := p.usecases.ListAllDiagnoses(actions.ListAllDiagnosesParams{
		RequestContext: ctx,
	})
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", "/diagnoses")
		pages.Diagnoses(diagnoses).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Diagnoses(diagnoses)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleStatisticsPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	_ = ctx

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", "/statistics")
		pages.Statistics().Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Statistics()).Render(r.Context(), w)
}
