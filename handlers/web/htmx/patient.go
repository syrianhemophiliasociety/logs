package htmx

import (
	"encoding/json"
	"errors"
	"net/http"
	"shs/actions"
	"shs/app"
	"shs/handlers/web/context"
	"shs/log"
	"shs/web/i18n"
	"shs/web/views/components"
)

type patientHtmx struct {
	usecases *actions.Actions
}

func NewPatientHtmx(usecases *actions.Actions) *patientHtmx {
	return &patientHtmx{
		usecases: usecases,
	}
}

func (p *patientHtmx) HandleFindPatients(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody actions.FindPatientsParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := p.usecases.FindPatients(reqBody)
	if errors.Is(err, app.ErrNotFound{}) {
		components.NotFoundError(i18n.StringsCtx(r.Context()).NavPatients).Render(r.Context(), w)
		return
	}
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	components.PatientsBrief(payload.Data).Render(r.Context(), w)
}

func (p *patientHtmx) HandlePatientUpdateView(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")
	patient, err := p.usecases.GetPatient(actions.GetPatientParams{
		ActionContext: ctx,
		PublicId:      patientId,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	_ = components.PatientUpdateProfile(patient.Data).
		Render(r.Context(), w)
}

func (p *patientHtmx) HandlePatientDetailsView(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")
	patient, err := p.usecases.GetPatient(actions.GetPatientParams{
		ActionContext: ctx,
		PublicId:      patientId,
	})
	if err != nil {
		components.GenericError("Something went wrong").
			Render(r.Context(), w)
		return
	}

	_ = components.PatientViewProfile(patient.Data).
		Render(r.Context(), w)
}
