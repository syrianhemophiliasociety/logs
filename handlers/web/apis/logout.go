package apis

import (
	"net/http"
	"shs/actions"
	"shs/config"
	"shs/handlers/middlewares/webauth"
)

type logoutApi struct {
	usecases *actions.Actions
}

func NewLogoutApi(usecases *actions.Actions) *logoutApi {
	return &logoutApi{
		usecases: usecases,
	}
}

func (l *logoutApi) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(webauth.CtxSessionTokenKey).(string)
	_ = l.usecases.Logout(sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:   webauth.SessionTokenKey,
		Value:  "",
		Path:   "/",
		Domain: config.Env().Hostname,
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
