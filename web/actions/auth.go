package actions

import (
	"net/http"
)

func (a *Actions) CheckAuth(sessionToken string) (Account, error) {
	return makeRequest[any, Account](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/me/auth",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
}

func (a *Actions) Logout(sessionToken string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/me/logout",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
	return err
}

func (a *Actions) SetRedirectPath(clientHash, path string) error {
	return a.cache.SetRedirectPath(clientHash, path)
}

func (a *Actions) GetRedirectPath(clientHash string) (string, error) {
	return a.cache.GetRedirectPath(clientHash)
}
