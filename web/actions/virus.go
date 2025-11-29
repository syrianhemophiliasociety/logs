package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Virus struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type ListAllVirusesParams struct {
	RequestContext
}

type ListAllVirusesPayload struct {
	Data []Virus `json:"data"`
}

func (a *Actions) ListAllViruses(params ListAllVirusesParams) ([]Virus, error) {
	payload, err := makeRequest[any, ListAllVirusesPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/virus/all",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
	if err != nil {
		return nil, err
	}

	return payload.Data, nil
}

type CreateVirusRequest struct {
	Name         string `json:"name"`
	BloodTestIds []uint `json:"blood_test_ids"`
}

func (v *CreateVirusRequest) UnmarshalJSON(payload []byte) error {
	var data map[string]any
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}

	var ok bool
	(*v).Name, ok = data["name"].(string)
	if !ok {
		return errors.New("missing name")
	}

	const bloodTestKey = "blood_test_id"
	switch data[bloodTestKey].(type) {
	case string:
		btIdInt, err := strconv.Atoi(data[bloodTestKey].(string))
		if err != nil {
			return err
		}
		(*v).BloodTestIds = []uint{uint(btIdInt)}

	case []any:
		for _, btId := range data[bloodTestKey].([]any) {
			btIdStr, ok := btId.(string)
			if !ok {
				return errors.New("invalid blood_test_id type")
			}
			btIdInt, err := strconv.Atoi(btIdStr)
			if err != nil {
				return err
			}
			(*v).BloodTestIds = append((*v).BloodTestIds, uint(btIdInt))
		}

	default:
		return errors.New("invalid blood_test_id value")
	}

	return nil
}

type CreateVirusParams struct {
	RequestContext
	NewVirus CreateVirusRequest `json:"new_virus"`
}

type CreateVirusPayload struct {
}

func (a *Actions) CreateVirus(params CreateVirusParams) (CreateVirusPayload, error) {
	payload, err := makeRequest[CreateVirusParams, CreateVirusPayload](makeRequestConfig[CreateVirusParams]{
		method:   http.MethodPost,
		endpoint: "/v1/virus",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: params,
	})
	if err != nil {
		return CreateVirusPayload{}, err
	}

	return payload, nil
}

type DeleteVirusParams struct {
	RequestContext
	VirusId uint
}

type DeleteVirusPayload struct {
}

func (a *Actions) DeleteVirus(params DeleteVirusParams) (DeleteVirusPayload, error) {
	payload, err := makeRequest[DeleteVirusParams, DeleteVirusPayload](makeRequestConfig[DeleteVirusParams]{
		method:   http.MethodDelete,
		endpoint: fmt.Sprintf("/v1/virus/%d", params.VirusId),
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: params,
	})
	if err != nil {
		return DeleteVirusPayload{}, err
	}

	return payload, nil
}
