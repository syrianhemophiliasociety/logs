package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/log"
	"strconv"
)

type virusApi struct {
	usecases *actions.Actions
}

func NewVirusApi(usecases *actions.Actions) *virusApi {
	return &virusApi{
		usecases: usecases,
	}
}

func (e *virusApi) HandleCreateVirus(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreateVirusParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := e.usecases.CreateVirus(reqBody)
	if err != nil {
		log.Errorf("[VIRUS API]: Failed to create virus: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *virusApi) HandleListViruses(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.ListAllViruses(actions.ListAllVirusesParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Errorf("[VIRUS API]: Failed to get viruss, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *virusApi) HandleDeleteVirus(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.DeleteVirus(actions.DeleteVirusParams{
		ActionContext: ctx,
		VirusId:       uint(id),
	})
	if err != nil {
		log.Errorf("[VIRUS API]: Failed to delete virus, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
