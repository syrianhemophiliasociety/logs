package actions

import "shs/app/models"

const (
	patientPermissions = models.AccountPermissionReadOwnVisit | models.AccountPermissionWriteOwnVisit

	secritaryPermissions = models.AccountPermissionReadPatient | models.AccountPermissionWritePatient |
		models.AccountPermissionReadMedicine | models.AccountPermissionWriteMedicine |
		models.AccountPermissionReadOtherVisits | models.AccountPermissionWriteOtherVisits |
		models.AccountPermissionReadBloodTest

	adminPermissions = secritaryPermissions |
		models.AccountPermissionReadAccounts | models.AccountPermissionWriteAccounts |
		models.AccountPermissionReadBloodTest | models.AccountPermissionWriteBloodTest |
		models.AccountPermissionReadMedicine | models.AccountPermissionWriteBloodTest |
		models.AccountPermissionReadVirus | models.AccountPermissionWriteVirus
)

type Account struct {
	Id          uint                      `json:"id"`
	DisplayName string                    `json:"display_name"`
	Username    string                    `json:"username"`
	Type        string                    `json:"type"`
	Permissions models.AccountPermissions `json:"permissions"`
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
	if !params.Account.HasPermission(models.AccountPermissionWriteAccounts) {
		return CreateSecritaryAccountPayload{}, ErrPermissionDenied{}
	}

	_, err := a.app.CreateAccount(models.Account{
		DisplayName: params.NewAccount.DisplayName,
		Username:    params.NewAccount.Username,
		Password:    params.NewAccount.Password,
		Type:        models.AccountTypeSecritary,
		Permissions: secritaryPermissions,
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
	if !params.Account.HasPermission(models.AccountPermissionWriteAdmins) {
		return CreateAdminAccountPayload{}, ErrPermissionDenied{}
	}

	_, err := a.app.CreateAccount(models.Account{
		DisplayName: params.NewAccount.DisplayName,
		Username:    params.NewAccount.Username,
		Password:    params.NewAccount.Password,
		Type:        models.AccountTypeAdmin,
		Permissions: adminPermissions,
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
	account, err := a.app.GetAccountById(params.AccountId)
	if err != nil {
		return DeleteAccountPayload{}, err
	}

	if account.Type == models.AccountTypeAdmin && !params.Account.HasPermission(models.AccountPermissionWriteAdmins) {
		return DeleteAccountPayload{}, ErrPermissionDenied{}
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
	if !params.Account.HasPermission(models.AccountPermissionReadAccounts) {
		return ListAllAccountsPayload{}, ErrPermissionDenied{}
	}

	var accounts []models.Account
	var err error
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
			Permissions: account.Permissions,
		})
	}

	return ListAllAccountsPayload{
		Data: outAccounts,
	}, nil
}
