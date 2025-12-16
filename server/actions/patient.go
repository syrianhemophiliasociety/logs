package actions

import (
	"encoding/base64"
	"shs/app"
	"shs/app/models"
	"shs/cardgen"
	"slices"
	"strings"
	"time"
)

type DiagnosisResult struct {
	Diagnosis
	Id          uint      `json:"id"`
	DiagnosisId uint      `json:"diagnosis_id"`
	DiagnosedAt time.Time `json:"diagnosed_at"`

	CreatedAt time.Time `json:"created_at"`
}

type BloodTestFilledField struct {
	BloodTestFieldId uint                 `json:"blood_test_field_id"`
	Name             string               `json:"name"`
	Unit             models.BlootTestUnit `json:"unit"`
	ValueNumber      float64              `json:"value_number"`
	ValueString      string               `json:"value_string"`
}

type BloodTestResult struct {
	Id           uint                   `json:"id"`
	Name         string                 `json:"name"`
	BloodTestId  uint                   `json:"blood_test_id"`
	FilledFields []BloodTestFilledField `json:"filled_fields"`
	Pending      bool                   `json:"pending"`
	CreatedAt    time.Time              `json:"created_at"`
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
	Id                  uint               `json:"id"`
	PublicId            string             `json:"public_id"`
	NationalId          string             `json:"national_id"`
	Nationality         string             `json:"nationality"`
	FirstName           string             `json:"first_name"`
	LastName            string             `json:"last_name"`
	FatherName          string             `json:"father_name"`
	MotherName          string             `json:"mother_name"`
	PlaceOfBirth        Address            `json:"place_of_birth"`
	DateOfBirth         time.Time          `json:"date_of_birth"`
	Residency           Address            `json:"residency"`
	Gender              bool               `json:"gender"`
	PhoneNumber         string             `json:"phone_number"`
	BATScore            uint               `json:"bat_score"`
	FamilyHistoryExists bool               `json:"family_history_exists"`
	FirstVisitReason    string             `json:"first_visit_reason"`
	Viruses             []Virus            `json:"viruses"`
	BloodTestResults    []BloodTestResult  `json:"blood_test_results"`
	JointsEvaluations   []JointsEvaluation `json:"joints_evaluations"`
	Diagnoses           []DiagnosisResult  `json:"diagnoses"`
}

func (p Patient) IntoModel() models.Patient {
	viruses := make([]models.Virus, 0, len(p.Viruses))
	for _, v := range p.Viruses {
		viruses = append(viruses, models.Virus{
			Id:   v.Id,
			Name: v.Name,
		})
	}

	bloodTestResults := make([]models.BloodTestResult, 0, len(p.BloodTestResults))
	for _, btr := range p.BloodTestResults {
		bloodTestResultFields := make([]models.BloodTestFilledField, 0, len(btr.FilledFields))
		for _, field := range btr.FilledFields {
			bloodTestResultFields = append(bloodTestResultFields, models.BloodTestFilledField{
				BloodTestResultId: btr.BloodTestId,
				BloodTestFieldId:  field.BloodTestFieldId,
				ValueNumber:       field.ValueNumber,
				ValueString:       field.ValueString,
			})
		}

		bloodTestResults = append(bloodTestResults, models.BloodTestResult{
			BloodTestId:  btr.BloodTestId,
			FilledFields: bloodTestResultFields,
			Pending:      btr.Pending,
			CreatedAt:    btr.CreatedAt,
		})
	}

	return models.Patient{
		Id:          p.Id,
		PublicId:    p.PublicId,
		NationalId:  p.NationalId,
		Nationality: p.Nationality,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		FatherName:  p.FatherName,
		MotherName:  p.MotherName,
		PlaceOfBirth: models.Address{
			Governorate: p.PlaceOfBirth.Governorate,
			Suburb:      p.PlaceOfBirth.Suburb,
			Street:      p.PlaceOfBirth.Street,
		},
		DateOfBirth: p.DateOfBirth,
		Residency: models.Address{
			Governorate: p.Residency.Governorate,
			Suburb:      p.Residency.Suburb,
			Street:      p.Residency.Street,
		},
		Gender:              p.Gender,
		PhoneNumber:         p.PhoneNumber,
		FamilyHistoryExists: p.FamilyHistoryExists,
		FirstVisitReason:    models.PatientFirstVisitReason(p.FirstVisitReason),
		BATScore:            p.BATScore,
		Viruses:             viruses,
		BloodTestResults:    bloodTestResults,
	}
}

func (p *Patient) FromModel(patient models.Patient) {
	(*p) = Patient{
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
		Gender:              patient.Gender,
		PhoneNumber:         patient.PhoneNumber,
		BATScore:            patient.BATScore,
		FamilyHistoryExists: patient.FamilyHistoryExists,
		FirstVisitReason:    string(patient.FirstVisitReason),
	}
}

func (p *Patient) WithBloodTestResults(patientBloodTestResults []models.BloodTestResult, bloodTests []models.BloodTest) {
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

	(*p).BloodTestResults = make([]BloodTestResult, 0, len(patientBloodTestResults))
	for _, btr := range patientBloodTestResults {
		fields := make([]BloodTestFilledField, 0, len(btr.FilledFields))
		for _, field := range btr.FilledFields {
			fields = append(fields, BloodTestFilledField{
				BloodTestFieldId: field.BloodTestFieldId,
				Name:             bloodTestFieldNames[field.BloodTestFieldId],
				Unit:             bloodTestFieldUnits[field.BloodTestFieldId],
				ValueNumber:      field.ValueNumber,
				ValueString:      field.ValueString,
			})
		}

		(*p).BloodTestResults = append((*p).BloodTestResults, BloodTestResult{
			Id:           btr.Id,
			BloodTestId:  btr.BloodTestId,
			Name:         bloodTestNames[btr.BloodTestId],
			FilledFields: fields,
			Pending:      btr.Pending,
			CreatedAt:    btr.CreatedAt,
		})
	}
}

func (p *Patient) WithJointsEvaluations(jointsEvaluations []models.JointsEvaluation) {
	for _, je := range jointsEvaluations {
		outJointsEvaluation := new(JointsEvaluation)
		outJointsEvaluation.FromModel(je)
		(*p).JointsEvaluations = append((*p).JointsEvaluations, *outJointsEvaluation)
	}
}

func (p *Patient) WithViruses(patientViruses []models.Virus, viruses []models.Virus) {
	(*p).Viruses = make([]Virus, 0, len(patientViruses))
	for _, v := range patientViruses {
		(*p).Viruses = append((*p).Viruses, Virus{
			Id:   v.Id,
			Name: v.Name,
		})
	}
}

func (p *Patient) WithDiagnoses(diagnosesResults []models.DiagnosisResult, diagnoses []models.Diagnosis) {
	diagnosisMapped := make(map[uint]Diagnosis)

	for _, d := range diagnoses {
		outD := new(Diagnosis)
		outD.FromModel(d)
		diagnosisMapped[d.Id] = *outD
	}

	(*p).Diagnoses = make([]DiagnosisResult, 0, len(diagnosesResults))
	for _, dr := range diagnosesResults {
		(*p).Diagnoses = append((*p).Diagnoses, DiagnosisResult{
			Diagnosis:   diagnosisMapped[dr.DiagnosisId],
			DiagnosisId: dr.DiagnosisId,
			Id:          dr.Id,
			DiagnosedAt: dr.DiagnosedAt,
			CreatedAt:   dr.CreatedAt,
		})
	}
}

type CreatePatientParams struct {
	ActionContext
	NewPatient Patient `json:"new_patient"`
}

type CreatePatientPayload struct {
	PatientPublicId string `json:"id"`
}

func (a *Actions) CreatePatient(params CreatePatientParams) (CreatePatientPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWritePatient) {
		return CreatePatientPayload{}, ErrPermissionDenied{}
	}

	newPatient := models.Patient{
		NationalId:          params.NewPatient.NationalId,
		Nationality:         params.NewPatient.Nationality,
		FirstName:           params.NewPatient.FirstName,
		LastName:            params.NewPatient.LastName,
		FatherName:          params.NewPatient.FatherName,
		MotherName:          params.NewPatient.MotherName,
		DateOfBirth:         params.NewPatient.DateOfBirth,
		Gender:              params.NewPatient.Gender,
		PhoneNumber:         params.NewPatient.PhoneNumber,
		BATScore:            params.NewPatient.BATScore,
		FirstVisitReason:    models.PatientFirstVisitReason(params.NewPatient.FirstVisitReason),
		Viruses:             []models.Virus{},
		BloodTestResults:    []models.BloodTestResult{},
		FamilyHistoryExists: params.NewPatient.FamilyHistoryExists,
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

	newPatient, err := a.app.CreatePatient(newPatient)
	if err != nil {
		return CreatePatientPayload{}, err
	}

	// INFO: in case of minors without a national id, the password will be the patient's phone number without the country code
	password := params.NewPatient.NationalId
	if password == "" {
		password = cleanPhoneNumberCountryCode(params.NewPatient.PhoneNumber)
	}

	_, err = a.app.CreateAccount(models.Account{
		DisplayName: newPatient.FirstName + " " + newPatient.LastName,
		Username:    newPatient.PublicId,
		Password:    password,
		Type:        models.AccountTypePatient,
		Permissions: patientPermissions,
	})
	if err != nil {
		return CreatePatientPayload{}, err
	}

	return CreatePatientPayload{
		PatientPublicId: newPatient.PublicId,
	}, nil
}

type CreatePatientBloodTestResultParams struct {
	ActionContext
	PatientPublicId string          `json:"patient_id"`
	BloodTest       BloodTestResult `json:"patient_blood_test"`
}

type CreatePatientBloodTestResultPayload struct {
}

func (a *Actions) CreatePatientBloodTestResult(params CreatePatientBloodTestResultParams) (CreatePatientBloodTestResultPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWritePatient) {
		return CreatePatientBloodTestResultPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetFullPatientByPublicId(params.PatientPublicId)
	if err != nil {
		return CreatePatientBloodTestResultPayload{}, err
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
		Pending:      params.BloodTest.Pending,
	})
	if err != nil {
		return CreatePatientBloodTestResultPayload{}, err
	}

	return CreatePatientBloodTestResultPayload{}, nil
}

type CreatePatientDiagnosisResultParams struct {
	ActionContext
	PatientPublicId string          `json:"patient_id"`
	Diagnosis       DiagnosisResult `json:"patient_diagnosis"`
}

type CreatePatientDiagnosisResultPayload struct {
}

func (a *Actions) CreatePatientDiagnosisResult(params CreatePatientDiagnosisResultParams) (CreatePatientDiagnosisResultPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWritePatient) {
		return CreatePatientDiagnosisResultPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetFullPatientByPublicId(params.PatientPublicId)
	if err != nil {
		return CreatePatientDiagnosisResultPayload{}, err
	}

	_, err = a.app.CreateDiagnosisResult(models.DiagnosisResult{
		DiagnosisId: params.Diagnosis.DiagnosisId,
		PatientId:   patient.Id,
		DiagnosedAt: params.Diagnosis.DiagnosedAt,
	})
	if err != nil {
		return CreatePatientDiagnosisResultPayload{}, err
	}

	return CreatePatientDiagnosisResultPayload{}, nil
}

type UpdatePatientPendingBloodTestResultParams struct {
	ActionContext
	BloodTestResultId uint
	PatientPublicId   string
	FilledFields      []BloodTestFilledField `json:"filled_fields"`
}

type UpdatePatientPendingBloodTestResultPayload struct {
}

func (a *Actions) UpdatePatientPendingBloodTestResult(params UpdatePatientPendingBloodTestResultParams) (UpdatePatientPendingBloodTestResultPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWritePatient) {
		return UpdatePatientPendingBloodTestResultPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetFullPatientByPublicId(params.PatientPublicId)
	if err != nil {
		return UpdatePatientPendingBloodTestResultPayload{}, err
	}

	if !slices.ContainsFunc(patient.BloodTestResults, func(btr models.BloodTestResult) bool {
		return btr.Id == params.BloodTestResultId
	}) {
		return UpdatePatientPendingBloodTestResultPayload{}, app.ErrNotFound{
			ResourceName: "blood_test_result",
		}
	}

	bloodTestResultFields := make([]models.BloodTestFilledField, 0, len(params.FilledFields))
	for _, field := range params.FilledFields {
		bloodTestResultFields = append(bloodTestResultFields, models.BloodTestFilledField{
			BloodTestFieldId: field.BloodTestFieldId,
			ValueNumber:      field.ValueNumber,
			ValueString:      field.ValueString,
		})
	}

	err = a.app.UpdatePatientPendingBloodTestResultFields(params.BloodTestResultId, bloodTestResultFields)
	if err != nil {
		return UpdatePatientPendingBloodTestResultPayload{}, err
	}

	return UpdatePatientPendingBloodTestResultPayload{}, nil
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

	if !params.Account.HasPermission(models.AccountPermissionReadPatient) {
		return FindPatientsPayload{}, ErrPermissionDenied{}
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
		outPatient := new(Patient)
		outPatient.FromModel(patient)
		outPatients = append(outPatients, *outPatient)
	}

	return FindPatientsPayload{
		Data: outPatients,
	}, nil
}

type ListLastPatientsParams struct {
	ActionContext
}

type ListLastPatientsPayload struct {
	Data []Patient `json:"data"`
}

func (a *Actions) ListLastPatients(params ListLastPatientsParams) (ListLastPatientsPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadPatient) {
		return ListLastPatientsPayload{}, ErrPermissionDenied{}
	}

	patients, err := a.app.ListLastPatients(50)
	if err != nil {
		return ListLastPatientsPayload{}, err
	}

	outPatients := make([]Patient, 0, len(patients))
	for _, patient := range patients {
		outPatient := new(Patient)
		outPatient.FromModel(patient)
		outPatients = append(outPatients, *outPatient)
	}

	return ListLastPatientsPayload{
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
	if !params.Account.HasPermission(models.AccountPermissionReadPatient) {
		return GetPatientPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetFullPatientByPublicId(params.PublicId)
	if err != nil {
		return GetPatientPayload{}, err
	}

	bloodTests, err := a.app.ListAllBloodTests()
	if err != nil {
		return GetPatientPayload{}, err
	}

	diagnoses, err := a.app.ListAllDiagnoses()
	if err != nil {
		return GetPatientPayload{}, err
	}

	diagnosesResults, err := a.app.ListPatientDiagnosisResults(patient.Id)
	if err != nil {
		return GetPatientPayload{}, err
	}

	outPatient := &Patient{}
	outPatient.FromModel(patient)
	outPatient.WithViruses(patient.Viruses, nil)
	outPatient.WithBloodTestResults(patient.BloodTestResults, bloodTests)
	outPatient.WithJointsEvaluations(patient.JointsEvaluations)
	outPatient.WithDiagnoses(diagnosesResults, diagnoses)

	return GetPatientPayload{
		Data: *outPatient,
	}, nil
}

type DeletePatientParams struct {
	ActionContext
	PublicId string
}

type DeletePatientPayload struct {
}

func (a *Actions) DeletePatient(params DeletePatientParams) (DeletePatientPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWritePatient) {
		return DeletePatientPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetMinimalPatientByPublicId(params.PublicId)
	if err != nil {
		return DeletePatientPayload{}, err
	}

	err = a.app.DeletePatient(patient.Id)
	if err != nil {
		return DeletePatientPayload{}, err
	}

	return DeletePatientPayload{}, nil
}

type GeneratePatientCardParams struct {
	ActionContext
	PatientId string
}

type GeneratePatientCardPayload struct {
	ImageBase64 string `json:"image_base_64"`
}

func (a *Actions) GeneratePatientCard(params GeneratePatientCardParams) (GeneratePatientCardPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadPatient) {
		return GeneratePatientCardPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetMinimalPatientByPublicId(params.PatientId)
	if err != nil {
		return GeneratePatientCardPayload{}, err
	}

	patientCard := cardgen.NewBuffer(nil)
	generator, err := cardgen.New(patientCard, patient)
	if err != nil {
		return GeneratePatientCardPayload{}, err
	}

	err = generator.Generate(false)
	if err != nil {
		return GeneratePatientCardPayload{}, err
	}
	err = generator.Finalize()
	if err != nil {
		return GeneratePatientCardPayload{}, err
	}

	b64Img := base64.StdEncoding.EncodeToString(patientCard.Bytes())

	return GeneratePatientCardPayload{
		ImageBase64: b64Img,
	}, nil
}
