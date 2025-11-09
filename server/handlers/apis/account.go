package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/log"
)

type accountApi struct {
	usecases *actions.Actions
}

func NewAccountApi(usecases *actions.Actions) *accountApi {
	return &accountApi{
		usecases: usecases,
	}
}

func (e *accountApi) HandleCreateAdminAccount(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreateAdminAccountParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := e.usecases.CreateAdminAccount(reqBody)
	if err != nil {
		log.Errorf("[ACCOUNT API]: Failed to create admin account: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *accountApi) HandleCreateSecritaryAccount(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreateSecritaryAccountParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := e.usecases.CreateSecritaryAccount(reqBody)
	if err != nil {
		log.Errorf("[ACCOUNT API]: Failed to create secritary account: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
