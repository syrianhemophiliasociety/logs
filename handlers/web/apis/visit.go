package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/handlers/web/context"
	"shs/log"
	"shs/web/i18n"
	"shs/web/views/components"
	"strconv"
)

type visitApi struct {
	usecases *actions.Actions
}

func NewVisitApi(usecases *actions.Actions) *visitApi {
	return &visitApi{
		usecases: usecases,
	}
}

type CreateTreatmentDetailsRequest struct {
	actions.TreatmentDetails
}

func (v *visitApi) HandleCreateTreatmentDetails(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody CreateTreatmentDetailsRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.CreateTreatmentDetails(actions.CreateTreatmentDetailsParams{
		ActionContext:    ctx,
		TreatmentDetails: reqBody.TreatmentDetails,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *visitApi) HandleDeleteTreatmentDetails(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	id := r.PathValue("id")
	intId, _ := strconv.Atoi(id)

	_, err = v.usecases.DeleteTreatmentDetails(actions.DeleteTreatmentDetailsParams{
		ActionContext: ctx,
		Id:            uint(intId),
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}
