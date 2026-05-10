package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/log"
	"strconv"
)

type bloodTestApi struct {
	usecases *actions.Actions
}

func NewBloodTestApi(usecases *actions.Actions) *bloodTestApi {
	return &bloodTestApi{
		usecases: usecases,
	}
}

func (e *bloodTestApi) HandleCreateBloodTest(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreateBloodTestParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := e.usecases.CreateBloodTest(reqBody)
	if err != nil {
		log.Errorf("[BLOODTEST API]: Failed to create blood test: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *bloodTestApi) HandleGetBloodTest(w http.ResponseWriter, r *http.Request) {
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

	payload, err := e.usecases.GetBloodTest(actions.GetBloodTestParams{
		ActionContext: ctx,
		BloodTestId:   uint(id),
	})
	if err != nil {
		log.Errorf("[BLOODTEST API]: Failed to get blood test, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *bloodTestApi) HandleListBloodTests(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.ListAllBloodTests(actions.ListAllBloodTestsParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Errorf("[BLOODTEST API]: Failed to get blood tests, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *bloodTestApi) HandleDeleteBloodTest(w http.ResponseWriter, r *http.Request) {
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

	payload, err := e.usecases.DeleteBloodTest(actions.DeleteBloodTestParams{
		ActionContext: ctx,
		BloodTestId:   uint(id),
	})
	if err != nil {
		log.Errorf("[BLOODTEST API]: Failed to delete blood test, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
