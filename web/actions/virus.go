package actions

import (
	"fmt"
	"net/http"
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

type CreateVirusParams struct {
	RequestContext
	NewVirus Virus `json:"new_virus"`
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
