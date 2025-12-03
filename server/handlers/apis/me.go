package apis

import (
	"encoding/json"
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
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(actions.Account{
		Id:          ctx.Account.Id,
		DisplayName: ctx.Account.DisplayName,
		Username:    ctx.Account.Username,
		Type:        string(ctx.Account.Type),
		Permissions: ctx.Account.Permissions,
	})
}

func (m *meApi) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Header["Authorization"]
	if !ok {
		return
	}
	_ = m.usecases.InvalidateAuthenticatedAccount(sessionToken[0])
}
