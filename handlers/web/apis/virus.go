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

type CreateVirusRequest struct {
	Name         string `json:"name"`
	BloodTestIds []uint `json:"blood_test_ids"`
}

func (v *CreateVirusRequest) UnmarshalJSON(payload []byte) error {
	var data map[string]any
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}

	var ok bool
	(*v).Name, ok = data["name"].(string)
	if !ok {
		return errors.New("missing name")
	}

	const bloodTestKey = "blood_test_id"
	switch data[bloodTestKey].(type) {
	case string:
		btIdInt, err := strconv.Atoi(data[bloodTestKey].(string))
		if err != nil {
			return err
		}
		(*v).BloodTestIds = []uint{uint(btIdInt)}

	case []any:
		for _, btId := range data[bloodTestKey].([]any) {
			btIdStr, ok := btId.(string)
			if !ok {
				return errors.New("invalid blood_test_id type")
			}
			btIdInt, err := strconv.Atoi(btIdStr)
			if err != nil {
				return err
			}
			(*v).BloodTestIds = append((*v).BloodTestIds, uint(btIdInt))
		}

	default:
		return errors.New("invalid blood_test_id value")
	}

	return nil
}

/////

type virusApi struct {
	usecases *actions.Actions
}

func NewVirusApi(usecases *actions.Actions) *virusApi {
	return &virusApi{
		usecases: usecases,
	}
}

func (v *virusApi) HandleCreateVirus(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody CreateVirusRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.CreateVirus(actions.CreateVirusParams{
		ActionContext: ctx,
		NewVirus: actions.Virus{
			Name:         reqBody.Name,
			BloodTestIds: reqBody.BloodTestIds,
		},
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *virusApi) HandleDeleteVirus(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	id := r.PathValue("id")
	intId, _ := strconv.Atoi(id)

	_, err = v.usecases.DeleteVirus(actions.DeleteVirusParams{
		ActionContext: ctx,
		VirusId:       uint(intId),
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}
