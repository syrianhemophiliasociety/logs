package actions

import (
	"fmt"
	"net/http"
	"shs-web/log"
)

type RequestBloodTest struct {
	Id         uint     `json:"id"`
	Name       string   `json:"name"`
	FieldNames []string `json:"blood_test_field_name"`
	FieldUnits []string `json:"blood_test_field_unit"`
}

type RequestBloodTestSingle struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	FieldName string `json:"blood_test_field_name"`
	FieldUnit string `json:"blood_test_field_unit"`
}

type BloodTestField struct {
	Id       uint   `json:"id"`
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

type CreateBloodTestParams struct {
	RequestContext
	NewBloodTest       RequestBloodTest
	NewBloodTestSingle RequestBloodTestSingle
}

type CreateBloodTestPayload struct {
}

func (a *Actions) CreateBloodTest(params CreateBloodTestParams) (CreateBloodTestPayload, error) {
	var newBloodTest BloodTest

	if params.NewBloodTest.Name != "" {
		newBloodTest.Name = params.NewBloodTest.Name
		for i := range len(params.NewBloodTest.FieldNames) {
			newBloodTest.Fields = append(newBloodTest.Fields, BloodTestField{
				Name: params.NewBloodTest.FieldNames[i],
				Unit: params.NewBloodTest.FieldUnits[i],
			})
		}
	}
	if params.NewBloodTestSingle.Name != "" {
		newBloodTest.Name = params.NewBloodTestSingle.Name
		newBloodTest.Fields = append(newBloodTest.Fields, BloodTestField{
			Name: params.NewBloodTestSingle.FieldName,
			Unit: params.NewBloodTestSingle.FieldUnit,
		})
	}

	log.Warningln("body", newBloodTest)

	payload, err := makeRequest[map[string]any, CreateBloodTestPayload](makeRequestConfig[map[string]any]{
		method:   http.MethodPost,
		endpoint: "/v1/bloodtest",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: map[string]any{
			"new_blood_test": newBloodTest,
		},
	})
	if err != nil {
		return CreateBloodTestPayload{}, err
	}

	return payload, nil
}

type DeleteBloodTestParams struct {
	RequestContext
	BloodTestId uint
}

type DeleteBloodTestPayload struct {
}

func (a *Actions) DeleteBloodTest(params DeleteBloodTestParams) (DeleteBloodTestPayload, error) {
	payload, err := makeRequest[DeleteBloodTestParams, DeleteBloodTestPayload](makeRequestConfig[DeleteBloodTestParams]{
		method:   http.MethodDelete,
		endpoint: fmt.Sprintf("/v1/bloodtest/%d", params.BloodTestId),
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: params,
	})
	if err != nil {
		return DeleteBloodTestPayload{}, err
	}

	return payload, nil
}
