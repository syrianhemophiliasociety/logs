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

func (m Medicine) IntoModel() models.Medicine {
	return models.Medicine{
		Name: m.Name,
		Dose: m.Dose,
		Unit: m.Unit,
	}
}

func (m *Medicine) FromModel(medicine models.Medicine) {
	(*m) = Medicine{
		Id:   medicine.Id,
		Name: medicine.Name,
		Dose: medicine.Dose,
		Unit: medicine.Unit,
	}
}

type CreateMedicinePayload struct {
}

func (a *Actions) CreateMedicine(params CreateMedicineParams) (CreateMedicinePayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteMedicine) {
		return CreateMedicinePayload{}, ErrPermissionDenied{}
	}

	_, err := a.app.CreateMedicine(params.NewMedicine.IntoModel())

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
	if !params.Account.HasPermission(models.AccountPermissionWriteMedicine) {
		return DeleteMedicinePayload{}, ErrPermissionDenied{}
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
	if !params.Account.HasPermission(models.AccountPermissionReadMedicine) {
		return ListAllMedicinePayload{}, ErrPermissionDenied{}
	}

	medicines, err := a.app.ListAllMedicines()
	if err != nil {
		return ListAllMedicinePayload{}, err
	}

	outMedicines := make([]Medicine, 0, len(medicines))
	for _, medicine := range medicines {
		outMedicine := new(Medicine)
		outMedicine.FromModel(medicine)
		outMedicines = append(outMedicines, *outMedicine)
	}

	return ListAllMedicinePayload{
		Data: outMedicines,
	}, nil
}
