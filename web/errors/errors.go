package errors

import "errors"

var (
	ErrInvalidToken             = errors.New("invalid-token")
	ErrExpiredToken             = errors.New("expired-token")
	ErrAccountNotFound          = errors.New("account-not-found")
	ErrProfileNotFound          = errors.New("profile-not-found")
	ErrAccountExists            = errors.New("account-exists")
	ErrProfileExists            = errors.New("profile-exists")
	ErrDifferentLoginMethodUsed = errors.New("different-login-method-used")
	ErrVerificationCodeExpired  = errors.New("verification-code-expired")
	ErrInvalidVerificationCode  = errors.New("invalid-verification-code")
	ErrInvalidSessionToken      = errors.New("invalid-session-token")

	ErrSomethingWentWrong = errors.New("something went wrong")
)
