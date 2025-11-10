package actions

import "net/http"

type LoginWithUsernameParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginWithUsernamePayload struct {
	SessionToken string `json:"session_token"`
}

func (a *Actions) LoginWithUsername(params LoginWithUsernameParams) (LoginWithUsernamePayload, error) {
	return makeRequest[LoginWithUsernameParams, LoginWithUsernamePayload](makeRequestConfig[LoginWithUsernameParams]{
		method:   http.MethodPost,
		endpoint: "/v1/login/username",
		body:     params,
	})
}
