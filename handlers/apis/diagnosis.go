package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/log"
	"strconv"
)

type diagnosisApi struct {
	usecases *actions.Actions
}

func NewDiagnosisApi(usecases *actions.Actions) *diagnosisApi {
	return &diagnosisApi{
		usecases: usecases,
	}
}

func (e *diagnosisApi) HandleCreateDiagnosis(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreateDiagnosisParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := e.usecases.CreateDiagnosis(reqBody)
	if err != nil {
		log.Errorf("[DIAGNOSIS API]: Failed to create diagnoses: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *diagnosisApi) HandleListDiagnosiss(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.ListAllDiagnoses(actions.ListAllDiagnosesParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Errorf("[DIAGNOSIS API]: Failed to get diagnoses, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *diagnosisApi) HandleDeleteDiagnosis(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.DeleteDiagnosis(actions.DeleteDiagnosisParams{
		ActionContext: ctx,
		DiagnosisId:   uint(id),
	})
	if err != nil {
		log.Errorf("[DIAGNOSIS API]: Failed to delete diagnosis, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
