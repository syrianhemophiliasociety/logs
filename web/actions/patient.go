package actions

import "net/http"

type Patient struct {
}

type ListAllPatientsParams struct {
	RequestContext
}

type ListAllPatientsPayload struct {
	Data []Patient `json:"data"`
}

func (a *Actions) ListAllPatients(params ListAllPatientsParams) ([]Patient, error) {
	return nil, nil
	payload, err := makeRequest[any, ListAllPatientsPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/patient/all",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
	if err != nil {
		return nil, err
	}

	return payload.Data, nil
}
