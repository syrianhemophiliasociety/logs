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
	"time"
)

type RequestMedicine struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	Dose         string `json:"dose"`
	Unit         string `json:"unit"`
	Amount       string `json:"amount"`
	ExpiresAt    string `json:"expires_at"`
	ReceivedAt   string `json:"received_at"`
	Manufacturer string `json:"manufacturer"`
	BatchNumber  string `json:"batch_number"`
	FactorType   string `json:"factor_type"`
}

func clusterFuckMedicineToActionsOne(reqMed RequestMedicine) (actions.Medicine, error) {
	dose, err := strconv.Atoi(reqMed.Dose)
	if err != nil {
		return actions.Medicine{}, err
	}

	amount, err := strconv.Atoi(reqMed.Amount)
	if err != nil {
		return actions.Medicine{}, err
	}

	expiresAt, err := time.Parse("2006-01-02", reqMed.ExpiresAt)
	if err != nil {
		return actions.Medicine{}, err
	}

	receivedAt, err := time.Parse("2006-01-02", reqMed.ReceivedAt)
	if err != nil {
		return actions.Medicine{}, err
	}

	return actions.Medicine{
		Name:         reqMed.Name,
		Dose:         dose,
		Unit:         reqMed.Unit,
		Amount:       amount,
		ExpiresAt:    expiresAt,
		ReceivedAt:   receivedAt,
		Manufacturer: reqMed.Manufacturer,
		BatchNumber:  reqMed.BatchNumber,
		FactorType:   reqMed.FactorType,
	}, nil
}

///

type medicineApi struct {
	usecases *actions.Actions
}

func NewMedicineApi(usecases *actions.Actions) *medicineApi {
	return &medicineApi{
		usecases: usecases,
	}
}

func (v *medicineApi) HandleCreateMedicine(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody RequestMedicine
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	mmeeeeed, err := clusterFuckMedicineToActionsOne(reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.CreateMedicine(actions.CreateMedicineParams{
		ActionContext: ctx,
		NewMedicine:   mmeeeeed,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *medicineApi) HandleUpdateMedicine(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	id := r.PathValue("id")
	intId, _ := strconv.Atoi(id)

	var reqBody struct {
		Amount string `json:"amount"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	medicineAmount, err := strconv.Atoi(reqBody.Amount)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.UpdateMedicine(actions.UpdateMedicineParams{
		ActionContext: ctx,
		MedicineId:    uint(intId),
		Amount:        medicineAmount,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *medicineApi) HandleDeleteMedicine(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	id := r.PathValue("id")
	intId, _ := strconv.Atoi(id)

	_, err = v.usecases.DeleteMedicine(actions.DeleteMedicineParams{
		ActionContext: ctx,
		MedicineId:    uint(intId),
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}
