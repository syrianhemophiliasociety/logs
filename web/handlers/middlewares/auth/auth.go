package auth

import (
	"context"
	"net/http"
	"shs-web/actions"
	"shs-web/handlers/middlewares/clienthash"
	"shs-web/handlers/middlewares/contenttype"
	"slices"
)

// Cookie keys
const (
	VerificationTokenKey = "verification-token"
	SessionTokenKey      = "token"
)

// Context keys
const (
	CtxSessionTokenKey = "session-token"
	AccountType        = "account-type"
)

var noAuthPaths = []string{"/login", "/signup"}

type Middleware struct {
	usecases *actions.Actions
}

// New returns a new auth middle ware instance.
func New(usecases *actions.Actions) *Middleware {
	return &Middleware{
		usecases: usecases,
	}
}

// AuthPage authenticates a page's handler.
func (a *Middleware) AuthPage(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmxRedirect := contenttype.IsNoLayoutPage(r)
		sessionToken, err := a.authenticate(r)
		authed := err == nil
		ctx := context.WithValue(r.Context(), CtxSessionTokenKey, sessionToken)

		switch {
		case authed && slices.Contains(noAuthPaths, r.URL.Path):
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		case !authed && slices.Contains(noAuthPaths, r.URL.Path):
			h(w, r.WithContext(ctx))
		case !authed && htmxRedirect:
			clientHash, ok := r.Context().Value(clienthash.ClientHashKey).(string)
			if ok {
				_ = a.usecases.SetRedirectPath(clientHash, r.URL.Path)
			}
			w.Header().Set("HX-Redirect", "/login")
		case !authed && !htmxRedirect:
			clientHash, ok := r.Context().Value(clienthash.ClientHashKey).(string)
			if ok {
				_ = a.usecases.SetRedirectPath(clientHash, r.URL.Path)
			}
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		default:
			h(w, r.WithContext(ctx))
		}
	}
}

// OptionalAuthPage authenticates a page's handler optionally (without redirection).
func (a *Middleware) OptionalAuthPage(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := a.authenticate(r)
		if err != nil {
			h(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), CtxSessionTokenKey, sessionToken)
		h(w, r.WithContext(ctx))
	}
}

// AuthApi authenticates an API's handler.
func (a *Middleware) AuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := a.authenticate(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), CtxSessionTokenKey, sessionToken)
		h(w, r.WithContext(ctx))
	}
}

// OptionalAuthApi authenticates a page's handler optionally (without 401).
func (a *Middleware) OptionalAuthApi(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := a.authenticate(r)
		if err != nil {
			h(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), CtxSessionTokenKey, sessionToken)
		h(w, r.WithContext(ctx))
	}
}

func (a *Middleware) authenticate(r *http.Request) (string, error) {
	sessionToken, err := r.Cookie(SessionTokenKey)
	if err != nil {
		return "", err
	}

	return sessionToken.Value, a.usecases.CheckAuth(sessionToken.Value)
}
