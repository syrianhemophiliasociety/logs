package actions

import (
	"fmt"
	"net/http"
	"time"
)

type Diagnosis struct {
	Id        uint   `json:"id"`
	GroupName string `json:"group_name"`
	Title     string `json:"title"`

	CreatedAt time.Time `json:"created_at"`
}

type CreateDiagnosisParams struct {
	RequestContext
	NewDiagnosis Diagnosis
}

type CreateDiagnosisPayload struct {
}

func (a *Actions) CreateDiagnosis(params CreateDiagnosisParams) (CreateDiagnosisPayload, error) {
	return makeRequest[map[string]any, CreateDiagnosisPayload](makeRequestConfig[map[string]any]{
		method:   http.MethodPost,
		endpoint: "/v1/diagnosis",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: map[string]any{
			"new_diagnosis": params.NewDiagnosis,
		},
	})
}

type ListAllDiagnosesParams struct {
	RequestContext
}

type ListAllDiagnosesPayload struct {
	Data []Diagnosis `json:"data"`
}

func (a *Actions) ListAllDiagnoses(params ListAllDiagnosesParams) (ListAllDiagnosesPayload, error) {
	return makeRequest[any, ListAllDiagnosesPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/diagnosis/all",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
}

type DeleteDiagnosisParams struct {
	RequestContext
	DiagnosisId uint
}

type DeleteDiagnosisPayload struct{}

func (a *Actions) DeleteDiagnosis(params DeleteDiagnosisParams) (DeleteDiagnosisPayload, error) {
	return makeRequest[any, DeleteDiagnosisPayload](makeRequestConfig[any]{
		method:   http.MethodDelete,
		endpoint: fmt.Sprintf("/v1/diagnosis/%d", params.DiagnosisId),
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: params,
	})
}
