package auth

import (
	"context"
	"net/http"
	"shs/actions"
)

// Context keys
const (
	AccountKey         = "account"
	CtxSessionTokenKey = "session-token"
)

type Middleware struct {
	usecases *actions.Actions
}

// New returns a new auth middleware instance.
func New(usecases *actions.Actions) *Middleware {
	return &Middleware{
		usecases: usecases,
	}
}

func (a *Middleware) AuthHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionToken, account, err := a.authenticate(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), AccountKey, account)
		ctx = context.WithValue(ctx, CtxSessionTokenKey, sessionToken)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthApi authenticates an API's handler.
func (a *Middleware) AuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken, account, err := a.authenticate(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), AccountKey, account)
		ctx = context.WithValue(ctx, CtxSessionTokenKey, sessionToken)
		h(w, r.WithContext(ctx))
	}
}

// OptionalAuthApi authenticates an API's handler optionally (without 401).
func (a *Middleware) OptionalAuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken, account, err := a.authenticate(r)
		if err != nil {
			h(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), AccountKey, account)
		ctx = context.WithValue(ctx, CtxSessionTokenKey, sessionToken)
		h(w, r.WithContext(ctx))
	}
}

func (a *Middleware) authenticate(r *http.Request) (string, actions.Account, error) {
	sessionToken := r.Header.Get("Authorization")
	if sessionToken == "" {
		return "", actions.Account{}, actions.ErrInvalidSessionToken{}
	}

	account, err := a.usecases.AuthenticateAccount(sessionToken)
	if err != nil {
		return "", actions.Account{}, err
	}

	return sessionToken, account, nil
}
