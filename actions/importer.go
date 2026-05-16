package actions

import (
	"encoding/csv"
	"fmt"
	"io"
	"shs/app/models"
	"shs/log"
	"slices"
	"strconv"
	"strings"
	"time"
)

type csvRow struct {
	FirstName             string
	LastName              string
	FatherName            string
	MotherName            string
	Nationality           string
	NationalID            string
	Gender                string
	DateOfBirth           time.Time
	PhoneNumber           string
	POB_Governorate       string
	POB_Suburb            string
	POB_Street            string
	Residency_Governorate string
	Residency_Suburb      string
	Residency_Street      string
	Diagnosis_GroupName   string
	Diagnosis_Title       string
	DateOfDiagnosis       time.Time
	BTFactorVIII          string
	BloodGroupABO         string
	BloodGroupRhD         string
	BTFactorIX            string
	BTVWFAg               string
	BTFactorV             string
	BTFactorX             string
	BTFibrinogen          string
	BTFactorVII           string
	BTInhibitorsScreening string
	BTInhibitorsTitrage   string
}

func tryParseTime(dateStr string) (time.Time, error) {
	layouts := []string{"2/1/2006", "02/01/2006", "2/01/2006", "02/1/2006"}
	for _, l := range layouts {
		if t, err := time.Parse(l, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}.Add(69 * time.Minute), fmt.Errorf("could not parse date: %s", dateStr)
}

func extractCsvRecords(csvFile io.Reader) ([]csvRow, error) {
	reader := csv.NewReader(csvFile)

	_, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var rows []csvRow

	for i := 0; ; i++ {
		if i == 0 {
			continue
		}
		column, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Warningf("Error reading row: %v", err)
			continue
		}

		dateOfBirth, _ := tryParseTime(column[7])
		dateOfDiagnosis, _ := tryParseTime(column[16])
		diagnosisSplit := strings.Split(strings.TrimSpace(column[15]), "#")
		diagnosisGroup := ""
		if len(diagnosisSplit) > 0 {
			diagnosisGroup = diagnosisSplit[0]
		}
		diagnosisTitle := ""
		if len(diagnosisSplit) > 1 {
			diagnosisTitle = diagnosisSplit[1]
		}

		r := csvRow{
			FirstName:             strings.TrimSpace(column[0]),
			LastName:              strings.TrimSpace(column[1]),
			FatherName:            strings.TrimSpace(column[2]),
			MotherName:            strings.TrimSpace(column[3]),
			Nationality:           strings.ToLower(strings.TrimSpace(column[4])),
			NationalID:            strings.TrimSpace(column[5]),
			Gender:                strings.ToLower(strings.TrimSpace(column[6])),
			DateOfBirth:           dateOfBirth,
			PhoneNumber:           strings.TrimSpace(column[8]),
			POB_Governorate:       strings.TrimSpace(column[9]),
			POB_Suburb:            strings.TrimSpace(column[10]),
			POB_Street:            strings.TrimSpace(column[11]),
			Residency_Governorate: strings.TrimSpace(column[12]),
			Residency_Suburb:      strings.TrimSpace(column[13]),
			Residency_Street:      strings.TrimSpace(column[14]),
			Diagnosis_GroupName:   diagnosisGroup,
			Diagnosis_Title:       diagnosisTitle,
			DateOfDiagnosis:       dateOfDiagnosis,
			BTFactorVIII:          strings.TrimSpace(column[17]),
			BloodGroupABO:         strings.TrimSpace(column[18]),
			BloodGroupRhD:         strings.TrimSpace(column[19]),
			BTFactorIX:            strings.TrimSpace(column[20]),
			BTVWFAg:               strings.TrimSpace(column[21]),
			BTFactorV:             strings.TrimSpace(column[22]),
			BTFactorX:             strings.TrimSpace(column[23]),
			BTFibrinogen:          strings.TrimSpace(column[24]),
			BTFactorVII:           strings.TrimSpace(column[25]),
			BTInhibitorsScreening: strings.TrimSpace(column[26]),
			BTInhibitorsTitrage:   strings.TrimSpace(column[27]),
		}

		rows = append(rows, r)
	}

	return rows, nil
}

type ImportPatientsFromCsvParams struct {
	ActionContext
	CsvFile io.Reader
}

type ImportPatientsFromCsvPayload struct {
	ImportCount     int       `json:"import_count"`
	IgnoredPatients []Patient `json:"ignored_patients"`
}

type patientBloodGroup struct {
	Id         uint
	ABOFieldId uint
	ABO        string
	RhFieldId  uint
	Rh         string
	CreatedAt  time.Time
}

type patientFactorVII struct {
	Id        uint
	FieldId   uint
	FactorVii string
	CreatedAt time.Time
}

type patientFactorV struct {
	Id        uint
	FieldId   uint
	FactorV   string
	CreatedAt time.Time
}

type patientFactorX struct {
	Id        uint
	FieldId   uint
	FactorX   string
	CreatedAt time.Time
}

type patientFactorIX struct {
	Id        uint
	FieldId   uint
	FactorIX  string
	CreatedAt time.Time
}

type patientVWFAg struct {
	Id        uint
	FieldId   uint
	VWFAg     string
	CreatedAt time.Time
}

type patientFibrinogen struct {
	Id        uint
	FieldId   uint
	Fibi      string
	CreatedAt time.Time
}

type patientInhibitors struct {
	Id        uint
	FieldId   uint
	Screening string
	Field2Id  uint
	Titrage   string
	CreatedAt time.Time
}

type patientFactorVIII struct {
	Id         uint
	FieldId    uint
	FactorViii string
	CreatedAt  time.Time
}

func (a *Actions) ImportPatientsFromCsv(params ImportPatientsFromCsvParams) (ImportPatientsFromCsvPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWritePatient) {
		return ImportPatientsFromCsvPayload{}, ErrPermissionDenied{}
	}
	if !params.Account.HasPermission(models.AccountPermissionWriteBloodTest) {
		return ImportPatientsFromCsvPayload{}, ErrPermissionDenied{}
	}
	if !params.Account.HasPermission(models.AccountPermissionWriteDiagnoses) {
		return ImportPatientsFromCsvPayload{}, ErrPermissionDenied{}
	}

	importRecords, err := extractCsvRecords(params.CsvFile)
	if err != nil {
		return ImportPatientsFromCsvPayload{}, err
	}

	patients := make([]models.Patient, 0, len(importRecords))
	patientDiagnoses := make(map[string]*models.Diagnosis)
	mPatientBloodGroup := make(map[string]*patientBloodGroup)
	mPatientFactorVIII := make(map[string]*patientFactorVIII)
	mPatientFactorVII := make(map[string]*patientFactorVII)
	mPatientFactorV := make(map[string]*patientFactorV)
	mPatientFactorIX := make(map[string]*patientFactorIX)
	mPatientFactorX := make(map[string]*patientFactorX)
	mPatientVWFAg := make(map[string]*patientVWFAg)
	mPatientFibrinogen := make(map[string]*patientFibrinogen)
	mPatientInhibitors := make(map[string]*patientInhibitors)

	for _, record := range importRecords {
		patient := models.Patient{
			NationalId:  record.NationalID,
			Nationality: record.Nationality,
			FirstName:   record.FirstName,
			LastName:    record.LastName,
			FatherName:  record.FatherName,
			MotherName:  record.MotherName,
			PlaceOfBirth: models.Address{
				Governorate: record.POB_Governorate,
				Suburb:      record.POB_Suburb,
				Street:      record.POB_Street,
			},
			DateOfBirth: record.DateOfBirth,
			Residency: models.Address{
				Governorate: record.Residency_Governorate,
				Suburb:      record.Residency_Suburb,
				Street:      record.Residency_Street,
			},
			Gender:                 record.Gender == "male",
			PhoneNumber:            record.PhoneNumber,
			PhoneNumberCountryCode: "+963",
			FamilyHistoryExists:    false,
			FirstVisitReason:       "",
		}
		patients = append(patients, patient)

		patientDiagnoses[patient.IndexId()] = &models.Diagnosis{
			GroupName: record.Diagnosis_GroupName,
			Title:     record.Diagnosis_Title,
			CreatedAt: record.DateOfDiagnosis,
		}

		mPatientBloodGroup[patient.IndexId()] = &patientBloodGroup{
			Id:        0,
			ABO:       record.BloodGroupABO,
			Rh:        record.BloodGroupRhD,
			CreatedAt: record.DateOfDiagnosis,
		}

		mPatientFactorVIII[patient.IndexId()] = &patientFactorVIII{
			Id:         0,
			FactorViii: record.BTFactorVIII,
			CreatedAt:  record.DateOfDiagnosis,
		}

		mPatientFactorVII[patient.IndexId()] = &patientFactorVII{
			Id:        0,
			FactorVii: record.BTFactorVII,
			CreatedAt: record.DateOfDiagnosis,
		}

		mPatientFactorV[patient.IndexId()] = &patientFactorV{
			Id:        0,
			FactorV:   record.BTFactorV,
			CreatedAt: record.DateOfDiagnosis,
		}

		mPatientFactorIX[patient.IndexId()] = &patientFactorIX{
			Id:        0,
			FactorIX:  record.BTFactorIX,
			CreatedAt: record.DateOfDiagnosis,
		}

		mPatientFactorX[patient.IndexId()] = &patientFactorX{
			Id:        0,
			FactorX:   record.BTFactorX,
			CreatedAt: record.DateOfDiagnosis,
		}

		mPatientVWFAg[patient.IndexId()] = &patientVWFAg{
			Id:        0,
			VWFAg:     record.BTVWFAg,
			CreatedAt: record.DateOfDiagnosis,
		}

		mPatientFibrinogen[patient.IndexId()] = &patientFibrinogen{
			Id:        0,
			Fibi:      record.BTFibrinogen,
			CreatedAt: record.DateOfDiagnosis,
		}

		mPatientInhibitors[patient.IndexId()] = &patientInhibitors{
			Id:        0,
			Screening: record.BTInhibitorsScreening,
			Titrage:   record.BTInhibitorsTitrage,
			CreatedAt: record.DateOfDiagnosis,
		}
	}

	inPatients := make([]models.Patient, 0, len(patients))
	ignoredPatients := make([]models.Patient, 0, len(patients))
	newPatients := make([]models.Patient, 0, len(patients))

	for i, patient := range patients {
		existingPatient, _ := a.app.FindPatientsByIndexFields(models.PatientIndexFields{
			FirstName:  patient.FirstName,
			LastName:   patient.LastName,
			FatherName: patient.FatherName,
			MotherName: patient.MotherName,
		})
		if len(existingPatient) > 0 {
			ignoredPatients = append(ignoredPatients, existingPatient...)
			continue
		}
		inPatients = append(inPatients, patients[i])
	}

	for i := range inPatients {
		inPatients[i].FillEmptyFieldsUsingPublicId()
		newPatient, err := a.app.CreatePatient(inPatients[i])
		if err != nil {
			log.Errorln("Failed to create patient: ", err)
			continue
		}

		// INFO: in case of minors without a national id, the password will be the patient's phone number without the country code
		password := inPatients[i].NationalId
		if password == "" {
			password = cleanPhoneNumberCountryCode(newPatient.PhoneNumber)
		}

		_, err = a.app.CreateAccount(models.Account{
			DisplayName: newPatient.FirstName + " " + newPatient.LastName,
			Username:    newPatient.PublicId,
			Password:    password,
			Type:        models.AccountTypePatient,
			Permissions: patientPermissions,
		})
		if err != nil {
			log.Errorln("Failed to create patient's account: ", err)
			continue
		}

		newPatients = append(newPatients, newPatient)
	}

	diagnoses, err := a.app.ListAllDiagnoses()
	if err != nil {
		log.Warningln("No diagnoses were found,", err)
	}
	for key, diagnosis := range patientDiagnoses {
		foundDiagnosisIdx := slices.IndexFunc(diagnoses, func(d models.Diagnosis) bool {
			return diagnosis.GroupName == d.GroupName &&
				diagnosis.Title == d.Title
		})
		if foundDiagnosisIdx > -1 {
			patientDiagnoses[key].Id = diagnoses[foundDiagnosisIdx].Id
		}
	}

	bloodTests, err := a.app.ListAllBloodTests()
	if err != nil {
		log.Warningln("No blood tests were found,", err)
	}
	for key := range mPatientBloodGroup {
		foundBtIdx := slices.IndexFunc(bloodTests, func(bt models.BloodTest) bool {
			return bt.Name == "Blood Group"
		})
		if foundBtIdx > -1 {
			mPatientBloodGroup[key].Id = bloodTests[foundBtIdx].Id
		}
		foundBtFieldAboIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "ABO"
		})
		if foundBtFieldAboIdx > -1 {
			mPatientBloodGroup[key].ABOFieldId = bloodTests[foundBtIdx].Fields[foundBtFieldAboIdx].Id
		}
		foundBtFieldRhIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "Rh(D)"
		})
		if foundBtFieldRhIdx > -1 {
			mPatientBloodGroup[key].RhFieldId = bloodTests[foundBtIdx].Fields[foundBtFieldRhIdx].Id
		}
	}

	for key := range mPatientFactorVIII {
		foundBtIdx := slices.IndexFunc(bloodTests, func(bt models.BloodTest) bool {
			return bt.Name == "Factor - VIII"
		})
		if foundBtIdx > -1 {
			mPatientFactorVIII[key].Id = bloodTests[foundBtIdx].Id
		}
		foundBtFieldAboIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "Factor - VIII"
		})
		if foundBtFieldAboIdx > -1 {
			mPatientFactorVIII[key].FieldId = bloodTests[foundBtIdx].Fields[foundBtFieldAboIdx].Id
		}
	}

	for key := range mPatientFactorVII {
		foundBtIdx := slices.IndexFunc(bloodTests, func(bt models.BloodTest) bool {
			return bt.Name == "Factor - VII"
		})
		if foundBtIdx > -1 {
			mPatientFactorVII[key].Id = bloodTests[foundBtIdx].Id
		}
		foundBtFieldAboIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "Factor - VII"
		})
		if foundBtFieldAboIdx > -1 {
			mPatientFactorVII[key].FieldId = bloodTests[foundBtIdx].Fields[foundBtFieldAboIdx].Id
		}
	}

	for key := range mPatientFactorV {
		foundBtIdx := slices.IndexFunc(bloodTests, func(bt models.BloodTest) bool {
			return bt.Name == "Factor - V"
		})
		if foundBtIdx > -1 {
			mPatientFactorV[key].Id = bloodTests[foundBtIdx].Id
		}
		foundBtFieldAboIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "Factor - V"
		})
		if foundBtFieldAboIdx > -1 {
			mPatientFactorV[key].FieldId = bloodTests[foundBtIdx].Fields[foundBtFieldAboIdx].Id
		}
	}

	for key := range mPatientFactorX {
		foundBtIdx := slices.IndexFunc(bloodTests, func(bt models.BloodTest) bool {
			return bt.Name == "Factor - X"
		})
		if foundBtIdx > -1 {
			mPatientFactorX[key].Id = bloodTests[foundBtIdx].Id
		}
		foundBtFieldAboIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "Factor - X"
		})
		if foundBtFieldAboIdx > -1 {
			mPatientFactorX[key].FieldId = bloodTests[foundBtIdx].Fields[foundBtFieldAboIdx].Id
		}
	}

	for key := range mPatientFactorIX {
		foundBtIdx := slices.IndexFunc(bloodTests, func(bt models.BloodTest) bool {
			return bt.Name == "Factor - IX"
		})
		if foundBtIdx > -1 {
			mPatientFactorIX[key].Id = bloodTests[foundBtIdx].Id
		}
		foundBtFieldAboIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "Factor - IX"
		})
		if foundBtFieldAboIdx > -1 {
			mPatientFactorIX[key].FieldId = bloodTests[foundBtIdx].Fields[foundBtFieldAboIdx].Id
		}
	}

	for key := range mPatientVWFAg {
		foundBtIdx := slices.IndexFunc(bloodTests, func(bt models.BloodTest) bool {
			return bt.Name == "VWF:Ag"
		})
		if foundBtIdx > -1 {
			mPatientVWFAg[key].Id = bloodTests[foundBtIdx].Id
		}
		foundBtFieldAboIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "VWF:Ag"
		})
		if foundBtFieldAboIdx > -1 {
			mPatientVWFAg[key].FieldId = bloodTests[foundBtIdx].Fields[foundBtFieldAboIdx].Id
		}
	}

	for key := range mPatientFibrinogen {
		foundBtIdx := slices.IndexFunc(bloodTests, func(bt models.BloodTest) bool {
			return bt.Name == "Fibrinogen"
		})
		if foundBtIdx > -1 {
			mPatientFibrinogen[key].Id = bloodTests[foundBtIdx].Id
		}
		foundBtFieldAboIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "Fibrinogen"
		})
		if foundBtFieldAboIdx > -1 {
			mPatientFibrinogen[key].FieldId = bloodTests[foundBtIdx].Fields[foundBtFieldAboIdx].Id
		}
	}

	for key := range mPatientInhibitors {
		foundBtIdx := slices.IndexFunc(bloodTests, func(bt models.BloodTest) bool {
			return bt.Name == "Inhibitors"
		})
		if foundBtIdx > -1 {
			mPatientInhibitors[key].Id = bloodTests[foundBtIdx].Id
		}
		foundBtFieldAboIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "Inhibitor Screening"
		})
		if foundBtFieldAboIdx > -1 {
			mPatientInhibitors[key].FieldId = bloodTests[foundBtIdx].Fields[foundBtFieldAboIdx].Id
		}
		foundBtField2AboIdx := slices.IndexFunc(bloodTests[foundBtIdx].Fields, func(btf models.BloodTestField) bool {
			return btf.Name == "Inhibitor Titrage"
		})
		if foundBtField2AboIdx > -1 {
			mPatientInhibitors[key].Field2Id = bloodTests[foundBtIdx].Fields[foundBtField2AboIdx].Id
		}
	}

	for _, patient := range newPatients {
		patientDiagnosis, ok := patientDiagnoses[patient.IndexId()]
		if !ok {
			log.Warningf("Diagnosis was not found for patient '%s'\n", patient.IndexId())
			continue
		}

		_, err := a.app.CreateDiagnosisResult(models.DiagnosisResult{
			DiagnosisId: patientDiagnosis.Id,
			PatientId:   patient.Id,
			CreatedAt:   patientDiagnosis.CreatedAt,
		})
		if err != nil {
			log.Warningf("Failed to assign '%s - %s' diagnosis to patient with id %s\n", patientDiagnosis.GroupName, patientDiagnosis.Title, patient.PublicId)
			continue
		}
	}

	for _, patient := range newPatients {
		// blood groups
		patientBloodGroup, ok := mPatientBloodGroup[patient.IndexId()]
		if !ok {
			log.Warningf("Blood group was not found for patient '%s'\n", patient.IndexId())
			continue
		}

		_, err = a.app.CreateBloodTestResult(models.BloodTestResult{
			CreatedAt:   patientBloodGroup.CreatedAt,
			BloodTestId: patientBloodGroup.Id,
			PatientId:   patient.Id,
			FilledFields: []models.BloodTestFilledField{
				{
					CreatedAt:        patientBloodGroup.CreatedAt,
					BloodTestFieldId: patientBloodGroup.RhFieldId,
					ValueString:      patientBloodGroup.Rh,
				},
				{
					CreatedAt:        patientBloodGroup.CreatedAt,
					BloodTestFieldId: patientBloodGroup.ABOFieldId,
					ValueString:      patientBloodGroup.ABO,
				},
			},
		})
		if err != nil {
			continue
		}
	}

	for _, patient := range newPatients {
		// factor viii
		patientFactor7, ok := mPatientFactorVIII[patient.IndexId()]
		if !ok {
			log.Warningf("Factor 7 was not found for patient '%s'\n", patient.IndexId())
			continue
		}

		patientFactor7Value, _ := strconv.ParseFloat(patientFactor7.FactorViii, 64)

		_, err = a.app.CreateBloodTestResult(models.BloodTestResult{
			CreatedAt:   patientFactor7.CreatedAt,
			BloodTestId: patientFactor7.Id,
			PatientId:   patient.Id,
			FilledFields: []models.BloodTestFilledField{
				{
					CreatedAt:        patientFactor7.CreatedAt,
					BloodTestFieldId: patientFactor7.FieldId,
					ValueString:      patientFactor7.FactorViii,
					ValueNumber:      patientFactor7Value,
				},
			},
		})
		if err != nil {
			continue
		}
	}

	for _, patient := range newPatients {
		// factor vii
		patientFactor6, ok := mPatientFactorVII[patient.IndexId()]
		if !ok {
			log.Warningf("Factor 6 was not found for patient '%s'\n", patient.IndexId())
			continue
		}

		patientFactor6Value, _ := strconv.ParseFloat(patientFactor6.FactorVii, 64)

		_, err = a.app.CreateBloodTestResult(models.BloodTestResult{
			CreatedAt:   patientFactor6.CreatedAt,
			BloodTestId: patientFactor6.Id,
			PatientId:   patient.Id,
			FilledFields: []models.BloodTestFilledField{
				{
					CreatedAt:        patientFactor6.CreatedAt,
					BloodTestFieldId: patientFactor6.FieldId,
					ValueString:      patientFactor6.FactorVii,
					ValueNumber:      patientFactor6Value,
				},
			},
		})
		if err != nil {
			continue
		}
	}

	for _, patient := range newPatients {
		// factor v
		patientFactor5, ok := mPatientFactorV[patient.IndexId()]
		if !ok {
			log.Warningf("Factor 5 was not found for patient '%s'\n", patient.IndexId())
			continue
		}

		patientFactor5Value, _ := strconv.ParseFloat(patientFactor5.FactorV, 64)

		_, err = a.app.CreateBloodTestResult(models.BloodTestResult{
			CreatedAt:   patientFactor5.CreatedAt,
			BloodTestId: patientFactor5.Id,
			PatientId:   patient.Id,
			FilledFields: []models.BloodTestFilledField{
				{
					CreatedAt:        patientFactor5.CreatedAt,
					BloodTestFieldId: patientFactor5.FieldId,
					ValueString:      patientFactor5.FactorV,
					ValueNumber:      patientFactor5Value,
				},
			},
		})
		if err != nil {
			continue
		}
	}

	for _, patient := range newPatients {
		// factor ix
		patientFactor9, ok := mPatientFactorIX[patient.IndexId()]
		if !ok {
			log.Warningf("Factor 9 was not found for patient '%s'\n", patient.IndexId())
			continue
		}

		patientFactor9Value, _ := strconv.ParseFloat(patientFactor9.FactorIX, 64)

		_, err = a.app.CreateBloodTestResult(models.BloodTestResult{
			CreatedAt:   patientFactor9.CreatedAt,
			BloodTestId: patientFactor9.Id,
			PatientId:   patient.Id,
			FilledFields: []models.BloodTestFilledField{
				{
					CreatedAt:        patientFactor9.CreatedAt,
					BloodTestFieldId: patientFactor9.FieldId,
					ValueString:      patientFactor9.FactorIX,
					ValueNumber:      patientFactor9Value,
				},
			},
		})
		if err != nil {
			continue
		}
	}

	for _, patient := range newPatients {
		// factor x
		patientFactor10, ok := mPatientFactorX[patient.IndexId()]
		if !ok {
			log.Warningf("Factor 10 was not found for patient '%s'\n", patient.IndexId())
			continue
		}

		patientFactor10Value, _ := strconv.ParseFloat(patientFactor10.FactorX, 64)

		_, err = a.app.CreateBloodTestResult(models.BloodTestResult{
			CreatedAt:   patientFactor10.CreatedAt,
			BloodTestId: patientFactor10.Id,
			PatientId:   patient.Id,
			FilledFields: []models.BloodTestFilledField{
				{
					CreatedAt:        patientFactor10.CreatedAt,
					BloodTestFieldId: patientFactor10.FieldId,
					ValueString:      patientFactor10.FactorX,
					ValueNumber:      patientFactor10Value,
				},
			},
		})
		if err != nil {
			continue
		}
	}

	for _, patient := range newPatients {
		// vwfag
		patientVWFAg, ok := mPatientVWFAg[patient.IndexId()]
		if !ok {
			log.Warningf("VWFAg not found for patient '%s'\n", patient.IndexId())
			continue
		}

		patientVWFAgValue, _ := strconv.ParseFloat(patientVWFAg.VWFAg, 64)

		_, err = a.app.CreateBloodTestResult(models.BloodTestResult{
			CreatedAt:   patientVWFAg.CreatedAt,
			BloodTestId: patientVWFAg.Id,
			PatientId:   patient.Id,
			FilledFields: []models.BloodTestFilledField{
				{
					CreatedAt:        patientVWFAg.CreatedAt,
					BloodTestFieldId: patientVWFAg.FieldId,
					ValueString:      patientVWFAg.VWFAg,
					ValueNumber:      patientVWFAgValue,
				},
			},
		})
		if err != nil {
			continue
		}
	}

	for _, patient := range newPatients {
		// inhibitors
		patientInhibitors, ok := mPatientInhibitors[patient.IndexId()]
		if !ok {
			log.Warningf("Inhibitors was not found for patient '%s'\n", patient.IndexId())
			continue
		}

		patientInTitValue, _ := strconv.ParseFloat(patientInhibitors.Titrage, 64)

		_, err = a.app.CreateBloodTestResult(models.BloodTestResult{
			CreatedAt:   patientInhibitors.CreatedAt,
			BloodTestId: patientInhibitors.Id,
			PatientId:   patient.Id,
			FilledFields: []models.BloodTestFilledField{
				{
					CreatedAt:        patientInhibitors.CreatedAt,
					BloodTestFieldId: patientInhibitors.FieldId,
					ValueString:      patientInhibitors.Screening,
				},
				{
					CreatedAt:        patientInhibitors.CreatedAt,
					BloodTestFieldId: patientInhibitors.Field2Id,
					ValueString:      patientInhibitors.Titrage,
					ValueNumber:      patientInTitValue,
				},
			},
		})
		if err != nil {
			continue
		}
	}

	outIgnoredPatients := make([]Patient, len(ignoredPatients))
	for i := range ignoredPatients {
		outIgnoredPatients[i].FromModel(ignoredPatients[i])
	}

	return ImportPatientsFromCsvPayload{
		ImportCount:     len(newPatients),
		IgnoredPatients: outIgnoredPatients,
	}, nil
}
