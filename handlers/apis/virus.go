package apis

import (
	"encoding/json"
	"net/http"
	"syrianhemophiliasociety/logs-web/actions"
	"syrianhemophiliasociety/logs-web/i18n"
	"syrianhemophiliasociety/logs-web/log"
	"syrianhemophiliasociety/logs-web/views/components"
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

func (v *virusApi) HandleCreateVirus(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody actions.CreateVirusRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.CreateVirus(actions.CreateVirusParams{
		RequestContext: ctx,
		NewVirus:       reqBody,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *virusApi) HandleDeleteVirus(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	id := r.PathValue("id")
	intId, _ := strconv.Atoi(id)

	_, err = v.usecases.DeleteVirus(actions.DeleteVirusParams{
		RequestContext: ctx,
		VirusId:        uint(intId),
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}
