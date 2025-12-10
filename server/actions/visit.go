package actions

import (
	"shs/app/models"
	"time"
)

type Visit struct {
	Reason             string               `json:"reason"`
	VisitedAt          time.Time            `json:"visited_at"`
	PrescribedMedicine []PrescribedMedicine `json:"prescribed_medicine"`
}

type CreatePatientVisitParams struct {
	ActionContext
	PatientId             string
	VisitReason           string `json:"visit_reason"`
	PrescribedMedicineIds []uint `json:"prescribed_medicine_ids"`
}

type CreatePatientVisitPayload struct {
}

func (a *Actions) CreatePatientVisit(params CreatePatientVisitParams) (CreatePatientVisitPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteOtherVisits) {
		return CreatePatientVisitPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetMinimalPatientByPublicId(params.PatientId)
	if err != nil {
		return CreatePatientVisitPayload{}, err
	}

	meds, err := a.app.ListMedicinesByIds(params.PrescribedMedicineIds)
	if err != nil {
		return CreatePatientVisitPayload{}, err
	}

	if len(meds) != len(params.PrescribedMedicineIds) {
		// TODO: do something
	}

	visit, err := a.app.CreatePatientVisit(models.Visit{
		PatientId: patient.Id,
		Reason:    models.VisitReason(params.VisitReason),
	})
	if err != nil {
		return CreatePatientVisitPayload{}, err
	}

	for _, medId := range params.PrescribedMedicineIds {
		_, err = a.app.CreatePrescribedMedicine(models.PrescribedMedicine{
			VisitId:    visit.Id,
			PatientId:  patient.Id,
			MedicineId: medId,
		})
		if err != nil {
			return CreatePatientVisitPayload{}, err
		}
	}

	return CreatePatientVisitPayload{}, nil
}

type PrescribedMedicine struct {
	Medicine
	PrescribedMedicineId uint      `json:"prescribed_medicine_id"`
	UsedAt               time.Time `json:"used_at"`
}

type GetPatientLastVisitParams struct {
	ActionContext
}

type GetPatientLastVisitPayload struct {
	Patient            Patient              `json:"patient"`
	VisitedAt          time.Time            `json:"visited_at"`
	PrescribedMedicine []PrescribedMedicine `json:"prescribed_medicine"`
}

func (a *Actions) GetPatientLastVisit(params GetPatientLastVisitParams) (GetPatientLastVisitPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadOtherVisits) {
		return GetPatientLastVisitPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetMinimalPatientByPublicId(params.Account.Username)
	if err != nil {
		return GetPatientLastVisitPayload{}, err
	}

	lastVisit, err := a.app.GetPatientLastVisit(patient.Id)
	if err != nil {
		return GetPatientLastVisitPayload{}, err
	}

	prescribedMeds, err := a.app.ListPatientVisitPrescribedMedicine(lastVisit.Id)
	if err != nil {
		return GetPatientLastVisitPayload{}, err
	}

	medsIds := make([]uint, 0, len(prescribedMeds))
	for _, pm := range prescribedMeds {
		medsIds = append(medsIds, pm.MedicineId)
	}

	meds, err := a.app.ListMedicinesByIds(medsIds)
	if err != nil {
		return GetPatientLastVisitPayload{}, err
	}

	medsMapped := make(map[uint]Medicine)
	for _, med := range meds {
		medsMapped[med.Id] = Medicine{
			Id:   med.Id,
			Name: med.Name,
			Dose: med.Dose,
			Unit: med.Unit,
		}
	}

	outMeds := make([]PrescribedMedicine, 0, len(prescribedMeds))
	for _, pm := range prescribedMeds {
		outMeds = append(outMeds, PrescribedMedicine{
			Medicine:             medsMapped[pm.MedicineId],
			PrescribedMedicineId: pm.Id,
			UsedAt:               pm.UsedAt,
		})
	}

	outPatient := new(Patient)
	outPatient.FromModel(patient)

	return GetPatientLastVisitPayload{
		Patient:            *outPatient,
		PrescribedMedicine: outMeds,
		VisitedAt:          lastVisit.CreatedAt,
	}, nil
}

type UseMedicineForVisitParams struct {
	ActionContext
	PrescribedMedicineId uint `json:"prescribed_medicine_id"`
}

type UseMedicineForVisitPayload struct {
}

func (a *Actions) UseMedicineForVisit(params UseMedicineForVisitParams) (UseMedicineForVisitPayload, error) {
	return UseMedicineForVisitPayload{}, nil
}

type ListPatientVisitsParams struct {
	ActionContext
	PatientId string
}

type ListPatientVisitsPayload struct {
	Data []Visit `json:"data"`
}

func (a *Actions) ListPatientVisits(params ListPatientVisitsParams) (ListPatientVisitsPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadOtherVisits) {
		return ListPatientVisitsPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetMinimalPatientByPublicId(params.PatientId)
	if err != nil {
		return ListPatientVisitsPayload{}, err
	}

	visits, err := a.app.ListPatientVisits(patient.Id)
	if err != nil {
		return ListPatientVisitsPayload{}, err
	}

	outVisits := make([]Visit, 0, len(visits))
	for _, visit := range visits {
		outVisits = append(outVisits, Visit{
			Reason:             string(visit.Reason),
			VisitedAt:          visit.CreatedAt,
			PrescribedMedicine: []PrescribedMedicine{},
		})
	}

	return ListPatientVisitsPayload{
		Data: outVisits,
	}, nil
}
