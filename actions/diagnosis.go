package actions

import (
	"shs/app/models"
	"time"
)

type Diagnosis struct {
	Id        uint   `json:"id"`
	GroupName string `json:"group_name"`
	Title     string `json:"title"`

	CreatedAt time.Time `json:"created_at"`
}

func (d Diagnosis) IntoModel() models.Diagnosis {
	return models.Diagnosis{
		GroupName: d.GroupName,
		Title:     d.Title,
	}
}

func (d *Diagnosis) FromModel(diagnosis models.Diagnosis) {
	(*d).Id = diagnosis.Id
	(*d).GroupName = diagnosis.GroupName
	(*d).Title = diagnosis.Title
	(*d).CreatedAt = diagnosis.CreatedAt
}

type CreateDiagnosisParams struct {
	ActionContext
	Diagnosis Diagnosis `json:"new_diagnosis"`
}

type CreateDiagnosisPayload struct {
}

func (a *Actions) CreateDiagnosis(params CreateDiagnosisParams) (CreateDiagnosisPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteDiagnoses) {
		return CreateDiagnosisPayload{}, ErrPermissionDenied{}
	}

	_, err := a.app.CreateDiagnosis(params.Diagnosis.IntoModel())
	if err != nil {
		return CreateDiagnosisPayload{}, err
	}

	return CreateDiagnosisPayload{}, nil
}

type ListAllDiagnosesParams struct {
	ActionContext
}

type ListAllDiagnosesPayload struct {
	Data []Diagnosis `json:"data"`
}

func (a *Actions) ListAllDiagnoses(params ListAllDiagnosesParams) (ListAllDiagnosesPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadDiagnoses) {
		return ListAllDiagnosesPayload{}, ErrPermissionDenied{}
	}

	diagnoses, err := a.app.ListAllDiagnoses()
	if err != nil {
		return ListAllDiagnosesPayload{}, err
	}

	outDiagnoses := make([]Diagnosis, 0, len(diagnoses))
	for _, d := range diagnoses {
		outDiagnosis := new(Diagnosis)
		outDiagnosis.FromModel(d)
		outDiagnoses = append(outDiagnoses, *outDiagnosis)
	}

	return ListAllDiagnosesPayload{
		Data: outDiagnoses,
	}, nil
}

type DeleteDiagnosisParams struct {
	ActionContext
	DiagnosisId uint `json:"blood_test_id"`
}

type DeleteDiagnosisPayload struct {
}

func (a *Actions) DeleteDiagnosis(params DeleteDiagnosisParams) (DeleteDiagnosisPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteDiagnoses) {
		return DeleteDiagnosisPayload{}, ErrPermissionDenied{}
	}

	err := a.app.DeleteDiagnosis(params.DiagnosisId)
	if err != nil {
		return DeleteDiagnosisPayload{}, err
	}

	return DeleteDiagnosisPayload{}, nil
}
