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

type ErrPermissionDenied struct{}

func (e ErrPermissionDenied) Error() string {
	return "permission-denied"
}

func (e ErrPermissionDenied) ClientStatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrPermissionDenied) ExtraData() map[string]any {
	return nil
}

func (e ErrPermissionDenied) ExposeToClients() bool {
	return true
}

type ErrValidation struct {
	Field string
}

func (e ErrValidation) Error() string {
	return "invalid-field"
}

func (e ErrValidation) ClientStatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrValidation) ExtraData() map[string]any {
	return map[string]any{
		"field_name": e.Field,
	}
}

func (e ErrValidation) ExposeToClients() bool {
	return true
}

type ErrInsufficientMedicine struct {
	MedicineName    string
	ExceedingAmount int
	LeftPackages    int
}

func (e ErrInsufficientMedicine) Error() string {
	return "insufficient-medicine-amount"
}

func (e ErrInsufficientMedicine) ClientStatusCode() int {
	return http.StatusForbidden
}

func (e ErrInsufficientMedicine) ExtraData() map[string]any {
	return map[string]any{
		"medicine_name":    e.MedicineName,
		"exceeding_amount": e.ExceedingAmount,
		"left_packages":    e.LeftPackages,
	}
}

func (e ErrInsufficientMedicine) ExposeToClients() bool {
	return true
}
