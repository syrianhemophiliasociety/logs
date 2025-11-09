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
		log.Errorf("[PATIENT API]: Failed to find patientes: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
