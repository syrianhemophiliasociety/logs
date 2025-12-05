package actions

import "shs/app/models"

type Cache interface {
	SetAuthenticatedAccount(sessionToken string, account models.Account) error
	GetAuthenticatedAccount(sessionToken string) (models.Account, error)
	InvalidateAuthenticatedAccount(sessionToken string) error
	InvalidateAuthenticatedAccountById(accountId uint) error
}
