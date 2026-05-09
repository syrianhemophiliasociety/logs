package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/log"
)

type addressApi struct {
	usecases *actions.Actions
}

func NewAddressApi(usecases *actions.Actions) *addressApi {
	return &addressApi{
		usecases: usecases,
	}
}

func (e *addressApi) HandleFindAddress(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	goveronate := r.PathValue("goveronate")
	suburb := r.PathValue("suburb")
	street := r.PathValue("street")

	reqBody := actions.GetAddressesAlikeParams{
		ActionContext: ctx,
		Address: actions.Address{
			Governorate: goveronate,
			Suburb:      suburb,
			Street:      street,
		},
	}

	payload, err := e.usecases.GetAddressesAlike(reqBody)
	if err != nil {
		log.Errorf("[ADDRESS API]: Failed to find addresses: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
