package actions

import (
	"shs/app/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	sessionTokenTtlDays = 60
)

func (a *Actions) AuthenticateAccount(sessionToken string) (models.Account, error) {
	token, err := a.jwt.Decode(sessionToken, JwtSessionToken)
	if err != nil {
		return models.Account{}, err
	}

	if !token.Payload.Valid() {
		return models.Account{}, ErrInvalidSessionToken{}
	}

	account, err := a.cache.GetAuthenticatedAccount(sessionToken)
	if err != nil {
		account, err = a.app.GetAccountByUsername(token.Payload.Username)
		if err != nil {
			return models.Account{}, err
		}

		err = a.cache.SetAuthenticatedAccount(sessionToken, account)
		if err != nil {
			return models.Account{}, err
		}
	}

	return account, nil
}

type TokenPayload struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (t TokenPayload) Valid() bool {
	return t.Name != "" && t.Username != ""
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
		return LoginWithUsernamePayload{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(params.Password))
	if err != nil {
		return LoginWithUsernamePayload{}, ErrInvalidLoginCredientials{}
	}

	sessionToken, err := a.jwt.Sign(TokenPayload{
		Name:     account.DisplayName,
		Username: account.Username,
	}, JwtSessionToken, time.Hour*24*sessionTokenTtlDays)
	if err != nil {
		return LoginWithUsernamePayload{}, err
	}

	return LoginWithUsernamePayload{
		SessionToken: sessionToken,
	}, nil
}

func (a *Actions) InvalidateAuthenticatedAccount(token string) error {
	return a.cache.InvalidateAuthenticatedAccount(token)
}
