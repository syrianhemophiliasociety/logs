package actions

import (
	"shs/app/models"
	"slices"
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

func (v *Visit) FromModel(visit models.Visit) {
	(*v) = Visit{
		Id:            visit.Id,
		Reason:        string(visit.Reason),
		ExtraNote:     visit.Notes,
		VisitedAt:     visit.CreatedAt,
		PatientWeight: visit.PatientWeight,
		PatientHeight: visit.PatientHeight,
	}
}

type VisitWithPatient struct {
	Visit   Visit   `json:"visit"`
	Patient Patient `json:"patient"`
}

type TreatmentDetails struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	ArabicTitle string `json:"arabic_title"`
	Type        string `json:"type"`
}

func (td *TreatmentDetails) FromModel(treatment models.TreatmentDetails) {
	(*td) = TreatmentDetails{
		Id:          treatment.Id,
		Title:       treatment.Title,
		ArabicTitle: treatment.ArabicTitle,
		Type:        treatment.Type,
	}
}

func (td *TreatmentDetails) IntoModel() models.TreatmentDetails {
	return models.TreatmentDetails{
		Title:       td.Title,
		ArabicTitle: td.ArabicTitle,
		Type:        td.Type,
	}
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

	patient, err := a.app.GetPatientByPublicId(params.PatientId)
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
	Medicine             Medicine
	TreatmentDetails     TreatmentDetails `json:"-"`
	TreatmentDetailsId   uint             `json:"treatment_details_id"`
	PrescribedMedicineId uint             `json:"prescribed_medicine_id"`
	UsedAt               time.Time        `json:"used_at"`
}

func (pm *PrescribedMedicine) FromModel(m models.PrescribedMedicine, med models.Medicine) {
	outMed := new(Medicine)
	outMed.FromModel(med)
	(*pm).Medicine = *outMed
	(*pm).PrescribedMedicineId = m.Id
	(*pm).UsedAt = m.UsedAt
	(*pm).TreatmentDetailsId = m.TreatmentDetailsId
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
	VisitId             uint                 `json:"visit_id"`
	Patient             Patient              `json:"patient"`
	VisitedAt           time.Time            `json:"visited_at"`
	PatientWeight       float64              `json:"patient_weight"`
	PatientHeight       float64              `json:"patient_height"`
	PrescribedMedicine  []PrescribedMedicine `json:"prescribed_medicine"`
	AvailableTreatments []TreatmentDetails   `json:"available_treatments"`
}

func (a *Actions) GetPatientLastVisit(params GetPatientLastVisitParams) (GetPatientLastVisitPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadOwnVisit) {
		return GetPatientLastVisitPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetPatientByPublicId(params.Account.Username)
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

	treatments, err := a.app.ListAllTreatmentDetails()
	if err != nil {
		return GetPatientLastVisitPayload{}, err
	}

	outTreatments := make([]TreatmentDetails, 0, len(treatments))
	for _, t := range treatments {
		outTreatment := new(TreatmentDetails)
		outTreatment.FromModel(t)
		outTreatments = append(outTreatments, *outTreatment)
	}

	return GetPatientLastVisitPayload{
		Patient:             *outPatient,
		PrescribedMedicine:  outMeds,
		VisitedAt:           lastVisit.CreatedAt,
		VisitId:             lastVisit.Id,
		PatientWeight:       lastVisit.PatientWeight,
		PatientHeight:       lastVisit.PatientHeight,
		AvailableTreatments: outTreatments,
	}, nil
}

type UseMedicineForVisitParams struct {
	ActionContext
	TreatmentId          uint `json:"visit_treatment_id"`
	PrescribedMedicineId uint `json:"prescribed_medicine_id"`
	VisitId              uint `json:"visit_id"`
}

type UseMedicineForVisitPayload struct {
}

func (a *Actions) UseMedicineForVisit(params UseMedicineForVisitParams) (UseMedicineForVisitPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteOwnVisit) {
		return UseMedicineForVisitPayload{}, ErrPermissionDenied{}
	}

	err := a.app.UseMedicineForVisit(params.PrescribedMedicineId, params.VisitId, params.TreatmentId)
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

	patient, err := a.app.GetPatientByPublicId(params.PatientId)
	if err != nil {
		return ListPatientVisitsPayload{}, err
	}

	visits, err := a.app.ListPatientVisits(patient.Id)
	if err != nil {
		return ListPatientVisitsPayload{}, err
	}

	treatments, _ := a.app.ListAllTreatmentDetails()
	treatmentsMapped := make(map[uint]TreatmentDetails, len(treatments))
	for _, t := range treatments {
		outTreatment := new(TreatmentDetails)
		outTreatment.FromModel(t)
		treatmentsMapped[t.Id] = *outTreatment
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
			outMed.TreatmentDetails = treatmentsMapped[pm.TreatmentDetailsId]
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

type ListAllVisitsParams struct {
	ActionContext
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	SortByVisitReason string    `json:"sort_by_visit_reason"`
}

type ListAllVisitsPayload struct {
	Data []VisitWithPatient `json:"data"`
}

func (a *Actions) ListAllVisits(params ListAllVisitsParams) (ListAllVisitsPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadOtherVisits) {
		return ListAllVisitsPayload{}, ErrPermissionDenied{}
	}

	visits, err := a.app.ListVisitsOnTimeRange(params.StartDate, time.Now())
	if err != nil {
		return ListAllVisitsPayload{}, err
	}

	outVisits := make([]VisitWithPatient, 0, len(visits))

	for _, visit := range visits {
		patient, err := a.app.GetPatientById(visit.PatientId)
		if err != nil {
			return ListAllVisitsPayload{}, err
		}

		prescribedMeds, err := a.app.ListPatientVisitPrescribedMedicine(visit.Id)
		if err != nil {
			return ListAllVisitsPayload{}, err
		}

		medsIds := make([]uint, 0, len(prescribedMeds))
		for _, pm := range prescribedMeds {
			medsIds = append(medsIds, pm.MedicineId)
		}

		meds, err := a.app.ListMedicinesByIds(medsIds)
		if err != nil {
			return ListAllVisitsPayload{}, err
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

		outVisits = append(outVisits, VisitWithPatient{
			Visit: Visit{
				Id:                 visit.Id,
				Reason:             string(visit.Reason),
				ExtraNote:          visit.Notes,
				VisitedAt:          visit.CreatedAt,
				PrescribedMedicine: outMeds,
				PatientWeight:      visit.PatientWeight,
				PatientHeight:      visit.PatientHeight,
			},
			Patient: *outPatient,
		})
	}

	if params.SortByVisitReason != "" {
		slices.SortFunc(outVisits, func(vi, vj VisitWithPatient) int {
			iMatches := vi.Visit.Reason == params.SortByVisitReason
			jMatches := vj.Visit.Reason == params.SortByVisitReason
			if iMatches == jMatches {
				return 0
			}

			if iMatches {
				return -1
			}

			return 1
		})
	}

	return ListAllVisitsPayload{
		Data: outVisits,
	}, nil
}

type CreateTreatmentDetailsParams struct {
	ActionContext
	TreatmentDetails TreatmentDetails
}

type CreateTreatmentDetailsPayload struct{}

func (a *Actions) CreateTreatmentDetails(params CreateTreatmentDetailsParams) (CreateTreatmentDetailsPayload, error) {
	// TODO: maybe create a new permission type
	if !params.Account.HasPermission(models.AccountPermissionWriteOtherVisits) {
		return CreateTreatmentDetailsPayload{}, ErrPermissionDenied{}
	}

	treatment := params.TreatmentDetails.IntoModel()
	_, err := a.app.CreateTreatmentDetails(treatment)
	if err != nil {
		return CreateTreatmentDetailsPayload{}, err
	}

	return CreateTreatmentDetailsPayload{}, nil
}

type ListAllTreatmentDetailsParams struct {
	ActionContext
}

type ListAllTreatmentDetailsPayload struct {
	Data []TreatmentDetails `json:"data"`
}

func (a *Actions) ListAllTreatmentDetails(params ListAllTreatmentDetailsParams) (ListAllTreatmentDetailsPayload, error) {
	// TODO: maybe create a new permission type
	if !params.Account.HasPermission(models.AccountPermissionReadOtherVisits) {
		return ListAllTreatmentDetailsPayload{}, ErrPermissionDenied{}
	}

	treatments, err := a.app.ListAllTreatmentDetails()
	if err != nil {
		return ListAllTreatmentDetailsPayload{}, err
	}

	outTreatments := make([]TreatmentDetails, 0, len(treatments))
	for _, t := range treatments {
		outTreatment := new(TreatmentDetails)
		outTreatment.FromModel(t)
		outTreatments = append(outTreatments, *outTreatment)
	}

	return ListAllTreatmentDetailsPayload{
		Data: outTreatments,
	}, nil
}

type DeleteTreatmentDetailsParams struct {
	ActionContext
	Id uint
}

type DeleteTreatmentDetailsPayload struct{}

func (a *Actions) DeleteTreatmentDetails(params DeleteTreatmentDetailsParams) (DeleteTreatmentDetailsPayload, error) {
	// TODO: maybe create a new permission type
	if !params.Account.HasPermission(models.AccountPermissionWriteOtherVisits) {
		return DeleteTreatmentDetailsPayload{}, ErrPermissionDenied{}
	}

	err := a.app.DeleteTreatmentDetails(params.Id)
	if err != nil {
		return DeleteTreatmentDetailsPayload{}, err
	}

	return DeleteTreatmentDetailsPayload{}, nil
}
