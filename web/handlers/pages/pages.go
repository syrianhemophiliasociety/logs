package pages

import (
	"net/http"
	"shs-web/actions"
	"shs-web/config"
	"shs-web/handlers/middlewares/contenttype"
	"shs-web/i18n"
	"shs-web/views/components"
	"shs-web/views/layouts"
	"shs-web/views/pages"
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavHome)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavAbout)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavPrivacy)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavLogin)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavViruses)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavMedicine)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavPatient)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavBloodTests)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavManagement)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavManagement)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavPatients)
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
	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavPatient)
		w.Header().Set("HX-Push-Url", "/patient/"+id)
		pages.Patient(patient, bloodTests, viruses, allMedicine, visits).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Patient(patient, bloodTests, viruses, allMedicine, visits)).Render(r.Context(), w)
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
		w.Header().Set("HX-Title", i18n.StringsCtx(r.Context()).NavPatient)
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
