package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/log"
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

func (e *patientApi) HandleCreatePatientBloodTest(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreatePatientBloodTestParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := e.usecases.CreatePatientBloodTest(reqBody)
	if err != nil {
		log.Errorf("[PATIENT API]: Failed to create patient's blood test: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *patientApi) HandleCreatePatientVirus(w http.ResponseWriter, r *http.Request) {
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

func (e *patientApi) HandleFindPatients(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	findParams := actions.FindPatientsParams{
		ActionContext: ctx,
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
