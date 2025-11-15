package actions

import (
	"errors"
	"shs/app/models"
)

type Medicine struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Dose int    `json:"dose"`
	Unit string `json:"unit"`
}

type CreateMedicineParams struct {
	ActionContext
	NewMedicine Medicine `json:"new_medicine"`
}

type CreateMedicinePayload struct {
}

func (a *Actions) CreateMedicine(params CreateMedicineParams) (CreateMedicinePayload, error) {
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
	if err != nil {
		return CreateMedicinePayload{}, err
	}

	_, err = a.app.CreateMedicine(models.Medicine{
		Name: params.NewMedicine.Name,
		Dose: params.NewMedicine.Dose,
		Unit: params.NewMedicine.Unit,
	})

	return CreateMedicinePayload{}, err
}

type UpdateMedicineParams struct {
	ActionContext
}

type UpdateMedicinePayload struct {
}

func (a *Actions) UpdateMedicine(params UpdateMedicineParams) (UpdateMedicinePayload, error) {
	return UpdateMedicinePayload{}, errors.New("not implemented")
}

type DeleteMedicineParams struct {
	ActionContext
	MedicineId uint `json:"medicine_id"`
}

type DeleteMedicinePayload struct {
}

func (a *Actions) DeleteMedicine(params DeleteMedicineParams) (DeleteMedicinePayload, error) {
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
	if err != nil {
		return DeleteMedicinePayload{}, err
	}

	return DeleteMedicinePayload{}, a.app.DeleteMedicine(params.MedicineId)
}

type ListAllMedicineParams struct {
	ActionContext
}

type ListAllMedicinePayload struct {
	Data []Medicine `json:"data"`
}

func (a *Actions) ListAllMedicine(params ListAllMedicineParams) (ListAllMedicinePayload, error) {
	medicines, err := a.app.ListAllMedicines()
	if err != nil {
		return ListAllMedicinePayload{}, err
	}

	outMedicines := make([]Medicine, 0, len(medicines))
	for _, medicine := range medicines {
		outMedicines = append(outMedicines, Medicine{
			Id:   medicine.Id,
			Name: medicine.Name,
			Dose: medicine.Dose,
			Unit: medicine.Unit,
		})
	}

	return ListAllMedicinePayload{
		Data: outMedicines,
	}, nil
}
