package actions

import (
	"fmt"
	"net/http"
	"strconv"
)

type RequestBloodTest struct {
	Id         uint     `json:"id"`
	Name       string   `json:"name"`
	FieldNames []string `json:"blood_test_field_name"`
	FieldUnits []string `json:"blood_test_field_unit"`
	MinValues  []string `json:"blood_test_field_min_value"`
	MaxValues  []string `json:"blood_test_field_max_value"`
}

type RequestBloodTestSingle struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	FieldName string `json:"blood_test_field_name"`
	FieldUnit string `json:"blood_test_field_unit"`
	MinValue  string `json:"blood_test_field_min_value"`
	MaxValue  string `json:"blood_test_field_max_value"`
}

type BloodTestField struct {
	Id             uint    `json:"id"`
	Name           string  `json:"name"`
	Unit           string  `json:"unit"`
	MinValueNumber float64 `json:"min_value_number"`
	MinValueString string  `json:"min_value_string"`
	MaxValueNumber float64 `json:"max_value_number"`
	MaxValueString string  `json:"max_value_string"`
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
			minValue, _ := strconv.ParseFloat(params.NewBloodTest.MinValues[i], 64)
			maxValue, _ := strconv.ParseFloat(params.NewBloodTest.MaxValues[i], 64)

			newBloodTest.Fields = append(newBloodTest.Fields, BloodTestField{
				Name:           params.NewBloodTest.FieldNames[i],
				Unit:           params.NewBloodTest.FieldUnits[i],
				MinValueString: params.NewBloodTest.MinValues[i],
				MinValueNumber: minValue,
				MaxValueString: params.NewBloodTest.MaxValues[i],
				MaxValueNumber: maxValue,
			})
		}
	}
	if params.NewBloodTestSingle.Name != "" {
		newBloodTest.Name = params.NewBloodTestSingle.Name
		minValue, _ := strconv.ParseFloat(params.NewBloodTestSingle.MinValue, 64)
		maxValue, _ := strconv.ParseFloat(params.NewBloodTestSingle.MaxValue, 64)
		newBloodTest.Fields = append(newBloodTest.Fields, BloodTestField{
			Name:           params.NewBloodTestSingle.FieldName,
			Unit:           params.NewBloodTestSingle.FieldUnit,
			MinValueString: params.NewBloodTestSingle.MinValue,
			MinValueNumber: minValue,
			MaxValueString: params.NewBloodTestSingle.MaxValue,
			MaxValueNumber: maxValue,
		})
	}

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
