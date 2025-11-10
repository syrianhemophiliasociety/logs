package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shs-web/actions"
	"shs-web/config"
	"shs-web/handlers/middlewares/auth"
	"shs-web/log"
	verrors "shs-web/views/errors"
	"time"
)

type usernameLoginApi struct {
	usecases *actions.Actions
}

func NewUsernameLoginApi(usecases *actions.Actions) *usernameLoginApi {
	return &usernameLoginApi{
		usecases: usecases,
	}
}

func (e *usernameLoginApi) HandleUsernameLogin(w http.ResponseWriter, r *http.Request) {
	var reqBody actions.LoginWithUsernameParams
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := e.usecases.LoginWithUsername(reqBody)
	if err != nil {
		log.Errorf("[USERNAME LOGIN API]: Failed to login user: %+v, error: %s\n", reqBody, err.Error())
		verrors.
			BugsBunnyError(fmt.Sprintf("No account associated with the username \"%s\" was found", reqBody.Username)).
			Render(r.Context(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionTokenKey,
		Value:    payload.SessionToken,
		HttpOnly: true,
		Path:     "/",
		Domain:   config.Env().Hostname,
		Expires:  time.Now().UTC().Add(time.Hour * 24 * 60),
	})

	w.Header().Set("HX-Redirect", "/")
}
