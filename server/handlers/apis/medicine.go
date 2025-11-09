package apis

import (
	"encoding/json"
	"net/http"
	"shs/actions"
	"shs/log"
	"strconv"
)

type medicineApi struct {
	usecases *actions.Actions
}

func NewMedicineApi(usecases *actions.Actions) *medicineApi {
	return &medicineApi{
		usecases: usecases,
	}
}

func (e *medicineApi) HandleCreateMedicine(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	var reqBody actions.CreateMedicineParams
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		handleErrorResponse(w, err)
		return
	}
	reqBody.ActionContext = ctx

	payload, err := e.usecases.CreateMedicine(reqBody)
	if err != nil {
		log.Errorf("[MEDICINE API]: Failed to create medicine: %+v, error: %s\n", reqBody, err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *medicineApi) HandleListMedicines(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	payload, err := e.usecases.ListAllMedicine(actions.ListAllMedicineParams{
		ActionContext: ctx,
	})
	if err != nil {
		log.Errorf("[MEDICINE API]: Failed to get medicines, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (e *medicineApi) HandleDeleteMedicine(w http.ResponseWriter, r *http.Request) {
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

	payload, err := e.usecases.DeleteMedicine(actions.DeleteMedicineParams{
		ActionContext: ctx,
		MedicineId:    uint(id),
	})
	if err != nil {
		log.Errorf("[MEDICINE API]: Failed to delete medicine, error: %s\n", err.Error())
		handleErrorResponse(w, err)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}
