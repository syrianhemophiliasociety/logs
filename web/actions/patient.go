package actions

import (
	"encoding/json"
	"errors"
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
	ValueNumber      uint   `json:"value_number"`
	ValueString      string `json:"value_string"`
}

type BloodTestResult struct {
	BloodTestId  uint                   `json:"blood_test_id"`
	Name         string                 `json:"name"`
	FilledFields []BloodTestFilledField `json:"filled_fields"`
	Pending      bool                   `json:"pending"`
}

type Address struct {
	Id          uint   `json:"id"`
	Governorate string `json:"governorate"`
	Suburb      string `json:"suburb"`
	Street      string `json:"street"`
}

type Patient struct {
	Id                  uint              `json:"id"`
	PublicId            string            `json:"public_id"`
	NationalId          string            `json:"national_id"`
	Nationality         string            `json:"nationality"`
	FirstName           string            `json:"first_name"`
	LastName            string            `json:"last_name"`
	FatherName          string            `json:"father_name"`
	MotherName          string            `json:"mother_name"`
	PlaceOfBirth        Address           `json:"place_of_birth"`
	DateOfBirth         time.Time         `json:"date_of_birth"`
	Residency           Address           `json:"residency"`
	Gender              bool              `json:"gender"`
	PhoneNumber         string            `json:"phone_number"`
	BATScore            uint              `json:"bat_score"`
	Viri                []Virus           `json:"viruses"`
	BloodTests          []BloodTestResult `json:"blood_test_results"`
	FamilyHistoryExists bool              `json:"family_history_exists"`
	FirstVisitReason    string            `json:"first_visit_reason"`
}

func (p Patient) FullName() string {
	return p.FirstName + " " + p.LastName
}

type PrescribedMedicine struct {
	Medicine
	PrescribedMedicineId uint      `json:"prescribed_medicine_id"`
	UsedAt               time.Time `json:"used_at"`
}

type Visit struct {
	Reason             string               `json:"reason"`
	VisitedAt          time.Time            `json:"visited_at"`
	PrescribedMedicine []PrescribedMedicine `json:"prescribed_medicine"`
}

//==================
// Get patient by id
//==================

type GetPatientParams struct {
	RequestContext
	PatientId string
}

type GetPatientPayload struct {
	Data Patient `json:"data"`
}

func (a *Actions) GetPatient(params GetPatientParams) (Patient, error) {
	payload, err := makeRequest[any, GetPatientPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/patient/" + params.PatientId,
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
	FamilyHistoryExists     string `json:"family_history_exists"`
	FirstVisitReason        string `json:"first_visit_reason"`
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
		Gender:              p.Gender == "male",
		PhoneNumber:         p.PhoneNumber,
		BATScore:            0,
		FirstVisitReason:    p.FirstVisitReason,
		FamilyHistoryExists: p.FamilyHistoryExists == "on",
	}
}

type CreatePatientParams struct {
	RequestContext
	NewPatient PatientRequest
}

type CreatePatientPayload struct {
	Id string `json:"id"`
}

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

//================
// List last patient
//================

type ListLastPatientsParams struct {
	RequestContext
}

type ListLastPatientsPayload struct {
	Data []Patient `json:"data"`
}

func (a *Actions) ListLastPatients(params ListLastPatientsParams) ([]Patient, error) {
	payload, err := makeRequest[any, ListLastPatientsPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/patients/last",
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
			ValueNumber:      uint(testResultInt),
			ValueString:      testResult,
		})
	}

	doTestLater, _ := data["do_later"].(string)

	for id, fields := range bloodTestsFields {
		(*p).BloodTests = append((*p).BloodTests, BloodTestResult{
			BloodTestId:  id,
			Name:         bloodTestNames[id],
			FilledFields: fields,
			Pending:      doTestLater == "on",
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

//================
// Delete patient
//================

type DeletePatientParams struct {
	RequestContext
	PatientId string
}

type DeletePatientPayload struct {
}

func (a *Actions) DeletePatient(params DeletePatientParams) (DeletePatientPayload, error) {
	return makeRequest[any, DeletePatientPayload](makeRequestConfig[any]{
		method:   http.MethodDelete,
		endpoint: "/v1/patient/" + params.PatientId,
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
}

//================================
// Check-up visits
//================================

type CreateCheckUpRequest struct {
	VisitReason         string
	VisitExtraDetails   string
	PrescribedMedicines []Medicine
}

func (v *CreateCheckUpRequest) UnmarshalJSON(payload []byte) error {
	var data map[string]any
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}

	var ok bool
	(*v).VisitReason, ok = data["visit_reason"].(string)
	if !ok {
		return errors.New("missing visit_reason")
	}
	(*v).VisitExtraDetails, _ = data["visit_extra_details"].(string)

	const medicineIdsKey = "medicine_id"
	switch data[medicineIdsKey].(type) {
	case string:
		mIdInt, err := strconv.Atoi(data[medicineIdsKey].(string))
		if err != nil {
			return err
		}
		(*v).PrescribedMedicines = []Medicine{
			{
				Id: uint(mIdInt),
			},
		}

	case []any:
		for _, mId := range data[medicineIdsKey].([]any) {
			mIdStr, ok := mId.(string)
			if !ok {
				return errors.New("invalid medicine_id type")
			}
			mIdInt, err := strconv.Atoi(mIdStr)
			if err != nil {
				return err
			}
			(*v).PrescribedMedicines = append((*v).PrescribedMedicines, Medicine{
				Id: uint(mIdInt),
			})
		}
	}

	const medicineAmountKey = "amount"
	switch data[medicineAmountKey].(type) {
	case string:
		mAmountInt, err := strconv.Atoi(data[medicineAmountKey].(string))
		if err != nil {
			return err
		}
		(*v).PrescribedMedicines[0].Amount = mAmountInt

	case []any:
		for i, mId := range data[medicineAmountKey].([]any) {
			mIdStr, ok := mId.(string)
			if !ok {
				return errors.New("invalid amount type")
			}
			mAmountInt, err := strconv.Atoi(mIdStr)
			if err != nil {
				return err
			}
			(*v).PrescribedMedicines[i].Amount = mAmountInt
		}
	}

	return nil
}

type CreatePatientCheckUpParams struct {
	RequestContext
	PatientId      string
	CheckUpRequest CreateCheckUpRequest
}

type CreatePatientCheckUpPayload struct {
}

func (a *Actions) CreatePatientCheckUp(params CreatePatientCheckUpParams) (CreatePatientCheckUpPayload, error) {
	return makeRequest[map[string]any, CreatePatientCheckUpPayload](makeRequestConfig[map[string]any]{
		method:   http.MethodPost,
		endpoint: "/v1/patient/" + params.PatientId + "/checkup",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
		body: map[string]any{
			"visit_reason":         params.CheckUpRequest.VisitReason,
			"visit_extra_details":  params.CheckUpRequest.VisitExtraDetails,
			"prescribed_medicines": params.CheckUpRequest.PrescribedMedicines,
		},
	})
}

//================
// Patient Card
//================

type GeneratePatientCardParams struct {
	RequestContext
	PatientId string
}

type GeneratePatientCardPayload struct {
	ImageBase64 string `json:"image_base_64"`
}

func (a *Actions) GeneratePatientCard(params GeneratePatientCardParams) (GeneratePatientCardPayload, error) {
	return makeRequest[any, GeneratePatientCardPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/patient/" + params.PatientId + "/card",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
}

//================================
// Patient Medications
//================================

type GetPatientLastVisitParams struct {
	RequestContext
}

type GetPatientLastVisitPayload struct {
	Patient            Patient              `json:"patient"`
	VisitedAt          time.Time            `json:"visited_at"`
	PrescribedMedicine []PrescribedMedicine `json:"prescribed_medicine"`
}

func (a *Actions) GetPatientLastVisit(params GetPatientLastVisitParams) (GetPatientLastVisitPayload, error) {
	return makeRequest[any, GetPatientLastVisitPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/patient/last-visit",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})
}

//================================
// Patient Visits
//================================

type ListPatientVisitsParams struct {
	RequestContext
	PatientId string
}

type ListPatientVisitsPayload struct {
	Data []Visit `json:"data"`
}

func (a *Actions) ListPatientVisits(params ListPatientVisitsParams) ([]Visit, error) {
	payload, err := makeRequest[any, ListPatientVisitsPayload](makeRequestConfig[any]{
		method:   http.MethodGet,
		endpoint: "/v1/patient/" + params.PatientId + "/visits",
		headers: map[string]string{
			"Authorization": params.SessionToken,
		},
	})

	if err != nil {
		return nil, err
	}

	return payload.Data, nil
}
