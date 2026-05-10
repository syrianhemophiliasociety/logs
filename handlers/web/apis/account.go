package apis

import (
	"encoding/json"
	"errors"
	"net/http"
	"shs/actions"
	"shs/handlers/web/context"
	"shs/log"
	"shs/web/i18n"
	"shs/web/views/components"
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
	ctx, err := context.Parse(r.Context())
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

	// TODO: fix this clusterfuckery! :) xP
	var payload1 actions.CreateSecritaryAccountPayload
	var payload2 actions.CreateAdminAccountPayload
	var payload3 actions.CreateJointologistAccountPayload
	switch reqBody.Type {
	case "secritary":
		payload1, err = v.usecases.CreateSecritaryAccount(actions.CreateSecritaryAccountParams{
			ActionContext: ctx,
			NewAccount:    reqBody,
		})
	case "admin":
		payload2, err = v.usecases.CreateAdminAccount(actions.CreateAdminAccountParams{
			ActionContext: ctx,
			NewAccount:    reqBody,
		})
	case "jointlogist":
		payload3, err = v.usecases.CreateJointologistAccount(actions.CreateJointologistAccountParams{
			ActionContext: ctx,
			NewAccount:    reqBody,
		})
	default:
		err = errors.New("idiot!")
	}
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	switch {
	case payload1.Id != 0:
		w.Header().Set("HX-Redirect", "/management/account/"+strconv.Itoa(int(payload1.Id)))
	case payload2.Id != 0:
		w.Header().Set("HX-Redirect", "/management/account/"+strconv.Itoa(int(payload2.Id)))
	case payload3.Id != 0:
		w.Header().Set("HX-Redirect", "/management/account/"+strconv.Itoa(int(payload3.Id)))
	}
}

func (v *accountApi) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
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

	var reqBody actions.Account
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.UpdateAccount(actions.UpdateAccountParams{
		ActionContext: ctx,
		AccountId:     uint(intId),
		NewAccount:    reqBody,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.StringsCtx(r.Context()).MessageSuccess)
}

func (v *accountApi) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
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
		ActionContext: ctx,
		AccountId:     uint(intId),
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	w.Header().Set("HX-Redirect", "/management")
}

func writeRawTextResponse(w http.ResponseWriter, msg string) error {
	w.Header().Set("HX-Trigger", `{"respDetails": "`+msg+`"}`)
	w.Write([]byte(msg))
	return nil
}
