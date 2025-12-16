package apis

import (
	"encoding/json"
	"io"
	"net/http"
	"shs-web/actions"
	"shs-web/i18n"
	"shs-web/log"
	"shs-web/views/components"
	"strconv"
)

type diagnosisApi struct {
	usecases *actions.Actions
}

func NewDiagnosisApi(usecases *actions.Actions) *diagnosisApi {
	return &diagnosisApi{
		usecases: usecases,
	}
}

func (v *diagnosisApi) HandleCreateDiagnosis(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody actions.Diagnosis
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
	}

	_, err = v.usecases.CreateDiagnosis(actions.CreateDiagnosisParams{
		RequestContext: ctx,
		NewDiagnosis:   reqBody,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *diagnosisApi) HandleDeleteDiagnosis(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	id := r.PathValue("id")
	intId, _ := strconv.Atoi(id)

	_, err = v.usecases.DeleteDiagnosis(actions.DeleteDiagnosisParams{
		RequestContext: ctx,
		DiagnosisId:    uint(intId),
	})
	if err != nil {
		writeRawTextResponse(w, i18n.Strings("en").ErrorSomethingWentWrong)
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}
