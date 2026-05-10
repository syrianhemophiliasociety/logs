package actions

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	sessionTokenTtlDays = 60
)

func (a *Actions) AuthenticateAccount(sessionToken string) (Account, error) {
	token, err := a.jwt.Decode(sessionToken, JwtSessionToken)
	if err != nil {
		return Account{}, err
	}

	if !token.Payload.Valid() {
		return Account{}, ErrInvalidSessionToken{}
	}

	// TODO: add some state for the session token to indicate logged out accounts
	account, err := a.cache.GetAuthenticatedAccount(sessionToken)
	if err != nil {
		dbaccount, err := a.app.GetAccountByUsername(token.Payload.Username)
		if err != nil {
			return Account{}, err
		}

		account = Account{
			Id:          dbaccount.Id,
			DisplayName: dbaccount.DisplayName,
			Username:    dbaccount.Username,
			Type:        string(dbaccount.Type),
			Password:    "",
			Permissions: dbaccount.Permissions,
		}

		err = a.cache.SetAuthenticatedAccount(sessionToken, account)
		if err != nil {
			return Account{}, err
		}
	}

	account.Password = ""

	return account, nil
}

func (a *Actions) CheckSessionToken(sessionToken string) error {
	_, err := a.cache.GetAuthenticatedAccount(sessionToken)
	if err != nil {
		return ErrInvalidSessionToken{}
	}

	return nil
}

type TokenPayload struct {
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (t TokenPayload) Valid() bool {
	return t.Name != "" && t.Username != "" && !t.CreatedAt.IsZero()
}

type LoginWithUsernameParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginWithUsernamePayload struct {
	SessionToken string `json:"session_token"`
}

func (a *Actions) LoginWithUsername(params LoginWithUsernameParams) (LoginWithUsernamePayload, error) {
	account, err := a.app.GetAccountByUsername(params.Username)
	if err != nil {
		return LoginWithUsernamePayload{}, ErrInvalidLoginCredientials{}
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(params.Password))
	if err != nil {
		return LoginWithUsernamePayload{}, ErrInvalidLoginCredientials{}
	}

	sessionToken, err := a.jwt.Sign(TokenPayload{
		Name:      account.DisplayName,
		Username:  account.Username,
		CreatedAt: time.Now().UTC(),
	}, JwtSessionToken, time.Hour*24*sessionTokenTtlDays)
	if err != nil {
		return LoginWithUsernamePayload{}, err
	}

	return LoginWithUsernamePayload{
		SessionToken: sessionToken,
	}, nil
}

func (a *Actions) Logout(token string) error {
	return a.cache.InvalidateAuthenticatedAccount(token)
}

func (a *Actions) InvalidateAuthenticatedAccount(token string) error {
	return a.cache.InvalidateAuthenticatedAccount(token)
}

func (a *Actions) SetRedirectPath(clientHash, path string) error {
	return a.cache.SetRedirectPath(clientHash, path)
}

func (a *Actions) GetRedirectPath(clientHash string) (string, error) {
	return a.cache.GetRedirectPath(clientHash)
}
