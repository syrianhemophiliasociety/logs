package actions

import (
	"shs/app/models"
)

const (
	patientPermissions = models.AccountPermissionReadOwnVisit | models.AccountPermissionWriteOwnVisit

	secritaryPermissions = models.AccountPermissionReadPatient | models.AccountPermissionWritePatient |
		models.AccountPermissionReadMedicine | models.AccountPermissionWriteMedicine |
		models.AccountPermissionReadOtherVisits | models.AccountPermissionWriteOtherVisits |
		models.AccountPermissionReadBloodTest |
		models.AccountPermissionReadVirus |
		models.AccountPermissionReadDiagnoses

	adminPermissions = secritaryPermissions |
		models.AccountPermissionReadAccounts | models.AccountPermissionWriteAccounts |
		models.AccountPermissionReadBloodTest | models.AccountPermissionWriteBloodTest |
		models.AccountPermissionReadMedicine | models.AccountPermissionWriteBloodTest |
		models.AccountPermissionReadVirus | models.AccountPermissionWriteVirus |
		models.AccountPermissionReadDiagnoses | models.AccountPermissionWriteDiagnoses |
		models.AccountPermissionReadJoints | models.AccountPermissionWriteJoints

	// aka jointologist
	snoopDoggPermissions = models.AccountPermissionReadPatient |
		models.AccountPermissionReadJoints | models.AccountPermissionWriteJoints
)

type Account struct {
	Id          uint                      `json:"id"`
	DisplayName string                    `json:"display_name"`
	Username    string                    `json:"username"`
	Type        string                    `json:"type"`
	Password    string                    `json:"password,omitempty"`
	Permissions models.AccountPermissions `json:"permissions"`
}

func (a Account) HasPermission(p models.AccountPermissions) bool {
	return a.Permissions&p != 0
}

func (a Account) Validate() error {
	if a.Username == "" {
		return ErrInvalidAccountUsername{}
	}
	if a.Password == "" {
		return ErrInvalidAccountPassword{}
	}
	if a.DisplayName == "" {
		return ErrInvalidAccountDisplayName{}
	}

	return nil
}

func (a *Account) FromModel(ma models.Account) {
	(*a) = Account{
		Id:          ma.Id,
		DisplayName: ma.DisplayName,
		Username:    ma.Username,
		Type:        string(ma.Type),
		Permissions: ma.Permissions,
	}
}

type CreateSecritaryAccountParams struct {
	ActionContext
	NewAccount Account `json:"new_account"`
}

type CreateSecritaryAccountPayload struct {
	Id uint `json:"id"`
}

func (a *Actions) CreateSecritaryAccount(params CreateSecritaryAccountParams) (CreateSecritaryAccountPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteAccounts) {
		return CreateSecritaryAccountPayload{}, ErrPermissionDenied{}
	}
	if err := params.NewAccount.Validate(); err != nil {
		return CreateSecritaryAccountPayload{}, err
	}

	newAccount, err := a.app.CreateAccount(models.Account{
		DisplayName: params.NewAccount.DisplayName,
		Username:    params.NewAccount.Username,
		Password:    params.NewAccount.Password,
		Type:        models.AccountTypeSecritary,
		Permissions: secritaryPermissions,
	})

	return CreateSecritaryAccountPayload{
		Id: newAccount.Id,
	}, err
}

type CreateAdminAccountParams struct {
	ActionContext
	NewAccount Account `json:"new_account"`
}

type CreateAdminAccountPayload struct {
	Id uint `json:"id"`
}

func (a *Actions) CreateAdminAccount(params CreateAdminAccountParams) (CreateAdminAccountPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteAccounts) {
		return CreateAdminAccountPayload{}, ErrPermissionDenied{}
	}
	if err := params.NewAccount.Validate(); err != nil {
		return CreateAdminAccountPayload{}, err
	}

	newAccount, err := a.app.CreateAccount(models.Account{
		DisplayName: params.NewAccount.DisplayName,
		Username:    params.NewAccount.Username,
		Password:    params.NewAccount.Password,
		Type:        models.AccountTypeAdmin,
		Permissions: adminPermissions,
	})

	return CreateAdminAccountPayload{
		Id: newAccount.Id,
	}, err
}

type CreateJointologistAccountParams struct {
	ActionContext
	NewAccount Account `json:"new_account"`
}

type CreateJointologistAccountPayload struct {
	Id uint `json:"id"`
}

func (a *Actions) CreateJointologistAccount(params CreateJointologistAccountParams) (CreateJointologistAccountPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteAccounts) {
		return CreateJointologistAccountPayload{}, ErrPermissionDenied{}
	}
	if err := params.NewAccount.Validate(); err != nil {
		return CreateJointologistAccountPayload{}, err
	}

	newAccount, err := a.app.CreateAccount(models.Account{
		DisplayName: params.NewAccount.DisplayName,
		Username:    params.NewAccount.Username,
		Password:    params.NewAccount.Password,
		Type:        models.AccountTypeJointologist,
		Permissions: snoopDoggPermissions,
	})

	return CreateJointologistAccountPayload{
		Id: newAccount.Id,
	}, err
}

type GetAccountParams struct {
	ActionContext
	AccountId uint
}

type GetAccountPayload struct {
	Account Account `json:"data"`
}

func (a *Actions) GetAccount(params GetAccountParams) (GetAccountPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadAccounts) {
		return GetAccountPayload{}, ErrPermissionDenied{}
	}

	account, err := a.app.GetAccountById(params.AccountId)
	if err != nil {
		return GetAccountPayload{}, err
	}

	outAccount := new(Account)
	outAccount.FromModel(account)

	return GetAccountPayload{
		Account: *outAccount,
	}, nil
}

type DeleteAccountParams struct {
	ActionContext
	AccountId uint
}

type DeleteAccountPayload struct {
}

func (a *Actions) DeleteAccount(params DeleteAccountParams) (DeleteAccountPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteAccounts) {
		return DeleteAccountPayload{}, ErrPermissionDenied{}
	}

	err := a.app.DeleteAccount(params.AccountId)
	if err != nil {
		return DeleteAccountPayload{}, err
	}

	return DeleteAccountPayload{}, nil
}

type UpdateAccountParams struct {
	ActionContext
	AccountId  uint
	NewAccount Account `json:"new_account"`
}

type UpdateAccountPayload struct {
}

func (a *Actions) UpdateAccount(params UpdateAccountParams) (UpdateAccountPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteAccounts) {
		return UpdateAccountPayload{}, ErrPermissionDenied{}
	}

	err := a.app.UpdateAccount(params.AccountId, models.Account{
		DisplayName: params.NewAccount.DisplayName,
		Username:    params.NewAccount.Username,
		Password:    params.NewAccount.Password,
		Permissions: params.NewAccount.Permissions,
	})
	if err != nil {
		return UpdateAccountPayload{}, err
	}

	err = a.cache.InvalidateAuthenticatedAccountById(params.AccountId)
	if err != nil {
		return UpdateAccountPayload{}, err
	}

	return UpdateAccountPayload{}, nil
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

	accounts, err := a.app.ListAllAccounts()
	if err != nil {
		return ListAllAccountsPayload{}, err
	}

	outAccounts := make([]Account, 0, len(accounts))
	for _, account := range accounts {
		outAccount := new(Account)
		outAccount.FromModel(account)
		outAccounts = append(outAccounts, *outAccount)
	}

	return ListAllAccountsPayload{
		Data: outAccounts,
	}, nil
}
