package pages

import (
	"fmt"
	"net/http"
	"shs/actions"
	"shs/app/models"
	"shs/config"
	"shs/handlers/middlewares/contenttype"
	"shs/handlers/web/context"
	"shs/web/i18n"
	"shs/web/views/components"
	"shs/web/views/layouts"
	"shs/web/views/pages"
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
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	viruses, err := p.usecases.ListAllViruses(actions.ListAllVirusesParams{
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTests, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavViruses)
		w.Header().Set("HX-Push-Url", "/viruses")
		pages.Viruses(viruses.Data, bloodTests.Data).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavViruses,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Viruses(viruses.Data, bloodTests.Data)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleMedicinesPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	medicines, err := p.usecases.ListAllMedicine(actions.ListAllMedicineParams{
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavMedicine)
		w.Header().Set("HX-Push-Url", "/medicines")
		pages.Medicines(medicines.Data).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavMedicine,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Medicines(medicines.Data)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleMedicinePage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
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
		ActionContext: ctx,
		MedicineId:    uint(intId),
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavMedicine)
		w.Header().Set("HX-Push-Url", "/medicine/"+id)
		pages.Medicine(medicine.Data).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavMedicine,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Medicine(medicine.Data)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleBloodTestsPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	bloodTests, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavBloodTests)
		w.Header().Set("HX-Push-Url", "/blood-tests")
		pages.BloodTests(bloodTests.Data).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavBloodTests,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.BloodTests(bloodTests.Data)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleBloodTestPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
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
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTestIndex := slices.IndexFunc(bloodTests.Data, func(bt actions.BloodTest) bool {
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
		pages.BloodTest(bloodTests.Data[bloodTestIndex]).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavBloodTests,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.BloodTest(bloodTests.Data[bloodTestIndex])).Render(r.Context(), w)
}

func (p *pagesHandler) HandleManagementPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	accounts, err := p.usecases.ListAllAccounts(actions.ListAllAccountsParams{
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavManagement)
		w.Header().Set("HX-Push-Url", "/management")
		pages.Management(accounts.Data).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavManagement,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Management(accounts.Data)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleAccountManagementPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
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
		ActionContext: ctx,
		AccountId:     uint(id),
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavManagement)
		w.Header().Set("HX-Push-Url", "/management/account/"+strconv.Itoa(int(account.Account.Id)))
		pages.Account(account.Account).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavManagement,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Account(account.Account)).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePatientsPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	bloodTests, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	viruses, err := p.usecases.ListAllViruses(actions.ListAllVirusesParams{
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	lastPatients, err := p.usecases.ListLastPatients(actions.ListLastPatientsParams{
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatients)
		w.Header().Set("HX-Push-Url", "/patients")
		pages.Patients(bloodTests.Data, viruses.Data, lastPatients.Data).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatients,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Patients(bloodTests.Data, viruses.Data, lastPatients.Data)).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePatientPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	id := r.PathValue("id")

	patient, err := p.usecases.GetPatient(actions.GetPatientParams{
		ActionContext: ctx,
		PublicId:      id,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	var (
		bloodTests  []actions.BloodTest
		viruses     []actions.Virus
		allMedicine []actions.Medicine
		visits      []actions.Visit
		diagnoses   []actions.Diagnosis
	)

	if ctx.Account.HasPermission(models.AccountPermissionReadBloodTest) {
		bloodTestsPL, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
			ActionContext: ctx,
		})
		if err != nil {
			components.GenericError("Something went wrong").
				Render(r.Context(), w)
			return
		}
		bloodTests = bloodTestsPL.Data
	}

	if ctx.Account.HasPermission(models.AccountPermissionReadVirus) {
		virusesPL, err := p.usecases.ListAllViruses(actions.ListAllVirusesParams{
			ActionContext: ctx,
		})
		if err != nil {
			components.GenericError("Something went wrong").
				Render(r.Context(), w)
			return
		}
		viruses = virusesPL.Data
	}

	if ctx.Account.HasPermission(models.AccountPermissionReadMedicine) {
		allMedicinePL, err := p.usecases.ListAllMedicine(actions.ListAllMedicineParams{
			ActionContext: ctx,
		})
		if err != nil {
			components.GenericError("Something went wrong").
				Render(r.Context(), w)
			return
		}
		allMedicine = allMedicinePL.Data
	}

	if ctx.Account.HasPermission(models.AccountPermissionReadOtherVisits) {
		visitsPL, err := p.usecases.ListPatientVisits(actions.ListPatientVisitsParams{
			ActionContext: ctx,
			PatientId:     patient.Data.PublicId,
		})
		if err != nil {
			components.GenericError("Something went wrong").
				Render(r.Context(), w)
			return
		}
		visits = visitsPL.Data
	}

	if ctx.Account.HasPermission(models.AccountPermissionReadDiagnoses) {
		diagnosesPL, err := p.usecases.ListAllDiagnoses(actions.ListAllDiagnosesParams{
			ActionContext: ctx,
		})
		if err != nil {
			components.GenericError("What do you think you're doing?").
				Render(r.Context(), w)
			return
		}
		diagnoses = diagnosesPL.Data
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", "/patient/"+id)
		pages.Patient(patient.Data, bloodTests, viruses, allMedicine, visits, diagnoses).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Patient(patient.Data, bloodTests, viruses, allMedicine, visits, diagnoses)).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePatientBloodTestResultPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	id := r.PathValue("id")
	btrId := r.PathValue("btr_id")

	patient, err := p.usecases.GetPatient(actions.GetPatientParams{
		ActionContext: ctx,
		PublicId:      id,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTestResultIndex := slices.IndexFunc(patient.Data.BloodTestResults, func(btr actions.BloodTestResult) bool {
		return strconv.Itoa(int(btr.Id)) == btrId
	})
	if bloodTestResultIndex < 0 {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTests, err := p.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	bloodTestIndex := slices.IndexFunc(bloodTests.Data, func(bt actions.BloodTest) bool {
		return bt.Id == patient.Data.BloodTestResults[bloodTestResultIndex].BloodTestId
	})
	if bloodTestIndex < 0 {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", fmt.Sprintf("/patient/%s/blood-test-result/%s", patient.Data.PublicId, btrId))
		pages.PatientBloodTestResult(patient.Data, patient.Data.BloodTestResults[bloodTestResultIndex], bloodTests.Data[bloodTestIndex]).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.PatientBloodTestResult(patient.Data, patient.Data.BloodTestResults[bloodTestResultIndex], bloodTests.Data[bloodTestIndex])).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePatientVisitPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
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
		ActionContext: ctx,
		PublicId:      id,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	visits, err := p.usecases.ListPatientVisits(actions.ListPatientVisitsParams{
		ActionContext: ctx,
		PatientId:     id,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	visitIndex := slices.IndexFunc(visits.Data, func(v actions.Visit) bool {
		return v.Id == uint(visitId)
	})
	if visitIndex < 0 {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", fmt.Sprintf("/patient/%s/visit/%d", patient.Data.PublicId, visitId))
		pages.PatientVisit(patient.Data, visits.Data[visitIndex]).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.PatientVisit(patient.Data, visits.Data[visitIndex])).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePatientMedicationsPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	payload, err := p.usecases.GetPatientLastVisit(actions.GetPatientLastVisitParams{
		ActionContext: ctx,
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
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	diagnoses, err := p.usecases.ListAllDiagnoses(actions.ListAllDiagnosesParams{
		ActionContext: ctx,
	})
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavPatient)
		w.Header().Set("HX-Push-Url", "/diagnoses")
		pages.Diagnoses(diagnoses.Data).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavPatient,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Diagnoses(diagnoses.Data)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleStatisticsPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
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

func (p *pagesHandler) HandleMedicinesUseLogsPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	_ = ctx

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavMedicineUseLogs)
		w.Header().Set("HX-Push-Url", "/medicines/logs")
		pages.MedicinesUseLogs().Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavMedicineUseLogs,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.MedicinesUseLogs()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleVisitsPage(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	payload, err := p.usecases.ListAllVisits(actions.ListAllVisitsParams{ActionContext: ctx})
	if err != nil {
		components.GenericError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", i18n.Strings("en").NavVisits)
		w.Header().Set("HX-Push-Url", "/visits")
		pages.Visits(payload.Data).Render(r.Context(), w)
		return
	}

	layouts.Default(layouts.PageProps{
		Title:    i18n.StringsCtx(r.Context()).NavVisits,
		Url:      config.Env().Hostname,
		ImageUrl: config.Env().Hostname + "/assets/favicon-32x32.png",
	}, pages.Visits(payload.Data)).Render(r.Context(), w)
}
