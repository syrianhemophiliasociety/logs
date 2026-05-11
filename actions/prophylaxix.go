package actions

import (
	"shs/app/models"
	"time"
)

const (
	prophylaxisFrequencyEvery4Weeks  float32 = 0.035 // 1/28
	prophylaxisFrequencyEvery2Weeks  float32 = 0.071 // 1/14
	prophylaxisFrequencyOnceInWeek   float32 = 0.142 // 1/7
	prophylaxisFrequencyTwiceInWeek  float32 = 0.285 // 2/7
	prophylaxisFrequencyThriceInWeek float32 = 0.428 // 3/7

	prophylaxisFrequencyNameEvery4Weeks  = "every4weeks"
	prophylaxisFrequencyNameEvery2Weeks  = "every2weeks"
	prophylaxisFrequencyNameOnceInWeek   = "once_in_week"
	prophylaxisFrequencyNameTwiceInWeek  = "twice_in_week"
	prophylaxisFrequencyNameThriceInWeek = "thrice_in_week"
)

var (
	prophylaxisFrequencyMapper = map[string]float32{
		prophylaxisFrequencyNameEvery4Weeks:  prophylaxisFrequencyEvery4Weeks,
		prophylaxisFrequencyNameEvery2Weeks:  prophylaxisFrequencyEvery2Weeks,
		prophylaxisFrequencyNameOnceInWeek:   prophylaxisFrequencyOnceInWeek,
		prophylaxisFrequencyNameTwiceInWeek:  prophylaxisFrequencyTwiceInWeek,
		prophylaxisFrequencyNameThriceInWeek: prophylaxisFrequencyThriceInWeek,
	}
	prophylaxisFrequencyMapperHuh = map[float32]string{
		prophylaxisFrequencyEvery4Weeks:  prophylaxisFrequencyNameEvery4Weeks,
		prophylaxisFrequencyEvery2Weeks:  prophylaxisFrequencyNameEvery2Weeks,
		prophylaxisFrequencyOnceInWeek:   prophylaxisFrequencyNameOnceInWeek,
		prophylaxisFrequencyTwiceInWeek:  prophylaxisFrequencyNameTwiceInWeek,
		prophylaxisFrequencyThriceInWeek: prophylaxisFrequencyNameThriceInWeek,
	}
)

func ProphylaxisFrequncyNumberToString(num float32) string {
	return prophylaxisFrequencyMapperHuh[num]
}

type Prophylaxis struct {
	Id                 uint      `json:"id"`
	Title              string    `json:"title"`
	FrequencyPerDays   string    `json:"frequency"`
	EndDate            time.Time `json:"end_date"`
	MedicineId         uint      `json:"medicine_id,omitempty"`
	PrescribedMedicine Medicine  `json:"prescribed_medicine"`
	MedicineDose       int       `json:"medicine_dose"`
}

func (pp *Prophylaxis) FromModel(p models.Prophylaxis) {
	(*pp).Id = p.Id
	(*pp).Title = p.Title
	(*pp).FrequencyPerDays = prophylaxisFrequencyMapperHuh[p.FrequencyPerDays]
	(*pp).EndDate = p.EndDate
	(*pp).MedicineDose = p.MedicineDose
	med := new(Medicine)
	med.FromModel(p.Medicine)
	(*pp).PrescribedMedicine = *med
}

func (pp Prophylaxis) IntoModel() models.Prophylaxis {
	return models.Prophylaxis{
		Id:               pp.Id,
		Title:            pp.Title,
		FrequencyPerDays: prophylaxisFrequencyMapper[pp.FrequencyPerDays],
		EndDate:          pp.EndDate,
		MedicineDose:     pp.MedicineDose,
		MedicineId:       pp.MedicineId,
	}
}

type CreatePatientProphylaxisParams struct {
	ActionContext
	PatientId   string
	Prophylaxis Prophylaxis `json:"joints_evaluation"`
}

type CreatePatientProphylaxisPayload struct {
}

func (a *Actions) CreatePatientProphylaxis(params CreatePatientProphylaxisParams) (CreatePatientProphylaxisPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWritePatient) {
		return CreatePatientProphylaxisPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetPatientByPublicId(params.PatientId)
	if err != nil {
		return CreatePatientProphylaxisPayload{}, err
	}

	je := params.Prophylaxis.IntoModel()
	je.PatientId = patient.Id

	_, err = a.app.CreateProphylaxis(je)
	if err != nil {
		return CreatePatientProphylaxisPayload{}, err
	}

	return CreatePatientProphylaxisPayload{}, nil
}

type ListPatientProphylaxesParams struct {
	ActionContext
	PatientId string
}

type ListPatientProphylaxesPayload struct {
	Data []Prophylaxis `json:"data"`
}

func (a *Actions) ListPatientProphylaxes(params ListPatientProphylaxesParams) (ListPatientProphylaxesPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadPatient) {
		return ListPatientProphylaxesPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetPatientByPublicId(params.PatientId)
	if err != nil {
		return ListPatientProphylaxesPayload{}, err
	}

	prophylaxes, err := a.app.ListProphylaxesForPatient(patient.Id)
	if err != nil {
		return ListPatientProphylaxesPayload{}, err
	}

	outProphylaxes := make([]Prophylaxis, 0, len(prophylaxes))
	for _, pp := range prophylaxes {
		outPP := new(Prophylaxis)
		outPP.FromModel(pp)
		outProphylaxes = append(outProphylaxes, *outPP)
	}

	return ListPatientProphylaxesPayload{
		Data: outProphylaxes,
	}, nil
}
