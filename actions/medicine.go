package actions

import (
	"fmt"
	"shs/app/models"
	"shs/log"
	"time"
)

type Medicine struct {
	Id           uint      `json:"id"`
	Name         string    `json:"name"`
	Dose         int       `json:"dose"`
	Unit         string    `json:"unit"`
	Amount       int       `json:"amount"`
	ExpiresAt    time.Time `json:"expires_at"`
	ReceivedAt   time.Time `json:"received_at"`
	Manufacturer string    `json:"manufacturer"`
	BatchNumber  string    `json:"batch_number"`
	FactorType   string    `json:"factor_type"`
	Factor       string    `json:"factor"`
}

func (m Medicine) DoseUnit() string {
	return fmt.Sprintf("%d %s", m.Dose, m.Unit)
}

type CreateMedicineParams struct {
	ActionContext
	NewMedicine Medicine `json:"new_medicine"`
}

func (m Medicine) IntoModel() models.Medicine {
	return models.Medicine{
		Name:         m.Name,
		Dose:         m.Dose,
		Unit:         m.Unit,
		Amount:       m.Amount,
		ExpiresAt:    m.ExpiresAt,
		ReceivedAt:   m.ReceivedAt,
		Manufacturer: m.Manufacturer,
		BatchNumber:  m.BatchNumber,
		FactorType:   m.FactorType,
		Factor:       m.Factor,
	}
}

func (m *Medicine) FromModel(medicine models.Medicine) {
	(*m) = Medicine{
		Id:           medicine.Id,
		Name:         medicine.Name,
		Dose:         medicine.Dose,
		Unit:         medicine.Unit,
		Amount:       medicine.Amount,
		ExpiresAt:    medicine.ExpiresAt,
		ReceivedAt:   medicine.ReceivedAt,
		Manufacturer: medicine.Manufacturer,
		BatchNumber:  medicine.BatchNumber,
		FactorType:   medicine.FactorType,
		Factor:       medicine.Factor,
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
	MedicineId uint `json:"medicine_id"`
	Amount     int  `json:"amount"`
}

type UpdateMedicinePayload struct {
}

func (a *Actions) UpdateMedicine(params UpdateMedicineParams) (UpdateMedicinePayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteMedicine) {
		return UpdateMedicinePayload{}, ErrPermissionDenied{}
	}

	err := a.app.UpdateMedicineAmount(params.MedicineId, params.Amount)

	return UpdateMedicinePayload{}, err
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

type GetMedicineParams struct {
	ActionContext
	MedicineId uint `json:"medicine_id"`
}

type GetMedicinePayload struct {
	Data Medicine `json:"data"`
}

func (a *Actions) GetMedicine(params GetMedicineParams) (GetMedicinePayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteMedicine) {
		return GetMedicinePayload{}, ErrPermissionDenied{}
	}

	medicine, err := a.app.GetMedicine(params.MedicineId)
	if err != nil {
		return GetMedicinePayload{}, err
	}

	outMedicine := new(Medicine)
	outMedicine.FromModel(medicine)

	return GetMedicinePayload{
		Data: *outMedicine,
	}, nil
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

type ListAllPrescribedMedicineParams struct {
	ActionContext
}

type ListAllPrescribedMedicinePayload struct {
	Data []PrescribedMedicineWithPatient `json:"data"`
}

func (a *Actions) ListAllPrescribedMedicine(params ListAllPrescribedMedicineParams) (ListAllPrescribedMedicinePayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadPatient) {
		return ListAllPrescribedMedicinePayload{}, ErrPermissionDenied{}
	}
	if !params.Account.HasPermission(models.AccountPermissionReadMedicine) {
		return ListAllPrescribedMedicinePayload{}, ErrPermissionDenied{}
	}

	prescribedMeds, err := a.app.ListAllPrescribedMedicines()
	if err != nil {
		return ListAllPrescribedMedicinePayload{}, err
	}

	medsIds := make([]uint, 0, len(prescribedMeds))
	for _, pm := range prescribedMeds {
		medsIds = append(medsIds, pm.MedicineId)
	}

	meds, err := a.app.ListMedicinesByIds(medsIds)
	if err != nil {
		return ListAllPrescribedMedicinePayload{}, err
	}

	medsMapped := make(map[uint]models.Medicine)
	for _, med := range meds {
		medsMapped[med.Id] = med
	}

	outData := make([]PrescribedMedicineWithPatient, 0, len(prescribedMeds))
	for _, medicine := range prescribedMeds {
		if medicine.UsedAt.IsZero() {
			continue
		}
		patient, err := a.app.GetPatientById(medicine.PatientId)
		if err != nil {
			log.Errorf("patient not found, error: %v\n", err)
			return ListAllPrescribedMedicinePayload{}, err
		}

		outPrescribedMedicine := new(PrescribedMedicine)
		outPrescribedMedicine.FromModel(medicine, medsMapped[medicine.MedicineId])
		outPatient := new(Patient)
		outPatient.FromModel(patient)

		outData = append(outData, PrescribedMedicineWithPatient{
			PrescribedMedicine: *outPrescribedMedicine,
			Patient:            *outPatient,
		})
	}

	return ListAllPrescribedMedicinePayload{
		Data: outData,
	}, nil
}
