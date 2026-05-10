package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/log"
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
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.LoginWithUsername(reqBody)
	if err != nil {
		log.Errorf("[USERNAME LOGIN API]: Failed to login user: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
