package apis

import (
	"net/http"
	"shs/actions"
)

type meApi struct {
	usecases *actions.Actions
}

func NewMeApi(usecases *actions.Actions) *meApi {
	return &meApi{
		usecases: usecases,
	}
}

func (u *meApi) HandleAuthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
}

func (m *meApi) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Header["Authorization"]
	if !ok {
		return
	}
	_ = m.usecases.InvalidateAuthenticatedAccount(sessionToken[0])
}
