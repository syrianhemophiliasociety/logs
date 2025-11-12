package actions

import (
	"fmt"
	"net/http"
	"shs-web/errors"
)

type Account struct {
	Id          uint   `json:"id"`
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Type        string `json:"type"`
}

type CreateAccountParams struct {
	RequestContext
	NewAccount Account `json:"new_account"`
}

type CreateAccountPayload struct {
}

func (a *Actions) CreateAccount(params CreateAccountParams) (CreateAccountPayload, error) {
	endpoint := ""
	switch params.NewAccount.Type {
	case "secritary":
		endpoint = "/v1/account/secritary"
	case "admin":
		endpoint = "/v1/account/admin"
	default:
		return CreateAccountPayload{}, errors.ErrSomethingWentWrong
	}

	payload, err := makeRequest[CreateAccountParams, CreateAccountPayload](makeRequestConfig[CreateAccountParams]{
		method:   http.MethodPost,
		endpoint: endpoint,
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: params,
	})
	if err != nil {
		return CreateAccountPayload{}, err
	}

	return payload, nil
}

type DeleteAccountParams struct {
	RequestContext
	AccountId uint
}

type DeleteAccountPayload struct {
}

func (a *Actions) DeleteAccount(params DeleteAccountParams) (DeleteAccountPayload, error) {
	payload, err := makeRequest[DeleteAccountParams, DeleteAccountPayload](makeRequestConfig[DeleteAccountParams]{
		method:   http.MethodDelete,
		endpoint: fmt.Sprintf("/v1/account/%d", params.AccountId),
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: params,
	})
	if err != nil {
		return DeleteAccountPayload{}, err
	}

	return payload, nil
}

type ListAllAccountsParams struct {
	RequestContext
}

type ListAllAccountsPayload struct {
	Data []Account `json:"data"`
}

func (a *Actions) ListAllAccounts(params ListAllAccountsParams) ([]Account, error) {
	payload, err := makeRequest[any, ListAllAccountsPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/account/all",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
	if err != nil {
		return nil, err
	}

	return payload.Data, nil
}
