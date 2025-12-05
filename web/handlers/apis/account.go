package apis

import (
	"encoding/json"
	"net/http"
	"shs-web/actions"
	"shs-web/i18n"
	"shs-web/log"
	"shs-web/views/components"
	"strconv"
)

type accountApi struct {
	usecases *actions.Actions
}

func NewAccountApi(usecases *actions.Actions) *accountApi {
	return &accountApi{
		usecases: usecases,
	}
}

func (v *accountApi) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody actions.Account
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	payload, err := v.usecases.CreateAccount(actions.CreateAccountParams{
		RequestContext: ctx,
		NewAccount:     reqBody,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	w.Header().Set("HX-Redirect", "/management/account/"+strconv.Itoa(int(payload.Id)))
}

func (v *accountApi) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	intId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		return
	}

	var reqBody actions.UpdateAccountRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.UpdateAccount(actions.UpdateAccountParams{
		RequestContext: ctx,
		AccountId:      uint(intId),
		NewAccount:     reqBody.Account,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	respRawText(w, i18n.StringsCtx(r.Context()).MessageSuccess)
}

func (v *accountApi) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	intId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		return
	}

	_, err = v.usecases.DeleteAccount(actions.DeleteAccountParams{
		RequestContext: ctx,
		AccountId:      uint(intId),
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	respRawText(w, i18n.StringsCtx(r.Context()).MessageSuccess)
}
