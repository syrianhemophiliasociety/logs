package apis

import "net/http"

type ErrUnauthorized struct{}

func (e ErrUnauthorized) Error() string {
	return "unauthorized"
}

func (e ErrUnauthorized) ClientStatusCode() int {
	return http.StatusUnauthorized
}

func (e ErrUnauthorized) ExtraData() map[string]any {
	return nil
}

func (e ErrUnauthorized) ExposeToClients() bool {
	return true
}

type ErrBadRequest struct {
	FieldName string
}

func (e ErrBadRequest) Error() string {
	return "bad-request"
}

func (e ErrBadRequest) ClientStatusCode() int {
	return http.StatusBadRequest
}

func (e ErrBadRequest) ExtraData() map[string]any {
	return map[string]any{
		"invalid_field": e.FieldName,
	}
}

func (e ErrBadRequest) ExposeToClients() bool {
	return true
}

type ErrInvalidFileType struct {
	Want string
	Got  string
}

func (e ErrInvalidFileType) Error() string {
	return "invalid-file-type"
}

func (e ErrInvalidFileType) ClientStatusCode() int {
	return http.StatusBadGateway
}

func (e ErrInvalidFileType) ExtraData() map[string]any {
	return map[string]any{
		"want": e.Want,
		"got":  e.Got,
	}
}

func (e ErrInvalidFileType) ExposeToClients() bool {
	return true
}
