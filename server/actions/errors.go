package actions

import "net/http"

type ErrInvalidLoginCredientials struct{}

func (e ErrInvalidLoginCredientials) Error() string {
	return "invalid-login-credentials"
}

func (e ErrInvalidLoginCredientials) ClientStatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrInvalidLoginCredientials) ExtraData() map[string]any {
	return nil
}

func (e ErrInvalidLoginCredientials) ExposeToClients() bool {
	return true
}

type ErrInvalidSessionToken struct{}

func (e ErrInvalidSessionToken) Error() string {
	return "invalid-session-token"
}

func (e ErrInvalidSessionToken) ClientStatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrInvalidSessionToken) ExtraData() map[string]any {
	return nil
}

func (e ErrInvalidSessionToken) ExposeToClients() bool {
	return true
}

type ErrInvalidVerificationToken struct{}

func (e ErrInvalidVerificationToken) Error() string {
	return "invalid-verification-code"
}

func (e ErrInvalidVerificationToken) ClientStatusCode() int {
	return http.StatusBadRequest
}

func (e ErrInvalidVerificationToken) ExtraData() map[string]any {
	return nil
}

func (e ErrInvalidVerificationToken) ExposeToClients() bool {
	return true
}
