package actions

import "net/http"

type RequestBloodTestField struct {
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	MinValue string `json:"min_value"`
	MaxValue string `json:"max_value"`
}

type BloodTestField struct {
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	MinValue uint   `json:"min_value"`
	MaxValue uint   `json:"max_value"`
}

type BloodTest struct {
	Id     uint             `json:"id"`
	Name   string           `json:"name"`
	Fields []BloodTestField `json:"fields"`
}

type ListAllBloodTestsParams struct {
	RequestContext
}

type ListAllBloodTestsPayload struct {
	Data []BloodTest `json:"data"`
}

func (a *Actions) ListAllBloodTests(params ListAllBloodTestsParams) ([]BloodTest, error) {
	payload, err := makeRequest[any, ListAllBloodTestsPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/bloodtest/all",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
	if err != nil {
		return nil, err
	}

	return payload.Data, nil
}
