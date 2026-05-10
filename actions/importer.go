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
	FactorVIII            string
	BloodGroupABO         string
	BloodGroupRhD         string
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
			FactorVIII:            strings.TrimSpace(column[17]),
			BloodGroupABO:         strings.TrimSpace(column[18]),
			BloodGroupRhD:         strings.TrimSpace(column[19]),
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
			Gender:              record.Gender == "male",
			PhoneNumber:         record.PhoneNumber,
			FamilyHistoryExists: false,
			FirstVisitReason:    "",
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
			FactorViii: record.FactorVIII,
			CreatedAt:  record.DateOfDiagnosis,
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

	for _, patient := range inPatients {
		newPatient, err := a.app.CreatePatient(patient)
		if err != nil {
			log.Errorln("Failed to create patient: ", err)
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

	outIgnoredPatients := make([]Patient, len(ignoredPatients))
	for i := range ignoredPatients {
		outIgnoredPatients[i].FromModel(ignoredPatients[i])
	}

	return ImportPatientsFromCsvPayload{
		ImportCount:     len(newPatients),
		IgnoredPatients: outIgnoredPatients,
	}, nil
}
