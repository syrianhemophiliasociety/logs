package actions

import (
	"net/http"
)

type Profile struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	PfpLink  string `json:"pfp_link"`
	Username string `json:"username"`
}

func (a *Actions) CheckAuth(sessionToken string) error {
	_, err := makeRequest[any, any](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/me/auth",
		headers: map[string]string{
			"Authorization": sessionToken,
		},
	})
	return err
}

func (a *Actions) Logout(sessionToken string) error {
	_, err := makeRequest[any, Profile](makeRequestConfig[any]{
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
