package app

import (
	"fmt"
	"net/http"
	"strings"
)

// Error is implemented for every error around here :)
type Error interface {
	error
	// ClientStatusCode the HTTP status for clients.
	ClientStatusCode() int
	// ExtraData any data that will be helpful for clients for better UX context.
	ExtraData() map[string]any
	// ExposeToClients reports whether to expose this error to clients or not.
	ExposeToClients() bool
}

type ErrNotFound struct {
	ResourceName string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s-not-found", strings.ToLower(e.ResourceName))
}

func (e ErrNotFound) ClientStatusCode() int {
	return http.StatusNotFound
}

func (e ErrNotFound) ExtraData() map[string]any {
	return nil
}

func (e ErrNotFound) ExposeToClients() bool {
	return true
}

type ErrExists struct {
	ResourceName string
}

func (e ErrExists) Error() string {
	return fmt.Sprintf("%s-exists", strings.ToLower(e.ResourceName))
}

func (e ErrExists) ClientStatusCode() int {
	return http.StatusConflict
}

func (e ErrExists) ExtraData() map[string]any {
	return nil
}

func (e ErrExists) ExposeToClients() bool {
	return true
}
