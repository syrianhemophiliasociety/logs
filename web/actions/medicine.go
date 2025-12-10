package actions

import (
	"fmt"
	"net/http"
	"strconv"
)

type Medicine struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Dose   int    `json:"dose"`
	Unit   string `json:"unit"`
	Amount int    `json:"amount"`
}

type ListAllMedicinesParams struct {
	RequestContext
}

type ListAllMedicinesPayload struct {
	Data []Medicine `json:"data"`
}

func (a *Actions) ListAllMedicines(params ListAllMedicinesParams) ([]Medicine, error) {
	payload, err := makeRequest[any, ListAllMedicinesPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/medicine/all",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
	if err != nil {
		return nil, err
	}

	return payload.Data, nil
}

type RequestMedicine struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Dose   string `json:"dose"`
	Unit   string `json:"unit"`
	Amount string `json:"amount"`
}

type CreateMedicineParams struct {
	RequestContext
	NewMedicine RequestMedicine `json:"new_medicine"`
}

type CreateMedicinePayload struct {
}

func (a *Actions) CreateMedicine(params CreateMedicineParams) (CreateMedicinePayload, error) {
	dose, err := strconv.Atoi(params.NewMedicine.Dose)
	if err != nil {
		return CreateMedicinePayload{}, err
	}

	amount, err := strconv.Atoi(params.NewMedicine.Amount)
	if err != nil {
		return CreateMedicinePayload{}, err
	}

	medicine := Medicine{
		Name:   params.NewMedicine.Name,
		Dose:   dose,
		Unit:   params.NewMedicine.Unit,
		Amount: amount,
	}

	payload, err := makeRequest[map[string]any, CreateMedicinePayload](makeRequestConfig[map[string]any]{
		method:   http.MethodPost,
		endpoint: "/v1/medicine",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: map[string]any{
			"new_medicine": medicine,
		},
	})
	if err != nil {
		return CreateMedicinePayload{}, err
	}

	return payload, nil
}

type DeleteMedicineParams struct {
	RequestContext
	MedicineId uint
}

type DeleteMedicinePayload struct {
}

func (a *Actions) DeleteMedicine(params DeleteMedicineParams) (DeleteMedicinePayload, error) {
	payload, err := makeRequest[DeleteMedicineParams, DeleteMedicinePayload](makeRequestConfig[DeleteMedicineParams]{
		method:   http.MethodDelete,
		endpoint: fmt.Sprintf("/v1/medicine/%d", params.MedicineId),
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: params,
	})
	if err != nil {
		return DeleteMedicinePayload{}, err
	}

	return payload, nil
}
