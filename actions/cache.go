package actions

type Cache interface {
	SetAuthenticatedAccount(sessionToken string, account Account) error
	GetAuthenticatedAccount(sessionToken string) (Account, error)
	InvalidateAuthenticatedAccount(sessionToken string) error
	InvalidateAuthenticatedAccountById(accountId uint) error
	SetRedirectPath(clientHash, path string) error
	GetRedirectPath(clientHash string) (string, error)
}
