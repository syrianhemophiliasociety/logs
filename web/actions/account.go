package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shs-web/errors"
	"strconv"

	goerrors "errors"
)

type AccountPermissions uint64

const (
	AccountPermissionReadAccounts AccountPermissions = 1 << iota
	AccountPermissionWriteAccounts
	AccountPermissionReadPatient
	AccountPermissionWritePatient
	AccountPermissionReadMedicine
	AccountPermissionWriteMedicine
	AccountPermissionReadVirus
	AccountPermissionWriteVirus
	AccountPermissionReadBloodTest
	AccountPermissionWriteBloodTest
	AccountPermissionReadOwnVisit
	AccountPermissionWriteOwnVisit
	AccountPermissionReadOtherVisits
	AccountPermissionWriteOtherVisits
	AccountPermissionReadDiagnoses
	AccountPermissionWriteDiagnoses
)

type Account struct {
	Id          uint               `json:"id"`
	DisplayName string             `json:"display_name"`
	Username    string             `json:"username"`
	Password    string             `json:"password"`
	Type        string             `json:"type"`
	Permissions AccountPermissions `json:"permissions"`
}

func (a Account) HasPermission(p AccountPermissions) bool {
	return a.Permissions&p != 0
}

type CreateAccountParams struct {
	RequestContext
	NewAccount Account `json:"new_account"`
}

type CreateAccountPayload struct {
	Id uint `json:"id"`
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

type GetAccountParams struct {
	RequestContext
	AccountId uint
}

type GetAccountPayload struct {
	Data Account `json:"data"`
}

func (a *Actions) GetAccount(params GetAccountParams) (Account, error) {
	payload, err := makeRequest[any, GetAccountPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/v1/account/%d", params.AccountId),
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
	if err != nil {
		return Account{}, err
	}

	return payload.Data, nil
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

type UpdateAccountParams struct {
	RequestContext
	AccountId  uint
	NewAccount Account `json:"new_account"`
}

type UpdateAccountPayload struct {
}

type UpdateAccountRequest struct {
	Account
}

func (a *UpdateAccountRequest) UnmarshalJSON(payload []byte) error {
	var data map[string]any
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}

	var ok bool
	(*a).DisplayName, ok = data["display_name"].(string)
	if !ok {
		return goerrors.New("invalid display_name value")
	}
	(*a).Username, ok = data["username"].(string)
	if !ok {
		return goerrors.New("invalid username value")
	}
	(*a).Password, ok = data["password"].(string)
	if !ok {
		return goerrors.New("invalid password value")
	}

	const permissionsKey = "permissions"
	switch data[permissionsKey].(type) {
	case string:
		p, err := strconv.Atoi(data[permissionsKey].(string))
		if err != nil {
			return err
		}
		if (p & (p - 1)) != 0 {
			return goerrors.New("invalid permissions value")
		}
		(*a).Permissions = AccountPermissions(p)

	case []any:
		for _, p := range data[permissionsKey].([]any) {
			pStr, ok := p.(string)
			if !ok {
				return goerrors.New("invalid permissions type")
			}
			pInt, err := strconv.Atoi(pStr)
			if err != nil {
				return err
			}
			if (pInt & (pInt - 1)) != 0 {
				return goerrors.New("invalid permissions value")
			}
			(*a).Permissions |= AccountPermissions(pInt)
		}

	default:
		return goerrors.New("invalid permissions value")
	}

	return nil
}

func (a *Actions) UpdateAccount(params UpdateAccountParams) (UpdateAccountPayload, error) {
	payload, err := makeRequest[map[string]any, UpdateAccountPayload](makeRequestConfig[map[string]any]{
		method:   http.MethodPut,
		endpoint: fmt.Sprintf("/v1/account/%d", params.AccountId),
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: map[string]any{
			"new_account": params.NewAccount,
		},
	})
	if err != nil {
		return UpdateAccountPayload{}, err
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
