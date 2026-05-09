package apis

import "net/http"

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
