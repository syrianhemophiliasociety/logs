package apis

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"shs/actions"
	"shs/handlers/web/context"
	"shs/log"
	"shs/web/i18n"
	"shs/web/views/components"
	"strconv"
	"strings"
	"time"
)

type PatientRequest struct {
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	FatherName              string `json:"father_name"`
	MotherName              string `json:"mother_name"`
	Nationality             string `json:"nationality"`
	NationalId              string `json:"national_id"`
	Gender                  string `json:"gender"`
	PhoneNumber             string `json:"phone_number"`
	PhoneNumberCountryCode  string `json:"phone_number_country_code"`
	DateOfBirth             string `json:"date_of_birth"`
	PlaceOfBirthGovernorate string `json:"place_of_birth_governorate"`
	PlaceOfBirthSuburb      string `json:"place_of_birth_suburb"`
	PlaceOfBirthStreet      string `json:"place_of_birth_street"`
	ResidencyGovernorate    string `json:"residency_governorate"`
	ResidencySuburb         string `json:"residency_suburb"`
	ResidencyStreet         string `json:"residency_street"`
	FamilyHistoryExists     string `json:"family_history_exists"`
	FirstVisitReason        string `json:"first_visit_reason"`
	WBDR                    string `json:"wbdr"`
}

func clusterFuckPatientToActionsOne(p PatientRequest) actions.Patient {
	dateOfBirth, _ := time.Parse("2006-01-02", p.DateOfBirth)
	return actions.Patient{
		NationalId:  p.NationalId,
		Nationality: p.Nationality,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		FatherName:  p.FatherName,
		MotherName:  p.MotherName,
		PlaceOfBirth: actions.Address{
			Governorate: p.PlaceOfBirthGovernorate,
			Suburb:      p.PlaceOfBirthSuburb,
			Street:      p.PlaceOfBirthStreet,
		},
		DateOfBirth: dateOfBirth,
		Residency: actions.Address{
			Governorate: p.ResidencyGovernorate,
			Suburb:      p.ResidencySuburb,
			Street:      p.ResidencyStreet,
		},
		Gender:                 p.Gender == "male",
		PhoneNumber:            p.PhoneNumber,
		PhoneNumberCountryCode: p.PhoneNumberCountryCode,
		BATScore:               0,
		FirstVisitReason:       p.FirstVisitReason,
		FamilyHistoryExists:    p.FamilyHistoryExists == "on",
		WBDR:                   p.WBDR,
	}
}

type PatientBloodTests struct {
	BloodTests []actions.BloodTestResult
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

	bloodTestsFields := make(map[uint][]actions.BloodTestFilledField)
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

		testResultInt, _ := strconv.ParseFloat(testResult, 64)

		bloodTestNames[uint(id)] = name
		bloodTestsFields[uint(id)] = append(bloodTestsFields[uint(id)], actions.BloodTestFilledField{
			BloodTestFieldId: uint(fieldId),
			Name:             fieldName,
			Unit:             "",
			ValueNumber:      testResultInt,
			ValueString:      testResult,
		})
	}

	doTestLater, _ := data["do_later"].(string)

	for id, fields := range bloodTestsFields {
		cleanedFields := make([]actions.BloodTestFilledField, 0, len(fields))
		for _, f := range fields {
			if f.ValueString == "" {
				continue
			}
			cleanedFields = append(cleanedFields, f)
		}

		(*p).BloodTests = append((*p).BloodTests, actions.BloodTestResult{
			BloodTestId:  id,
			Name:         bloodTestNames[id],
			FilledFields: cleanedFields,
			Pending:      doTestLater == "on",
		})
	}

	return nil
}

type PatientDiagnosisRequest struct {
	DiagnosisId string `json:"diagnosis_id"`
	DiagnosedAt string `json:"diagnosed_at"`
}

type CreateCheckUpRequest struct {
	VisitReason         string
	VisitExtraDetails   string
	PatientWeight       float64
	PatientHeight       float64
	PrescribedMedicines []actions.Medicine
}

func (v *CreateCheckUpRequest) UnmarshalJSON(payload []byte) error {
	var data map[string]any
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}

	const (
		visitReasonKey       = "visit_reason"
		visitExtraDetailsKey = "visit_extra_details"
		medicineIdsKey       = "medicine_id"
		medicineAmountKey    = "amount"
		patientWeightKey     = "patient_weight"
		patientHeightKey     = "patient_height"
	)

	var ok bool
	(*v).VisitReason, ok = data[visitReasonKey].(string)
	if !ok {
		return errors.New("missing visit_reason")
	}
	(*v).VisitExtraDetails, _ = data[visitExtraDetailsKey].(string)
	weight, _ := data[patientWeightKey].(string)
	height, _ := data[patientHeightKey].(string)

	(*v).PatientWeight, _ = strconv.ParseFloat(weight, 64)
	(*v).PatientHeight, _ = strconv.ParseFloat(height, 64)

	_, ok = data[medicineIdsKey]
	if !ok {
		return nil
	}
	_, ok = data[medicineAmountKey]
	if !ok {
		return nil
	}

	switch data[medicineIdsKey].(type) {
	case string:
		mIdInt, err := strconv.Atoi(data[medicineIdsKey].(string))
		if err != nil {
			return err
		}
		(*v).PrescribedMedicines = []actions.Medicine{
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
			(*v).PrescribedMedicines = append((*v).PrescribedMedicines, actions.Medicine{
				Id: uint(mIdInt),
			})
		}
	}

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

////

type patientApi struct {
	usecases *actions.Actions
}

func NewPatientApi(usecases *actions.Actions) *patientApi {
	return &patientApi{
		usecases: usecases,
	}
}

func (v *patientApi) HandleCreatePatient(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody PatientRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	payload, err := v.usecases.CreatePatient(actions.CreatePatientParams{
		ActionContext: ctx,
		NewPatient:    clusterFuckPatientToActionsOne(reqBody),
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	w.Header().Set("HX-Redirect", "/patient/"+payload.PatientPublicId)
}

func (v *patientApi) HandleUpdatePatient(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	id := r.PathValue("id")

	var reqBody PatientRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.UpdatePatient(actions.UpdatePatientParams{
		ActionContext:   ctx,
		NewPatient:      clusterFuckPatientToActionsOne(reqBody),
		PatientPublicId: id,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}
}

func (v *patientApi) HandleCreatePatientBloodTestResult(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")

	var reqBody PatientBloodTests
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.CreatePatientBloodTestResult(actions.CreatePatientBloodTestResultParams{
		ActionContext:   ctx,
		PatientPublicId: patientId,
		BloodTest:       reqBody.BloodTests[0],
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *patientApi) HandleCreatePatientDiagnosisResult(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")

	var reqBody PatientDiagnosisRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	diagnosisId, err := strconv.Atoi(reqBody.DiagnosisId)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	diagnosedAt, err := time.Parse("2006-01-02", reqBody.DiagnosedAt)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.CreatePatientDiagnosisResult(actions.CreatePatientDiagnosisResultParams{
		ActionContext:   ctx,
		PatientPublicId: patientId,
		Diagnosis: actions.DiagnosisResult{
			DiagnosisId: uint(diagnosisId),
			DiagnosedAt: diagnosedAt,
		},
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *patientApi) HandleCreatePatientCheckUp(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")

	var reqBody CreateCheckUpRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		writeRawTextResponse(w, i18n.Strings("en").ErrorSomethingWentWrong)
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.CreatePatientVisit(actions.CreatePatientVisitParams{
		ActionContext:       ctx,
		PatientId:           patientId,
		VisitReason:         reqBody.VisitReason,
		VisitExtraDetails:   reqBody.VisitExtraDetails,
		PrescribedMedicines: reqBody.PrescribedMedicines,
	})
	if errors.Is(err, actions.ErrInsufficientMedicine{}) {
		imErr := err.(actions.ErrInsufficientMedicine)
		writeRawTextResponse(w, i18n.Strings("en").ErrorInsufficientMedicineAmountFmt(imErr.MedicineName, imErr.ExceedingAmount, imErr.LeftPackages))
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorInsufficientMedicineAmountFmt(imErr.MedicineName, imErr.ExceedingAmount, imErr.LeftPackages)).Render(r.Context(), w)
		log.Errorln(err)
		return
	}
	if err != nil {
		writeRawTextResponse(w, i18n.Strings("en").ErrorSomethingWentWrong)
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *patientApi) HandleGenerateCard(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")

	payload, err := v.usecases.GeneratePatientCard(actions.GeneratePatientCardParams{
		ActionContext: ctx,
		PatientId:     patientId,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	w.Write([]byte(payload.ImageBase64))
}

func (v *patientApi) HandleDeletePatient(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")

	_, err = v.usecases.DeletePatient(actions.DeletePatientParams{
		ActionContext: ctx,
		PublicId:      patientId,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *patientApi) HandleUpdatePatientPendingBloodTestResult(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")
	btrIdStr := r.PathValue("btr_id")
	btrId, err := strconv.Atoi(btrIdStr)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody PatientBloodTests
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.UpdatePatientPendingBloodTestResult(actions.UpdatePatientPendingBloodTestResultParams{
		ActionContext:     ctx,
		PatientPublicId:   patientId,
		BloodTestResultId: uint(btrId),
		FilledFields:      reqBody.BloodTests[0].FilledFields,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *patientApi) HandleCreatePatientJointsEvaluation(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	patientId := r.PathValue("id")

	var reqBody actions.JointsEvaluation
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.CreatePatientJointsEvaluation(actions.CreatePatientJointsEvaluationParams{
		ActionContext:    ctx,
		PatientId:        patientId,
		JointsEvaluation: reqBody,
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func validateFileType(r io.ReadSeeker, wantedTypes ...string) error {
	reader := bufio.NewReader(r)

	bytes, err := reader.Peek(256)
	if err != nil && err != io.EOF {
		return err
	}
	r.Seek(0, 0)

	fileType := http.DetectContentType(bytes)
	for _, wantedType := range wantedTypes {
		if strings.Contains(fileType, wantedType) {
			return nil
		}
	}

	return ErrInvalidFileType{
		Want: strings.Join(wantedTypes, ","),
		Got:  fileType,
	}
}

func (v *patientApi) HandleUploadImportPatientsFromCsv(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}
	r.ParseMultipartForm(32 << 20) // 32 MB

	file, _, err := r.FormFile("patient_records")
	if err != nil {
		log.Warningf("upload error: %v", err)
		w.Write([]byte("upload failed"))
		return
	}
	defer file.Close()

	if err := validateFileType(file, "text/plain", "application/vnd.ms-excel"); err != nil {
		log.Errorln(err.(ErrInvalidFileType).Got)
		w.Write([]byte("invalid file type"))
		return
	}

	payload, err := v.usecases.ImportPatientsFromCsv(actions.ImportPatientsFromCsvParams{
		ActionContext: ctx,
		CsvFile:       file,
	})
	if len(payload.IgnoredPatients) > 0 {
		w.Write([]byte("Ignored Patients:<br/><ul>"))
		for _, patient := range payload.IgnoredPatients {
			fmt.Fprintf(w, "<li>%s %s Son of %s and %s</li>", patient.FirstName, patient.LastName, patient.FatherName, patient.MotherName)
		}
		w.Write([]byte("</ul>"))
	}
}

func (v *patientApi) HandlePatientUseMedicine(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	visitIdStr := r.PathValue("visit_id")
	medIdStr := r.PathValue("med_id")
	visitId, err := strconv.Atoi(visitIdStr)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}
	medId, err := strconv.Atoi(medIdStr)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	_, err = v.usecases.UseMedicineForVisit(actions.UseMedicineForVisitParams{
		ActionContext:        ctx,
		VisitId:              uint(visitId),
		PrescribedMedicineId: uint(medId),
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}
