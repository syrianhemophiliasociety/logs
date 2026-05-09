package app

import (
	"shs/app/models"

	"golang.org/x/crypto/bcrypt"
)

func (a *App) GetAccountByUsername(username string) (models.Account, error) {
	return a.repo.GetAccountByUsername(username)
}

func (a *App) GetAccountById(id uint) (models.Account, error) {
	return a.repo.GetAccount(id)
}

func (a *App) CreateAccount(account models.Account) (models.Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.Account{}, err
	}

	account.Password = string(hashedPassword)
	return a.repo.CreateAccount(account)
}

func (a *App) ListAllAccounts() ([]models.Account, error) {
	return a.repo.ListAllAccounts()
}

func (a *App) UpdateAccount(id uint, newAccount models.Account) error {
	oldAccount, err := a.repo.GetAccount(id)
	if err != nil {
		return err
	}

	if newAccount.Username != "" && newAccount.Username != oldAccount.Username {
		err := a.repo.UpdateAccountUsername(id, newAccount.Username)
		if err != nil {
			return err
		}
	}

	if newAccount.DisplayName != "" && newAccount.DisplayName != oldAccount.DisplayName {
		err := a.repo.UpdateAccountDisplayName(id, newAccount.DisplayName)
		if err != nil {
			return err
		}
	}

	if newAccount.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newAccount.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		if err = bcrypt.CompareHashAndPassword([]byte(oldAccount.Password), []byte(newAccount.Password)); err != nil {
			err = a.repo.UpdateAccountPassword(id, string(hashedPassword))
			if err != nil {
				return err
			}
		}
	}

	if newAccount.Permissions != 0 && newAccount.Permissions != oldAccount.Permissions {
		err := a.repo.UpdateAccountPermissions(id, newAccount.Permissions)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) DeleteAccount(id uint) error {
	return a.repo.DeleteAccount(id)
}
