package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/log"
	"strconv"
)

type patientApi struct {
	usecases *actions.Actions
}

func NewPatientApi(usecases *actions.Actions) *patientApi {
	return &patientApi{
		usecases: usecases,
	}
}

func (e *patientApi) HandleCreatePatient(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreatePatientParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := e.usecases.CreatePatient(reqBody)
	if err != nil {
		log.Errorf("[PATIENT API]: Failed to create patient: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleCreatePatientBloodTestResult(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreatePatientBloodTestResultParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := e.usecases.CreatePatientBloodTestResult(reqBody)
	if err != nil {
		log.Errorf("[PATIENT API]: Failed to create patient's blood test: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleCreatePatientDiagnosisResult(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreatePatientDiagnosisResultParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := e.usecases.CreatePatientDiagnosisResult(reqBody)
	if err != nil {
		log.Errorf("[PATIENT API]: Failed to create patient's diagnosis result: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleListLastPatients(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.ListLastPatients(actions.ListLastPatientsParams{
		ActionContext: ctx,
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleFindPatients(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	findParams := actions.FindPatientsParams{
		ActionContext: ctx,
		PublicId:      r.PathValue("public_id"),
		FirstName:     r.PathValue("first_name"),
		LastName:      r.PathValue("last_name"),
		FatherName:    r.PathValue("father_name"),
		MotherName:    r.PathValue("mother_name"),
		NationalId:    r.PathValue("national_id"),
		PhoneNumber:   r.PathValue("phone_number"),
	}

	payload, err := e.usecases.FindPatients(findParams)
	if err != nil {
		log.Errorf("[PATIENT API]: Failed to find patientes: %+v, error: %s\n", findParams, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleGetPatient(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	params := actions.GetPatientParams{
		ActionContext: ctx,
		PublicId:      r.PathValue("id"),
	}

	payload, err := e.usecases.GetPatient(params)
	if err != nil {
		log.Errorf("[PATIENT API]: Failed to get patient: %+v, error: %s\n", params, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleDeletePatient(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	params := actions.DeletePatientParams{
		ActionContext: ctx,
		PublicId:      r.PathValue("id"),
	}

	payload, err := e.usecases.DeletePatient(params)
	if err != nil {
		log.Errorf("[PATIENT API]: Failed to delete patient: %+v, error: %s\n", params, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleCheckUp(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreatePatientVisitParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx
	reqBody.PatientId = r.PathValue("id")

	payload, err := e.usecases.CreatePatientVisit(reqBody)
	if err != nil {
		log.Errorf("[PATIENT API]: Failed to create patient visit: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleGenerateCard(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.GeneratePatientCard(actions.GeneratePatientCardParams{
		ActionContext: ctx,
		PatientId:     r.PathValue("id"),
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleGetPatientLastVisit(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.GetPatientLastVisit(actions.GetPatientLastVisitParams{
		ActionContext: ctx,
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleListPatientVisits(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.ListPatientVisits(actions.ListPatientVisitsParams{
		ActionContext: ctx,
		PatientId:     r.PathValue("id"),
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleUpdatePendingBloodTestResult(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	btrId, err := strconv.Atoi(r.PathValue("btr_id"))
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var params actions.UpdatePatientPendingBloodTestResultParams
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	params.ActionContext = ctx
	params.PatientPublicId = r.PathValue("id")
	params.BloodTestResultId = uint(btrId)

	payload, err := e.usecases.UpdatePatientPendingBloodTestResult(params)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleCreatePatientJointsEvaluation(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreatePatientJointsEvaluationParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx
	reqBody.PatientId = r.PathValue("id")

	payload, err := e.usecases.CreatePatientJointsEvaluation(reqBody)
	if err != nil {
		log.Errorf("[PATIENT API]: Failed to create patient's joints evaluation: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleListPatientJointsEvaluations(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.ListPatientJointsEvaluations(actions.ListPatientJointsEvaluationsParams{
		ActionContext: ctx,
		PatientId:     r.PathValue("id"),
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleUsePrescribedMedicineForVisit(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	visitId, err := strconv.Atoi(r.PathValue("visit_id"))
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	medId, err := strconv.Atoi(r.PathValue("med_id"))
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.UseMedicineForVisit(actions.UseMedicineForVisitParams{
		ActionContext:        ctx,
		VisitId:              uint(visitId),
		PrescribedMedicineId: uint(medId),
	})
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
