package actions

import "shs/app/models"

type Account struct {
	Id          uint   `json:"id"`
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	Type        string `json:"type"`
}

type createAccountParams struct {
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type CreateSecritaryAccountParams struct {
	ActionContext
	NewAccount createAccountParams `json:"new_account"`
}

type CreateSecritaryAccountPayload struct {
}

func (a *Actions) CreateSecritaryAccount(params CreateSecritaryAccountParams) (CreateSecritaryAccountPayload, error) {
	err := params.Account.CheckType(models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
	if err != nil {
		return CreateSecritaryAccountPayload{}, err
	}

	_, err = a.app.CreateAccount(models.Account{
		DisplayName: params.NewAccount.DisplayName,
		Username:    params.NewAccount.Username,
		Password:    params.NewAccount.Password,
		Type:        models.AccountTypeSecritary,
	})

	return CreateSecritaryAccountPayload{}, err
}

type CreateAdminAccountParams struct {
	ActionContext
	NewAccount createAccountParams `json:"new_account"`
}

type CreateAdminAccountPayload struct {
}

func (a *Actions) CreateAdminAccount(params CreateAdminAccountParams) (CreateAdminAccountPayload, error) {
	err := params.Account.CheckType(models.AccountTypeSuperAdmin)
	if err != nil {
		return CreateAdminAccountPayload{}, err
	}

	_, err = a.app.CreateAccount(models.Account{
		DisplayName: params.NewAccount.DisplayName,
		Username:    params.NewAccount.Username,
		Password:    params.NewAccount.Password,
		Type:        models.AccountTypeAdmin,
	})

	return CreateAdminAccountPayload{}, err
}

type DeleteAccountParams struct {
	ActionContext
	AccountId uint
}

type DeleteAccountPayload struct {
}

func (a *Actions) DeleteAccount(params DeleteAccountParams) (DeleteAccountPayload, error) {
	err := params.Account.CheckType(models.AccountTypeAdmin)
	if err != nil {
		return DeleteAccountPayload{}, err
	}

	err = a.app.DeleteAccount(params.AccountId)
	if err != nil {
		return DeleteAccountPayload{}, err
	}

	return DeleteAccountPayload{}, nil
}

type ListAllAccountsParams struct {
	ActionContext
}

type ListAllAccountsPayload struct {
	Data []Account `json:"data"`
}

func (a *Actions) ListAllAccounts(params ListAllAccountsParams) (ListAllAccountsPayload, error) {
	err := params.Account.CheckType(models.AccountTypeAdmin)
	if err != nil {
		return ListAllAccountsPayload{}, err
	}

	var accounts []models.Account
	switch params.Account.Type {
	case models.AccountTypeAdmin:
		accounts, err = a.app.ListAllAccountsForAdmin()
	case models.AccountTypeSuperAdmin:
		accounts, err = a.app.ListAllAccountsForSuperAdmin()
	}
	if err != nil {
		return ListAllAccountsPayload{}, err
	}

	outAccounts := make([]Account, 0, len(accounts))
	for _, account := range accounts {
		outAccounts = append(outAccounts, Account{
			Id:          account.Id,
			DisplayName: account.DisplayName,
			Username:    account.Username,
			Type:        string(account.Type),
		})
	}

	return ListAllAccountsPayload{
		Data: outAccounts,
	}, nil
}
