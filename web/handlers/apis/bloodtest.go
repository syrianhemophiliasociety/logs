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

type bloodTestApi struct {
	usecases *actions.Actions
}

func NewBloodTestApi(usecases *actions.Actions) *bloodTestApi {
	return &bloodTestApi{
		usecases: usecases,
	}
}

func (v *bloodTestApi) HandleCreateBloodTest(w http.ResponseWriter, r *http.Request) {
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

	var reqBody actions.RequestBloodTest
	var reqBody2 actions.RequestBloodTestSingle
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		err = json.Unmarshal(body, &reqBody2)
		if err != nil {
			components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
			log.Errorln(err)
			return
		}
	}

	log.Warningln("body json", string(body))
	log.Warningf("body json %+v\n", reqBody)

	_, err = v.usecases.CreateBloodTest(actions.CreateBloodTestParams{
		RequestContext:     ctx,
		NewBloodTest:       reqBody,
		NewBloodTestSingle: reqBody2,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	respRawText(w, i18n.StringsCtx(r.Context()).MessageSuccess)
}

func (v *bloodTestApi) HandleDeleteBloodTest(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	id := r.PathValue("id")
	intId, _ := strconv.Atoi(id)

	_, err = v.usecases.DeleteBloodTest(actions.DeleteBloodTestParams{
		RequestContext: ctx,
		BloodTestId:    uint(intId),
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	respRawText(w, i18n.StringsCtx(r.Context()).MessageSuccess)
}
