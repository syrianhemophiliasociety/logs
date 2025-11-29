package actions

import "shs/app/models"

type CreatePatientVisitParams struct {
	ActionContext
	PatientId             string
	VisitReason           string `json:"visit_reason"`
	PrescribedMedicineIds []uint `json:"prescribed_medicine_ids"`
}

type CreatePatientVisitPayload struct {
}

func (a *Actions) CreatePatientVisit(params CreatePatientVisitParams) (CreatePatientVisitPayload, error) {
	err := params.Account.CheckType(models.AccountTypeAdmin, models.AccountTypeSecritary)
	if err != nil {
		return CreatePatientVisitPayload{}, err
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

	_, err = a.app.CreatePatientVisit(models.Visit{
		PatientId:           patient.Id,
		Reason:              models.VisitReason(params.VisitReason),
		PrescribedMedicines: meds,
	})
	if err != nil {
		return CreatePatientVisitPayload{}, err
	}

	return CreatePatientVisitPayload{}, nil
}
