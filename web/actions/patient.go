package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//================
// Types
//================

type BloodTestFilledField struct {
	BloodTestFieldId uint   `json:"blood_test_field_id"`
	Name             string `json:"name"`
	Unit             string `json:"unit"`
	MinValue         uint   `json:"min_value"`
	MaxValue         uint   `json:"max_value"`
	ValueNumber      uint   `json:"value_number"`
	ValueString      string `json:"value_string"`
}

type BloodTestResult struct {
	BloodTestId  uint                   `json:"blood_test_id"`
	Name         string                 `json:"name"`
	FilledFields []BloodTestFilledField `json:"filled_fields"`
}

type Address struct {
	Id          uint   `json:"id"`
	Governorate string `json:"governorate"`
	Suburb      string `json:"suburb"`
	Street      string `json:"street"`
}

type Patient struct {
	Id           uint              `json:"id"`
	PublicId     string            `json:"public_id"`
	NationalId   string            `json:"national_id"`
	Nationality  string            `json:"nationality"`
	FirstName    string            `json:"first_name"`
	LastName     string            `json:"last_name"`
	FatherName   string            `json:"father_name"`
	MotherName   string            `json:"mother_name"`
	PlaceOfBirth Address           `json:"place_of_birth"`
	DateOfBirth  time.Time         `json:"date_of_birth"`
	Residency    Address           `json:"residency"`
	Gender       bool              `json:"gender"`
	PhoneNumber  string            `json:"phone_number"`
	BATScore     uint              `json:"bat_score"`
	Viri         []Virus           `json:"viruses"`
	BloodTests   []BloodTestResult `json:"blood_test_results"`
}

func (p Patient) FullName() string {
	return p.FirstName + " " + p.LastName
}

//==================
// Get patient by id
//==================

type GetPatientParams struct {
	RequestContext
	PatientId int
}

type GetPatientPayload struct {
	Data Patient `json:"data"`
}

func (a *Actions) GetPatient(params GetPatientParams) (Patient, error) {
	payload, err := makeRequest[any, GetPatientPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/patient/" + strconv.Itoa(params.PatientId),
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
	if err != nil {
		return Patient{}, err
	}

	return payload.Data, nil
}

//================
// Create patient
//================

type PatientRequest struct {
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	FatherName              string `json:"father_name"`
	MotherName              string `json:"mother_name"`
	Nationality             string `json:"nationality"`
	NationalId              string `json:"national_id"`
	Gender                  string `json:"gender"`
	PhoneNumber             string `json:"phone_number"`
	DateOfBirth             string `json:"date_of_birth"`
	PlaceOfBirthGovernorate string `json:"place_of_birth_governorate"`
	PlaceOfBirthSuburb      string `json:"place_of_birth_suburb"`
	PlaceOfBirthStreet      string `json:"place_of_birth_street"`
	ResidencyGovernorate    string `json:"residency_governorate"`
	ResidencySuburb         string `json:"residency_suburb"`
	ResidencyStreet         string `json:"residency_street"`
}

func (p PatientRequest) IntoPatient() Patient {
	dateOfBirth, _ := time.Parse("2006-01-02", p.DateOfBirth)
	return Patient{
		NationalId:  p.NationalId,
		Nationality: p.Nationality,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		FatherName:  p.FatherName,
		MotherName:  p.MotherName,
		PlaceOfBirth: Address{
			Governorate: p.PlaceOfBirthGovernorate,
			Suburb:      p.PlaceOfBirthSuburb,
			Street:      p.PlaceOfBirthStreet,
		},
		DateOfBirth: dateOfBirth,
		Residency: Address{
			Governorate: p.ResidencyGovernorate,
			Suburb:      p.ResidencySuburb,
			Street:      p.ResidencyStreet,
		},
		Gender:      p.Gender == "male",
		PhoneNumber: p.PhoneNumber,
		BATScore:    0,
	}
}

type CreatePatientParams struct {
	RequestContext
	NewPatient PatientRequest
}

type CreatePatientPayload struct{}

func (a *Actions) CreatePatient(params CreatePatientParams) (CreatePatientPayload, error) {
	payload, err := makeRequest[map[string]any, CreatePatientPayload](makeRequestConfig[map[string]any]{
		method:   http.MethodPost,
		endpoint: "/v1/patient",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: map[string]any{
			"new_patient": params.NewPatient.IntoPatient(),
		},
	})
	if err != nil {
		return CreatePatientPayload{}, err
	}

	return payload, nil
}

//================
// Find patient
//================

type FindPatientsParams struct {
	RequestContext
	PublicId     string  `json:"public_id"`
	NationalId   string  `json:"national_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	FatherName   string  `json:"father_name"`
	MotherName   string  `json:"mother_name"`
	PlaceOfBirth Address `json:"place_of_birth"`
	Residency    Address `json:"residency"`
	PhoneNumber  string  `json:"phone_number"`
}

type FindPatientsPayload struct {
	Data []Patient `json:"data"`
}

func (a *Actions) FindPatients(params FindPatientsParams) ([]Patient, error) {
	if params.FirstName == "" {
		params.FirstName = " "
	}
	if params.LastName == "" {
		params.LastName = " "
	}
	if params.FatherName == "" {
		params.FatherName = " "
	}
	if params.MotherName == "" {
		params.MotherName = " "
	}
	if params.NationalId == "" {
		params.NationalId = " "
	}
	if params.PhoneNumber == "" {
		params.PhoneNumber = " "
	}
	if params.PublicId == "" {
		params.PublicId = " "
	}

	payload, err := makeRequest[any, FindPatientsPayload](makeRequestConfig[any]{
		method: http.MethodGet,
		endpoint: fmt.Sprintf(
			"/v1/patients/public-id/%s/first-name/%s/last-name/%s/father-name/%s/mother-name/%s/national-id/%s/phone-number/%s",
			params.PublicId, params.FirstName, params.LastName, params.FatherName, params.MotherName, params.NationalId, params.PhoneNumber),
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
	if err != nil {
		return nil, err
	}

	return payload.Data, nil
}

//================================================
// Create patient non personal details
//================================================

type PatientBloodTests struct {
	BloodTests []BloodTestResult
}

func (p *PatientBloodTests) UnmarshalJSON(payload []byte) error {
	var data map[string]any
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}

	const bloodTestResultFieldValue = "blood_test_result_value#"

	getBloodTestMeta := func(key string) (name, fieldName string, id, fieldId int) {
		stuff := strings.Split(strings.TrimPrefix(key, bloodTestResultFieldValue), "#")
		id, _ = strconv.Atoi(stuff[0])
		fieldId, _ = strconv.Atoi(stuff[2])
		name = stuff[1]
		fieldName = stuff[3]

		return
	}

	bloodTestsFields := make(map[uint][]BloodTestFilledField)
	bloodTestNames := make(map[uint]string)
	for k, v := range data {
		if !strings.HasPrefix(k, bloodTestResultFieldValue) {
			continue
		}

		name, fieldName, id, fieldId := getBloodTestMeta(k)
		testResult, ok := v.(string)
		if !ok {
			continue
		}

		testResultInt, _ := strconv.Atoi(testResult)

		bloodTestNames[uint(id)] = name
		bloodTestsFields[uint(id)] = append(bloodTestsFields[uint(id)], BloodTestFilledField{
			BloodTestFieldId: uint(fieldId),
			Name:             fieldName,
			Unit:             "",
			MinValue:         0,
			MaxValue:         0,
			ValueNumber:      uint(testResultInt),
			ValueString:      testResult,
		})
	}

	for id, fields := range bloodTestsFields {
		p.BloodTests = append(p.BloodTests, BloodTestResult{
			BloodTestId:  id,
			Name:         bloodTestNames[id],
			FilledFields: fields,
		})
	}

	return nil
}

type CreatePatientBloodTestParams struct {
	RequestContext
	PatientId        string
	PatientBloodTest BloodTestResult
}

type CreatePatientBloodTestPayload struct {
}

func (a *Actions) CreatePatientBloodTest(params CreatePatientBloodTestParams) (CreatePatientBloodTestPayload, error) {
	return makeRequest[map[string]any, CreatePatientBloodTestPayload](makeRequestConfig[map[string]any]{
		method:   http.MethodPost,
		endpoint: "/v1/patient/bloodtest",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: map[string]any{
			"patient_id":         params.PatientId,
			"patient_blood_test": params.PatientBloodTest,
		},
	})
}

type PatientViruses struct {
	Viruses []Virus
}

func (p *PatientViruses) UnmarshalJSON(payload []byte) error {
	var data map[string]any
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}

	const virusPrefix = "virus-"

	viruses := make([]Virus, 0)
	for k, v := range data {
		if !strings.HasPrefix(k, virusPrefix) {
			continue
		}

		virusStr := strings.Split(strings.TrimPrefix(k, virusPrefix), "-")
		virusId, _ := strconv.Atoi(virusStr[0])
		if v == "on" {
			viruses = append(viruses, Virus{
				Id:   uint(virusId),
				Name: virusStr[1],
			})
		}
	}

	(*p).Viruses = viruses

	return nil
}
