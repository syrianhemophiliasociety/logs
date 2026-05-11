package apis

import (
	"encoding/json"
	"errors"
	"net/http"
	"shs/actions"
	"shs/app/models"
	"shs/handlers/web/context"
	"shs/log"
	"shs/web/i18n"
	"shs/web/views/components"
	"strconv"
)

type Account struct {
	Id          uint                      `json:"id"`
	DisplayName string                    `json:"display_name"`
	Username    string                    `json:"username"`
	Type        string                    `json:"type"`
	Password    string                    `json:"password,omitempty"`
	Permissions models.AccountPermissions `json:"permissions"`
}

type UpdateAccountRequest struct {
	Account
}

func (a *UpdateAccountRequest) UnmarshalJSON(payload []byte) error {
	var data map[string]any
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}

	var ok bool
	(*a).DisplayName, ok = data["display_name"].(string)
	if !ok {
		return errors.New("invalid display_name value")
	}
	(*a).Username, ok = data["username"].(string)
	if !ok {
		return errors.New("invalid username value")
	}
	(*a).Password, ok = data["password"].(string)
	if !ok {
		return errors.New("invalid password value")
	}

	const permissionsKey = "permissions"
	switch data[permissionsKey].(type) {
	case string:
		p, err := strconv.Atoi(data[permissionsKey].(string))
		if err != nil {
			return err
		}
		if (p & (p - 1)) != 0 {
			return errors.New("invalid permissions value")
		}
		(*a).Permissions = models.AccountPermissions(p)

	case []any:
		for _, p := range data[permissionsKey].([]any) {
			pStr, ok := p.(string)
			if !ok {
				return errors.New("invalid permissions type")
			}
			pInt, err := strconv.Atoi(pStr)
			if err != nil {
				return err
			}
			if (pInt & (pInt - 1)) != 0 {
				return errors.New("invalid permissions value")
			}
			(*a).Permissions |= models.AccountPermissions(pInt)
		}

	default:
		return errors.New("invalid permissions value")
	}

	return nil
}

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

	var params UpdateAccountRequest
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.UpdateAccount(actions.UpdateAccountParams{
		ActionContext: ctx,
		AccountId:     uint(intId),
		NewAccount: actions.Account{
			Id:          uint(intId),
			DisplayName: params.DisplayName,
			Username:    params.Username,
			Type:        params.Type,
			Password:    params.Password,
			Permissions: params.Permissions,
		},
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
