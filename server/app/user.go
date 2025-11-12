package app

import (
	"shs/app/models"

	"golang.org/x/crypto/bcrypt"
)

func (a *App) GetAccountByUsername(username string) (models.Account, error) {
	return a.repo.GetAccountByUsername(username)
}

func (a *App) CreateAccount(account models.Account) (models.Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.Account{}, err
	}

	account.Password = string(hashedPassword)
	return a.repo.CreateAccount(account)
}

func (a *App) ListAllAccountsForAdmin() ([]models.Account, error) {
	return a.repo.ListAllAccounts([]models.AccountType{models.AccountTypeSecritary})
}

func (a *App) ListAllAccountsForSuperAdmin() ([]models.Account, error) {
	return a.repo.ListAllAccounts([]models.AccountType{models.AccountTypeSecritary, models.AccountTypeAdmin})
}

func (a *App) DeleteAccount(id uint) error {
	return a.repo.DeleteAccount(id)
}
