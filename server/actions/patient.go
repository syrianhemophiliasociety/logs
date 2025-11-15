package actions

import (
	"shs/app/models"
	"shs/nanoid"
	"slices"
	"time"
)

type BloodTestFilledField struct {
	Name        string               `json:"name"`
	Unit        models.BlootTestUnit `json:"unit"`
	MinValue    uint                 `json:"min_value"`
	MaxValue    uint                 `json:"max_value"`
	ValueNumber uint                 `json:"value_number"`
	ValueString string               `json:"value_string"`
}

type BloodTestResult struct {
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
	BloodTests   []BloodTestResult `json:"blood_tests"`
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
	}

	placesOfBirth, _ := a.app.GetAllAddressesALike(models.Address{
		Governorate: params.NewPatient.PlaceOfBirth.Governorate,
		Suburb:      params.NewPatient.PlaceOfBirth.Suburb,
		Street:      params.NewPatient.PlaceOfBirth.Street,
	})

	if len(placesOfBirth) == 1 {
		newPatient.Residency.Id = placesOfBirth[0].Id
		newPatient.ResidencyId = placesOfBirth[0].Id
	}

	allViri, err := a.app.ListAllViri()
	if err != nil {
		return CreatePatientPayload{}, err
	}

	for _, virus := range params.NewPatient.Viri {
		matchedVirusIndex := slices.IndexFunc(allViri, func(v models.Virus) bool {
			return v.Name == virus.Name
		})
		if matchedVirusIndex < 0 {
			continue
		}
		newPatient.Viri = append(newPatient.Viri, allViri[matchedVirusIndex])
	}

	// TODO: store blood tests :)

	_, err = a.app.CreatePatient(newPatient)
	if err != nil {
		return CreatePatientPayload{}, err
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
