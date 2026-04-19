package actions

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

//================================
// Import Patient
//================================

type ImportPatientsFromCsvParams struct {
	RequestContext
	CsvFile io.Reader
}

type ImportPatientsFromCsvPayload struct {
	ImportCount     int       `json:"import_count"`
	IgnoredPatients []Patient `json:"ignored_patients"`
}

func (a *Actions) ImportPatientsFromCsv(params ImportPatientsFromCsvParams) (ImportPatientsFromCsvPayload, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("patient_records", "patient_records.csv")
	if err != nil {
		return ImportPatientsFromCsvPayload{}, err
	}

	_, err = io.Copy(part, params.CsvFile)
	if err != nil {
		return ImportPatientsFromCsvPayload{}, err
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, getRequestUrl("/v1/patients/import/csv"), body)
	if err != nil {
		return ImportPatientsFromCsvPayload{}, err
	}
	req.Header.Set("Authorization", params.SessionToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ImportPatientsFromCsvPayload{}, err
	}
	defer resp.Body.Close()

	var respBody ImportPatientsFromCsvPayload
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return ImportPatientsFromCsvPayload{}, err
	}

	return respBody, nil
}
