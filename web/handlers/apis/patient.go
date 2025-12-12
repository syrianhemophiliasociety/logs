package apis

import (
	"encoding/json"
	"net/http"
	"shs-web/actions"
	"shs-web/errors"
	"shs-web/i18n"
	"shs-web/log"
	"shs-web/views/components"
)

type patientApi struct {
	usecases *actions.Actions
}

func NewPatientApi(usecases *actions.Actions) *patientApi {
	return &patientApi{
		usecases: usecases,
	}
}

func (v *patientApi) HandleCreatePatient(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody actions.PatientRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	payload, err := v.usecases.CreatePatient(actions.CreatePatientParams{
		RequestContext: ctx,
		NewPatient:     reqBody,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	w.Header().Set("HX-Redirect", "/patient/"+payload.Id)
}

func (v *patientApi) HandleFindPatients(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
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
	reqBody.RequestContext = ctx

	payload, err := v.usecases.FindPatients(reqBody)
	if errors.Is(err, errors.ErrPatientNotFound) {
		components.NotFoundError(i18n.StringsCtx(r.Context()).NavPatients).Render(r.Context(), w)
		return
	}
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	components.PatientsBrief(payload).Render(r.Context(), w)
}

func (v *patientApi) HandleAddPatientBloodTest(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")

	var reqBody actions.PatientBloodTests
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.CreatePatientBloodTest(actions.CreatePatientBloodTestParams{
		RequestContext:   ctx,
		PatientId:        patientId,
		PatientBloodTest: reqBody.BloodTests[0],
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *patientApi) HandleCreatePatientCheckUp(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")

	var reqBody actions.CreateCheckUpRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		writeRawTextResponse(w, i18n.Strings("en").ErrorSomethingWentWrong)
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.CreatePatientCheckUp(actions.CreatePatientCheckUpParams{
		RequestContext: ctx,
		PatientId:      patientId,
		CheckUpRequest: reqBody,
	})
	if errors.Is(err, errors.ErrInsufficientMedicineAmount{}) {
		imErr := err.(errors.ErrInsufficientMedicineAmount)
		writeRawTextResponse(w, i18n.Strings("en").ErrorInsufficientMedicineAmountFmt(imErr.MedicineName, imErr.ExceedingAmount, imErr.LeftPackages))
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorInsufficientMedicineAmountFmt(imErr.MedicineName, imErr.ExceedingAmount, imErr.LeftPackages)).Render(r.Context(), w)
		log.Errorln(err)
		return
	}
	if err != nil {
		writeRawTextResponse(w, i18n.Strings("en").ErrorSomethingWentWrong)
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *patientApi) HandleGenerateCard(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")

	payload, err := v.usecases.GeneratePatientCard(actions.GeneratePatientCardParams{
		RequestContext: ctx,
		PatientId:      patientId,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	w.Write([]byte(payload.ImageBase64))
}

func (v *patientApi) HandleDeletePatient(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")

	_, err = v.usecases.DeletePatient(actions.DeletePatientParams{
		RequestContext: ctx,
		PatientId:      patientId,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *patientApi) HandleUpdatePatientPendingBloodTestResult(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")
	btrId := r.PathValue("btr_id")

	var reqBody actions.PatientBloodTests
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.UpdatePatientPendingBloodTest(actions.UpdatePatientPendingBloodTestParams{
		RequestContext:    ctx,
		PatientId:         patientId,
		BloodTestResultId: btrId,
		FilledFields:      reqBody.BloodTests[0].FilledFields,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}
