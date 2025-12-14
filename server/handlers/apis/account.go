package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/log"
	"strconv"
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

func (e *accountApi) HandleListAllAccounts(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.ListAllAccounts(actions.ListAllAccountsParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Errorf("[ACCOUNT API]: Failed to get accounts, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *accountApi) HandleGetAccount(w http.ResponseWriter, r *http.Request) {
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

	payload, err := e.usecases.GetAccount(actions.GetAccountParams{
		ActionContext: ctx,
		AccountId:     uint(id),
	})
	if err != nil {
		log.Errorf("[ACCOUNT API]: Failed to get account, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *accountApi) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) {
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

	payload, err := e.usecases.DeleteAccount(actions.DeleteAccountParams{
		ActionContext: ctx,
		AccountId:     uint(id),
	})
	if err != nil {
		log.Errorf("[ACCOUNT API]: Failed to delete account, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *accountApi) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) {
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

	log.Debugln(id)

	var params actions.UpdateAccountParams
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.UpdateAccount(actions.UpdateAccountParams{
		ActionContext: ctx,
		AccountId:     uint(id),
		NewAccount:    params.NewAccount,
	})
	if err != nil {
		log.Errorf("[ACCOUNT API]: Failed to update account, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
