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
	ErrPatientNotFound          = errors.New("patient-not-found")

	ErrSomethingWentWrong = errors.New("something went wrong")
)

type SHSError interface {
	error
	Id() string
}

type ErrInsufficientMedicineAmount struct {
	MedicineName    string `json:"medicine_name"`
	ExceedingAmount int    `json:"exceeding_amount"`
	LeftPackages    int    `json:"left_packages"`
}

func (e ErrInsufficientMedicineAmount) Error() string {
	return "insufficient-medicine-amount"
}

func (e ErrInsufficientMedicineAmount) Id() string {
	return "insufficient-medicine-amount"
}

func IsShs(err error) bool {
	var shsError SHSError
	return errors.As(err, &shsError)
}

func Is(err, target error) bool {
	if IsShs(err) {
		return true
	}

	return errors.Is(err, target)
}
