package actions

import (
	"shs/app"
	"shs/app/models"
	"shs/log"
	"shs/nanoid"
	"slices"
	"strings"
	"time"
)

type BloodTestFilledField struct {
	BloodTestFieldId uint                 `json:"blood_test_field_id"`
	Name             string               `json:"name"`
	Unit             models.BlootTestUnit `json:"unit"`
	MinValue         uint                 `json:"min_value"`
	MaxValue         uint                 `json:"max_value"`
	ValueNumber      uint                 `json:"value_number"`
	ValueString      string               `json:"value_string"`
}

type BloodTestResult struct {
	Name         string                 `json:"name"`
	BloodTestId  uint                   `json:"blood_test_id"`
	FilledFields []BloodTestFilledField `json:"filled_fields"`
}

type Address struct {
	Id          uint   `json:"id"`
	Governorate string `json:"governorate"`
	Suburb      string `json:"suburb"`
	Street      string `json:"street"`
}

func (a Address) IntoModel() models.Address {
	return models.Address{
		Id:          a.Id,
		Governorate: a.Governorate,
		Suburb:      a.Suburb,
		Street:      a.Street,
	}
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

type CreatePatientParams struct {
	ActionContext
	NewPatient Patient `json:"new_patient"`
}

type CreatePatientPayload struct {
}

func (a *Actions) CreatePatient(params CreatePatientParams) (CreatePatientPayload, error) {
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
	if err != nil {
		return CreatePatientPayload{}, err
	}

	newPatient := models.Patient{
		PublicId:    nanoid.New(),
		NationalId:  params.NewPatient.NationalId,
		Nationality: params.NewPatient.Nationality,
		FirstName:   params.NewPatient.FirstName,
		LastName:    params.NewPatient.LastName,
		FatherName:  params.NewPatient.FatherName,
		MotherName:  params.NewPatient.MotherName,
		DateOfBirth: params.NewPatient.DateOfBirth,
		Gender:      params.NewPatient.Gender,
		PhoneNumber: params.NewPatient.PhoneNumber,
		BATScore:    params.NewPatient.BATScore,
		Viri:        []models.Virus{},
		BloodTests:  []models.BloodTestResult{},
	}

	residencyAddresses, _ := a.app.GetAllAddressesALike(models.Address{
		Governorate: params.NewPatient.Residency.Governorate,
		Suburb:      params.NewPatient.Residency.Suburb,
		Street:      params.NewPatient.Residency.Street,
	})

	if len(residencyAddresses) == 1 {
		newPatient.Residency.Id = residencyAddresses[0].Id
		newPatient.ResidencyId = residencyAddresses[0].Id
	} else {
		residency, err := a.app.CreateAddress(params.NewPatient.Residency.IntoModel())
		if err != nil {
			return CreatePatientPayload{}, err
		}
		newPatient.Residency = residency
	}

	placesOfBirth, _ := a.app.GetAllAddressesALike(models.Address{
		Governorate: params.NewPatient.PlaceOfBirth.Governorate,
		Suburb:      params.NewPatient.PlaceOfBirth.Suburb,
		Street:      params.NewPatient.PlaceOfBirth.Street,
	})

	if len(placesOfBirth) == 1 {
		newPatient.PlaceOfBirth.Id = placesOfBirth[0].Id
		newPatient.PlaceOfBirthId = placesOfBirth[0].Id
	} else {
		placeOfBirth, err := a.app.CreateAddress(params.NewPatient.PlaceOfBirth.IntoModel())
		if err != nil {
			return CreatePatientPayload{}, err
		}
		newPatient.PlaceOfBirth = placeOfBirth
	}

	allViri, err := a.app.ListAllViri()
	if err != nil {
		return CreatePatientPayload{}, err
	}

	for _, virus := range params.NewPatient.Viri {
		matchedVirusIndex := slices.IndexFunc(allViri, func(v models.Virus) bool {
			return v.Id == virus.Id
		})
		if matchedVirusIndex < 0 {
			continue
		}
		newPatient.Viri = append(newPatient.Viri, allViri[matchedVirusIndex])
	}

	newPatient, err = a.app.CreatePatient(newPatient)
	if err != nil {
		return CreatePatientPayload{}, err
	}

	for _, btr := range params.NewPatient.BloodTests {
		bloodTestResultFields := make([]models.BloodTestFilledField, 0, len(btr.FilledFields))
		for _, field := range btr.FilledFields {
			bloodTestResultFields = append(bloodTestResultFields, models.BloodTestFilledField{
				BloodTestResultId: btr.BloodTestId,
				BloodTestFieldId:  field.BloodTestFieldId,
				ValueNumber:       field.ValueNumber,
				ValueString:       field.ValueString,
			})
		}

		_, err := a.app.CreateBloodTestResult(models.BloodTestResult{
			BloodTestId:  btr.BloodTestId,
			PatientId:    newPatient.Id,
			FilledFields: bloodTestResultFields,
		})
		if err != nil {
			log.Errorf("failed creating blood test result: %v\n", err)
		}
	}

	_, err = a.app.CreateAccount(models.Account{
		DisplayName: newPatient.FirstName + " " + newPatient.LastName,
		Username:    newPatient.PublicId,
		Password:    newPatient.NationalId,
		Type:        models.AccountTypePatient,
	})
	if err != nil {
		return CreatePatientPayload{}, err
	}

	return CreatePatientPayload{}, nil
}

type CreatePatientBloodTestParams struct {
	ActionContext
	PatientPublicId string          `json:"patient_id"`
	BloodTest       BloodTestResult `json:"patient_blood_test"`
}

type CreatePatientBloodTestPayload struct {
}

func (a *Actions) CreatePatientBloodTest(params CreatePatientBloodTestParams) (CreatePatientBloodTestPayload, error) {
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
	if err != nil {
		return CreatePatientBloodTestPayload{}, err
	}

	patient, err := a.app.GetPatientByPublicId(params.PatientPublicId)
	if err != nil {
		return CreatePatientBloodTestPayload{}, err
	}

	bloodTestResultFields := make([]models.BloodTestFilledField, 0, len(params.BloodTest.FilledFields))
	for _, field := range params.BloodTest.FilledFields {
		bloodTestResultFields = append(bloodTestResultFields, models.BloodTestFilledField{
			BloodTestFieldId: field.BloodTestFieldId,
			ValueNumber:      field.ValueNumber,
			ValueString:      field.ValueString,
		})
	}

	_, err = a.app.CreateBloodTestResult(models.BloodTestResult{
		BloodTestId:  params.BloodTest.BloodTestId,
		PatientId:    patient.Id,
		FilledFields: bloodTestResultFields,
	})
	if err != nil {
		return CreatePatientBloodTestPayload{}, err
	}

	return CreatePatientBloodTestPayload{}, nil
}

type FindPatientsParams struct {
	ActionContext
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

func (p *FindPatientsParams) clean() {
	p.PublicId = strings.TrimSpace(p.PublicId)
	p.NationalId = strings.TrimSpace(p.NationalId)
	p.FirstName = strings.TrimSpace(p.FirstName)
	p.LastName = strings.TrimSpace(p.LastName)
	p.FatherName = strings.TrimSpace(p.FatherName)
	p.MotherName = strings.TrimSpace(p.MotherName)
	p.PhoneNumber = strings.TrimSpace(p.PhoneNumber)
}

func (p *FindPatientsParams) empty() bool {
	return p.PublicId == "" && p.NationalId == "" &&
		p.FirstName == "" && p.LastName == "" &&
		p.FatherName == "" && p.MotherName == "" &&
		p.PhoneNumber == ""
}

type FindPatientsPayload struct {
	Data []Patient `json:"data"`
}

func (a *Actions) FindPatients(params FindPatientsParams) (FindPatientsPayload, error) {
	params.clean()

	if params.empty() {
		return FindPatientsPayload{}, app.ErrNotFound{
			ResourceName: "patient",
		}
	}

	patients, err := a.app.FindPatientsByIndexFields(models.PatientIndexFields{
		PublicId:     params.PublicId,
		NationalId:   params.NationalId,
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		FatherName:   params.FatherName,
		MotherName:   params.MotherName,
		PlaceOfBirth: models.Address{},
		Residency:    models.Address{},
		PhoneNumber:  params.PhoneNumber,
	})
	if err != nil {
		return FindPatientsPayload{}, err
	}

	outPatients := make([]Patient, 0, len(patients))
	for _, patient := range patients {
		outPatients = append(outPatients, Patient{
			Id:          patient.Id,
			PublicId:    patient.PublicId,
			NationalId:  patient.NationalId,
			Nationality: patient.NationalId,
			FirstName:   patient.FirstName,
			LastName:    patient.LastName,
			FatherName:  patient.FatherName,
			MotherName:  patient.MotherName,
			PlaceOfBirth: Address{
				Id:          patient.PlaceOfBirth.Id,
				Governorate: patient.PlaceOfBirth.Governorate,
				Suburb:      patient.PlaceOfBirth.Suburb,
				Street:      patient.PlaceOfBirth.Street,
			},
			DateOfBirth: patient.DateOfBirth,
			Residency: Address{
				Id:          patient.Residency.Id,
				Governorate: patient.Residency.Governorate,
				Suburb:      patient.Residency.Suburb,
				Street:      patient.Residency.Street,
			},
			Gender:      patient.Gender,
			PhoneNumber: patient.PhoneNumber,
			BATScore:    patient.BATScore,
		})
	}

	return FindPatientsPayload{
		Data: outPatients,
	}, nil
}

type GetPatientParams struct {
	ActionContext
	PublicId string
}

type GetPatientPayload struct {
	Data Patient `json:"data"`
}

func (a *Actions) GetPatient(params GetPatientParams) (GetPatientPayload, error) {
	err := checkAccountType(params.Account, models.AccountTypeSuperAdmin, models.AccountTypeAdmin, models.AccountTypeSecritary)
	if err != nil {
		return GetPatientPayload{}, err
	}

	patient, err := a.app.GetPatientByPublicId(params.PublicId)
	if err != nil {
		return GetPatientPayload{}, err
	}

	viruses := make([]Virus, 0, len(patient.Viri))
	for _, v := range patient.Viri {
		viruses = append(viruses, Virus{
			Id:   v.Id,
			Name: v.Name,
		})
	}

	bloodTests, err := a.app.ListAllBloodTests()
	if err != nil {
		return GetPatientPayload{}, err
	}

	bloodTestNames := make(map[uint]string)
	bloodTestFieldNames := make(map[uint]string)
	bloodTestFieldUnits := make(map[uint]models.BlootTestUnit)

	for _, bt := range bloodTests {
		bloodTestNames[bt.Id] = bt.Name
		for _, field := range bt.Fields {
			bloodTestFieldNames[field.Id] = field.Name
			bloodTestFieldUnits[field.Id] = field.Unit
		}
	}

	bloodTestResults := make([]BloodTestResult, 0, len(patient.BloodTests))
	for _, bt := range patient.BloodTests {
		fields := make([]BloodTestFilledField, 0, len(bt.FilledFields))
		for _, field := range bt.FilledFields {
			fields = append(fields, BloodTestFilledField{
				BloodTestFieldId: field.Id,
				Name:             bloodTestFieldNames[field.BloodTestFieldId],
				Unit:             bloodTestFieldUnits[field.BloodTestFieldId],
				MinValue:         0,
				MaxValue:         0,
				ValueNumber:      field.ValueNumber,
				ValueString:      field.ValueString,
			})
		}

		bloodTestResults = append(bloodTestResults, BloodTestResult{
			BloodTestId:  bt.BloodTestId,
			Name:         bloodTestNames[bt.BloodTestId],
			FilledFields: fields,
		})
	}

	outPatient := Patient{
		Id:          patient.Id,
		PublicId:    patient.PublicId,
		NationalId:  patient.NationalId,
		Nationality: patient.Nationality,
		FirstName:   patient.FirstName,
		LastName:    patient.LastName,
		FatherName:  patient.FatherName,
		MotherName:  patient.MotherName,
		PlaceOfBirth: Address{
			Id:          patient.PlaceOfBirth.Id,
			Governorate: patient.PlaceOfBirth.Governorate,
			Suburb:      patient.PlaceOfBirth.Suburb,
			Street:      patient.PlaceOfBirth.Street,
		},
		DateOfBirth: patient.DateOfBirth,
		Residency: Address{
			Id:          patient.Residency.Id,
			Governorate: patient.Residency.Governorate,
			Suburb:      patient.Residency.Suburb,
			Street:      patient.Residency.Street,
		},
		Gender:      patient.Gender,
		PhoneNumber: patient.PhoneNumber,
		BATScore:    patient.BATScore,
		Viri:        viruses,
		BloodTests:  bloodTestResults,
	}

	return GetPatientPayload{
		Data: outPatient,
	}, nil
}
