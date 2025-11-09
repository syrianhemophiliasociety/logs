package actions

import "shs/app/models"

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
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
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
	err := checkAccountType(params.Account, models.AccountTypeSuperAdmin)
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
