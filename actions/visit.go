package actions

import (
	"shs/app/models"
	"time"
)

type Visit struct {
	Id                 uint                 `json:"id"`
	Reason             string               `json:"reason"`
	ExtraNote          string               `json:"extra_note"`
	VisitedAt          time.Time            `json:"visited_at"`
	PatientWeight      float64              `json:"patient_weight"`
	PatientHeight      float64              `json:"patient_height"`
	PrescribedMedicine []PrescribedMedicine `json:"prescribed_medicine"`
}

type CreatePatientVisitParams struct {
	ActionContext
	PatientId           string
	VisitReason         string     `json:"visit_reason"`
	VisitExtraDetails   string     `json:"visit_extra_details"`
	PatientWeight       float64    `json:"patient_weight"`
	PatientHeight       float64    `json:"patient_height"`
	PrescribedMedicines []Medicine `json:"prescribed_medicines"`
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

	medIds := make([]uint, 0, len(params.PrescribedMedicines))
	for _, med := range params.PrescribedMedicines {
		medIds = append(medIds, med.Id)
	}

	meds, err := a.app.ListMedicinesByIds(medIds)
	if err != nil {
		return CreatePatientVisitPayload{}, err
	}

	prescribedMedicinesAmount := make(map[uint]int)
	for _, med := range params.PrescribedMedicines {
		prescribedMedicinesAmount[med.Id] += med.Amount
	}

	for _, med := range meds {
		if prescribedMedicinesAmount[med.Id] > med.Amount {
			return CreatePatientVisitPayload{}, ErrInsufficientMedicine{
				MedicineName:    med.Name,
				ExceedingAmount: prescribedMedicinesAmount[med.Id],
				LeftPackages:    med.Amount,
			}
		}
	}

	visit, err := a.app.CreatePatientVisit(models.Visit{
		PatientId:     patient.Id,
		Reason:        models.VisitReason(params.VisitReason),
		Notes:         params.VisitExtraDetails,
		PatientWeight: params.PatientWeight,
		PatientHeight: params.PatientHeight,
	})
	if err != nil {
		return CreatePatientVisitPayload{}, err
	}

	for _, med := range params.PrescribedMedicines {
		for range med.Amount {
			_, err = a.app.CreatePrescribedMedicine(models.PrescribedMedicine{
				VisitId:    visit.Id,
				PatientId:  patient.Id,
				MedicineId: med.Id,
			})
		}
		if err != nil {
			return CreatePatientVisitPayload{}, err
		}
		err = a.app.DecrementMedicineAmount(med.Id, med.Amount)
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

func (pm *PrescribedMedicine) FromModel(m models.PrescribedMedicine, med models.Medicine) {
	outMed := new(Medicine)
	outMed.FromModel(med)
	(*pm).Medicine = *outMed
	(*pm).PrescribedMedicineId = m.Id
	(*pm).UsedAt = m.UsedAt
}

func (pm PrescribedMedicine) IntoModel(visitId, patientId, medicineId uint) models.PrescribedMedicine {
	return models.PrescribedMedicine{
		VisitId:    visitId,
		PatientId:  patientId,
		MedicineId: medicineId,
	}
}

type GetPatientLastVisitParams struct {
	ActionContext
}

type GetPatientLastVisitPayload struct {
	VisitId            uint                 `json:"visit_id"`
	Patient            Patient              `json:"patient"`
	VisitedAt          time.Time            `json:"visited_at"`
	PatientWeight      float64              `json:"patient_weight"`
	PatientHeight      float64              `json:"patient_height"`
	PrescribedMedicine []PrescribedMedicine `json:"prescribed_medicine"`
}

func (a *Actions) GetPatientLastVisit(params GetPatientLastVisitParams) (GetPatientLastVisitPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadOwnVisit) {
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

	medsMapped := make(map[uint]models.Medicine)
	for _, med := range meds {
		medsMapped[med.Id] = med
	}

	outMeds := make([]PrescribedMedicine, 0, len(prescribedMeds))
	for _, pm := range prescribedMeds {
		outMed := new(PrescribedMedicine)
		outMed.FromModel(pm, medsMapped[pm.MedicineId])
		outMeds = append(outMeds, *outMed)
	}

	outPatient := new(Patient)
	outPatient.FromModel(patient)

	return GetPatientLastVisitPayload{
		Patient:            *outPatient,
		PrescribedMedicine: outMeds,
		VisitedAt:          lastVisit.CreatedAt,
		VisitId:            lastVisit.Id,
		PatientWeight:      lastVisit.PatientWeight,
		PatientHeight:      lastVisit.PatientHeight,
	}, nil
}

type UseMedicineForVisitParams struct {
	ActionContext
	PrescribedMedicineId uint `json:"prescribed_medicine_id"`
	VisitId              uint `json:"visit_id"`
}

type UseMedicineForVisitPayload struct {
}

func (a *Actions) UseMedicineForVisit(params UseMedicineForVisitParams) (UseMedicineForVisitPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteOwnVisit) {
		return UseMedicineForVisitPayload{}, ErrPermissionDenied{}
	}

	err := a.app.UseMedicineForVisit(params.PrescribedMedicineId, params.VisitId)
	if err != nil {
		return UseMedicineForVisitPayload{}, err
	}

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
		prescribedMeds, err := a.app.ListPatientVisitPrescribedMedicine(visit.Id)
		if err != nil {
			return ListPatientVisitsPayload{}, err
		}

		medsIds := make([]uint, 0, len(prescribedMeds))
		for _, pm := range prescribedMeds {
			medsIds = append(medsIds, pm.MedicineId)
		}

		meds, err := a.app.ListMedicinesByIds(medsIds)
		if err != nil {
			return ListPatientVisitsPayload{}, err
		}

		medsMapped := make(map[uint]models.Medicine)
		for _, med := range meds {
			medsMapped[med.Id] = med
		}

		outMeds := make([]PrescribedMedicine, 0, len(prescribedMeds))
		for _, pm := range prescribedMeds {
			outMed := new(PrescribedMedicine)
			outMed.FromModel(pm, medsMapped[pm.MedicineId])
			outMeds = append(outMeds, *outMed)
		}

		outVisits = append(outVisits, Visit{
			Id:                 visit.Id,
			Reason:             string(visit.Reason),
			ExtraNote:          visit.Notes,
			VisitedAt:          visit.CreatedAt,
			PrescribedMedicine: outMeds,
			PatientWeight:      visit.PatientWeight,
			PatientHeight:      visit.PatientHeight,
		})
	}

	return ListPatientVisitsPayload{
		Data: outVisits,
	}, nil
}
