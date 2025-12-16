package mariadb

import (
	"errors"
	"fmt"
	"shs/app"
	"shs/app/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	client *gorm.DB
}

func New() (*Repository, error) {
	conn, err := dbConnector()
	if err != nil {
		return nil, err
	}

	return &Repository{
		client: conn,
	}, nil
}

// --------------------------------
// App Repository
// --------------------------------

func (r *Repository) GetAccount(id uint) (models.Account, error) {
	var account models.Account

	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			First(&account, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Account{}, &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (r *Repository) GetAccountByUsername(username string) (models.Account, error) {
	var account models.Account

	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			First(&account, "username = ?", username).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Account{}, &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (r *Repository) CreateAccount(account models.Account) (models.Account, error) {
	account.CreatedAt = time.Now().UTC()
	account.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			Create(&account).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Account{}, &app.ErrExists{
			ResourceName: "account",
		}
	}
	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func (r *Repository) ListAllAccounts() ([]models.Account, error) {
	var accounts []models.Account

	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			Where("type NOT IN ?", []models.AccountType{models.AccountTypeSuperAdmin, models.AccountTypePatient}).
			Find(&accounts).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *Repository) DeleteAccount(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			Delete(&models.Account{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateAccountPermissions(id uint, permissions models.AccountPermissions) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			Where("id = ?", id).
			Update("permissions", permissions).
			Update("updated_at", time.Now().UTC()).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateAccountDisplayName(id uint, name string) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			Where("id = ?", id).
			Update("display_name", name).
			Update("updated_at", time.Now().UTC()).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateAccountUsername(id uint, username string) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			Where("id = ?", id).
			Update("username", username).
			Update("updated_at", time.Now().UTC()).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateAccountPassword(id uint, password string) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Account)).
			Where("id = ?", id).
			Update("password", password).
			Update("updated_at", time.Now().UTC()).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "account",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateBloodTest(bt models.BloodTest) (models.BloodTest, error) {
	bt.CreatedAt = time.Now().UTC()
	bt.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTest)).
			Create(&bt).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.BloodTest{}, &app.ErrExists{
			ResourceName: "blood_test",
		}
	}
	if err != nil {
		return models.BloodTest{}, err
	}

	return bt, nil
}

func (r *Repository) DeleteBloodTest(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTestField)).
			Delete(&models.BloodTestField{BloodTestId: id}, "blood_test_id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "blood_test",
		}
	}
	if err != nil {
		return err
	}

	err = tryWrapDbError(
		r.client.
			Model(new(models.BloodTest)).
			Delete(&models.BloodTest{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "blood_test",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBloodTest(id uint) (models.BloodTest, error) {
	var bt models.BloodTest

	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTest)).
			Preload("Fields").
			First(&bt, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.BloodTest{}, &app.ErrNotFound{
			ResourceName: "blood_test",
		}
	}
	if err != nil {
		return models.BloodTest{}, err
	}

	return bt, nil
}

func (r *Repository) UpdateBloodTest(id uint, bt models.BloodTest) (models.BloodTest, error) {
	return models.BloodTest{}, errors.New("not implemented")
}

func (r *Repository) ListAllBloodTests() ([]models.BloodTest, error) {
	var bloodTests []models.BloodTest

	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTest)).
			Preload("Fields").
			Find(&bloodTests).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return bloodTests, nil
}

func (r *Repository) CreateBloodTestResult(btResult models.BloodTestResult) (models.BloodTestResult, error) {
	btResult.CreatedAt = time.Now().UTC()
	btResult.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTestResult)).
			Create(&btResult).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.BloodTestResult{}, &app.ErrExists{
			ResourceName: "blood_test_result",
		}
	}
	if err != nil {
		return models.BloodTestResult{}, err
	}

	return btResult, nil
}

func (r *Repository) ListPatientBloodTestResults(patientId uint) ([]models.BloodTestResult, error) {
	var bloodTestResults []models.BloodTestResult

	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTestResult)).
			Preload("FilledFields").
			Where("patient_id = ?", patientId).
			Find(&bloodTestResults).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return bloodTestResults, nil
}

func (r *Repository) SetBloodTestResultPending(id uint, pending bool) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTestResult)).
			Where("id = ?", id).
			Update("pending", pending).
			Error,
	)

	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "blood_test_result",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateBloodTestResultFilledFields(filledFields []models.BloodTestFilledField) error {
	for i := range filledFields {
		filledFields[i].CreatedAt = time.Now().UTC()
		filledFields[i].UpdatedAt = time.Now().UTC()
	}

	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTestFilledField)).
			Create(filledFields).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return &app.ErrExists{
			ResourceName: "blood_test_result_field",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateBloodTestResultCreatedAt(id uint, ts time.Time) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.BloodTestResult)).
			Where("id = ?", id).
			Update("created_at", ts).
			Error,
	)

	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "blood_test_result",
		}
	}
	if err != nil {
		return err
	}

	return nil

}

func (r *Repository) CreateVirus(virus models.Virus) (models.Virus, error) {
	virus.CreatedAt = time.Now().UTC()
	virus.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Virus)).
			Create(&virus).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Virus{}, &app.ErrExists{
			ResourceName: "virus",
		}
	}
	if err != nil {
		return models.Virus{}, err
	}

	return virus, nil
}

func (r *Repository) DeleteVirus(id uint) error {
	err := tryWrapDbError(
		r.client.
			Exec("DELETE FROM identifying_blood_tests WHERE virus_id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "virus",
		}
	}
	if err != nil {
		return err
	}

	err = tryWrapDbError(
		r.client.
			Model(new(models.Virus)).
			Delete(&models.Virus{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "virus",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListAllViruses() ([]models.Virus, error) {
	var viruses []models.Virus

	err := tryWrapDbError(
		r.client.
			Model(new(models.Virus)).
			Preload("IdentifyingBloodTests").
			Find(&viruses).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return viruses, nil
}

func (r *Repository) ListVirusesForPatient(patientId uint) ([]models.Virus, error) {
	viruses := make([]models.Virus, 0)

	query := fmt.Sprintf(`SELECT %s.id, viruses.name
	FROM viruses
		JOIN has_viruses ON %s.id = has_viruses.virus_id
	WHERE has_viruses.patient_id = ?`, models.Virus{}.TableName(), models.Virus{}.TableName())

	err := tryWrapDbError(
		r.client.
			Raw(query, patientId).
			Find(&viruses).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return viruses, nil
}

func (r *Repository) CreateMedicine(medicine models.Medicine) (models.Medicine, error) {
	medicine.CreatedAt = time.Now().UTC()
	medicine.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Medicine)).
			Create(&medicine).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Medicine{}, &app.ErrExists{
			ResourceName: "medicine",
		}
	}
	if err != nil {
		return models.Medicine{}, err
	}

	return medicine, nil

}

func (r *Repository) DeleteMedicine(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Medicine)).
			Delete(&models.Medicine{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "medicine",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListAllMedicines() ([]models.Medicine, error) {
	var medicines []models.Medicine

	err := tryWrapDbError(
		r.client.
			Model(new(models.Medicine)).
			Find(&medicines).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return medicines, nil
}

func (r *Repository) ListMedicinesByIds(ids []uint) ([]models.Medicine, error) {
	var medicines []models.Medicine

	err := tryWrapDbError(
		r.client.
			Model(new(models.Medicine)).
			Where("id IN ?", ids).
			Find(&medicines).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return medicines, nil
}

func (r *Repository) UpdateMedicineAmount(id uint, newAmount int) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Medicine)).
			Where("id = ?", id).
			Update("amount", newAmount).
			Error,
	)

	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "medicine",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetMedicine(id uint) (models.Medicine, error) {
	var medicine models.Medicine

	err := tryWrapDbError(
		r.client.
			Model(new(models.Medicine)).
			First(&medicine, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Medicine{}, &app.ErrNotFound{
			ResourceName: "medicine",
		}
	}
	if err != nil {
		return models.Medicine{}, err
	}

	return medicine, nil
}

func (r *Repository) DecrementMedicineAmount(id uint, amount int) error {
	err := tryWrapDbError(
		r.client.
			Exec(fmt.Sprintf("UPDATE %s SET amount = amount - ? WHERE id = ?;", models.Medicine{}.TableName()), amount, id).
			Error,
	)

	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "medicine",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) findOrCreateLastPatientId() (models.PatientId, error) {
	var patientIds []models.PatientId
	err := tryWrapDbError(
		r.client.
			Model(new(models.PatientId)).
			Order("id DESC").
			Limit(1).
			Find(&patientIds).
			Error,
	)
	lastPatientId := models.PatientId{
		PublicId: 1,
	}
	if len(patientIds) > 0 {
		lastPatientId = patientIds[0]
	}
	if err != nil {
		err = tryWrapDbError(
			r.client.
				Model(new(models.PatientId)).
				Create(&lastPatientId).
				Error,
		)
		return lastPatientId, err
	}

	err = tryWrapDbError(
		r.client.
			Model(new(models.PatientId)).
			Create(&models.PatientId{
				PublicId: lastPatientId.PublicId + 1,
			}).
			Error,
	)
	if err != nil {
		return models.PatientId{}, err
	}

	return lastPatientId, nil
}

func (r *Repository) CreatePatient(patient models.Patient) (models.Patient, error) {
	lastPatientId, err := r.findOrCreateLastPatientId()
	if err != nil {
		return models.Patient{}, err
	}

	patient.PublicId = fmt.Sprintf("%06d", lastPatientId.PublicId)
	patient.CreatedAt = time.Now().UTC()
	patient.UpdatedAt = time.Now().UTC()

	if patient.NationalId == "" {
		patient.NationalId = "please_change_" + patient.PublicId
	}

	err = tryWrapDbError(
		r.client.
			Model(new(models.Patient)).
			Create(&patient).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Patient{}, &app.ErrExists{
			ResourceName: "patient",
		}
	}
	if err != nil {
		return models.Patient{}, err
	}

	return patient, nil
}

func (r *Repository) GetPatientById(id uint) (models.Patient, error) {
	var patient models.Patient

	err := tryWrapDbError(
		r.client.
			Model(new(models.Patient)).
			Preload("Residency").
			Preload("PlaceOfBirth").
			First(&patient, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Patient{}, &app.ErrNotFound{
			ResourceName: "patient",
		}
	}
	if err != nil {
		return models.Patient{}, err
	}

	return patient, nil
}

func (r *Repository) GetPatientByPublicId(publicId string) (models.Patient, error) {
	var patient models.Patient

	err := tryWrapDbError(
		r.client.
			Model(new(models.Patient)).
			Preload("Residency").
			Preload("PlaceOfBirth").
			First(&patient, "public_id = ?", publicId).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Patient{}, &app.ErrNotFound{
			ResourceName: "patient",
		}
	}
	if err != nil {
		return models.Patient{}, err
	}

	return patient, nil
}

func (r *Repository) FindPatientsByVisitDateRange(from, to time.Time) ([]models.Patient, error) {
	return nil, errors.New("not inmplemented")
}

func (r *Repository) FindPatientsByFields(patientIndexFields models.PatientIndexFields) ([]models.Patient, error) {
	findQuery := make([]string, 0, 9)
	findArgs := make([]any, 0, 9)
	if patientIndexFields.FirstName != "" {
		findQuery = append(findQuery, "LOWER(first_name) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(patientIndexFields.FirstName))
	}
	if patientIndexFields.LastName != "" {
		findQuery = append(findQuery, "LOWER(last_name) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(patientIndexFields.LastName))
	}
	if patientIndexFields.FatherName != "" {
		findQuery = append(findQuery, "LOWER(father_name) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(patientIndexFields.FatherName))
	}
	if patientIndexFields.MotherName != "" {
		findQuery = append(findQuery, "LOWER(mother_name) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(patientIndexFields.MotherName))
	}
	if patientIndexFields.PhoneNumber != "" {
		findQuery = append(findQuery, "LOWER(phone_number) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(patientIndexFields.PhoneNumber))
	}
	if patientIndexFields.NationalId != "" {
		findQuery = append(findQuery, "national_id = ?")
		findArgs = append(findArgs, patientIndexFields.NationalId)
	}
	if patientIndexFields.PublicId != "" {
		findQuery = append(findQuery, "public_id = ?")
		findArgs = append(findArgs, patientIndexFields.PublicId)
	}

	var patients []models.Patient

	err := tryWrapDbError(
		r.client.
			Model(new(models.Patient)).
			Where(strings.Join(findQuery, " AND "), findArgs...).
			Find(&patients).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "patient",
		}
	}
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *Repository) ListLastPatients(limit int) ([]models.Patient, error) {
	var patients []models.Patient

	err := tryWrapDbError(
		r.client.
			Model(new(models.Patient)).
			Order("created_at DESC").
			Limit(limit).
			Find(&patients).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *Repository) DeletePatient(id uint) error {
	err := tryWrapDbError(
		r.client.
			Exec("DELETE FROM did_blood_tests WHERE patient_id = ?", id).
			Error,
	)
	if err != nil {
		return err
	}

	err = tryWrapDbError(
		r.client.
			Exec("DELETE FROM has_viruses WHERE patient_id = ?", id).
			Error,
	)
	if err != nil {
		return err
	}

	err = tryWrapDbError(
		r.client.
			Model(new(models.Patient)).
			Delete(&models.Patient{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "patient",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreatePatientVisit(visit models.Visit) (models.Visit, error) {
	visit.CreatedAt = time.Now().UTC()
	visit.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Visit)).
			Create(&visit).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Visit{}, &app.ErrExists{
			ResourceName: "visit",
		}
	}
	if err != nil {
		return models.Visit{}, err
	}

	return visit, nil
}

func (r *Repository) ListPatientVisits(patientId uint) ([]models.Visit, error) {
	var visits []models.Visit

	err := tryWrapDbError(
		r.client.
			Model(new(models.Visit)).
			Where("patient_id = ?", patientId).
			Find(&visits).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return visits, nil
}

func (r *Repository) GetPatientVisit(visitId uint) (models.Visit, error) {
	var visit models.Visit

	err := tryWrapDbError(
		r.client.
			Model(new(models.Visit)).
			Where("id = ?", visitId).
			First(&visit).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Visit{}, &app.ErrNotFound{
			ResourceName: "visit",
		}
	}
	if err != nil {
		return models.Visit{}, err
	}

	return visit, nil
}

func (r *Repository) CreatePrescribedMedicine(pm models.PrescribedMedicine) (models.PrescribedMedicine, error) {
	pm.CreatedAt = time.Now().UTC()
	pm.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.PrescribedMedicine)).
			Create(&pm).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.PrescribedMedicine{}, &app.ErrExists{
			ResourceName: "prescribed_medicine",
		}
	}
	if err != nil {
		return models.PrescribedMedicine{}, err
	}

	return pm, nil
}

func (r *Repository) GetPatientLastVisit(patientId uint) (models.Visit, error) {
	var visits []models.Visit

	err := tryWrapDbError(
		r.client.
			Model(new(models.Visit)).
			Where("patient_id = ?", patientId).
			Order("created_at DESC").
			Limit(1).
			Find(&visits).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return models.Visit{}, &app.ErrNotFound{
			ResourceName: "visit",
		}
	}
	if err != nil {
		return models.Visit{}, err
	}

	if len(visits) == 0 {
		return models.Visit{}, &app.ErrNotFound{
			ResourceName: "visit",
		}
	}

	return visits[0], nil
}

func (r *Repository) ListPatientVisitPrescribedMedicine(visitId uint) ([]models.PrescribedMedicine, error) {
	var prescribedMeds []models.PrescribedMedicine

	err := tryWrapDbError(
		r.client.
			Model(new(models.PrescribedMedicine)).
			Where("visit_id = ?", visitId).
			Find(&prescribedMeds).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return prescribedMeds, nil
}

func (r *Repository) UseMedicineForVisit(prescribedMedicineId, visitId uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.PrescribedMedicine)).
			Where("id = ? AND visit_id = ?", prescribedMedicineId, visitId).
			Update("used_at", time.Now().UTC()).
			Error,
	)

	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "prescribed_medicine",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateAddress(address models.Address) (models.Address, error) {
	address.CreatedAt = time.Now().UTC()
	address.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Address)).
			Create(&address).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Address{}, &app.ErrExists{
			ResourceName: "address",
		}
	}
	if err != nil {
		return models.Address{}, err
	}

	return address, nil
}

func (r *Repository) GetAllAddresses() ([]models.Address, error) {
	var addresses []models.Address

	err := tryWrapDbError(
		r.client.
			Model(new(models.Address)).
			Find(&addresses).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "address",
		}
	}
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *Repository) GetAllAddressesALike(searchAddress models.Address) ([]models.Address, error) {
	findQuery := make([]string, 0, 3)
	findArgs := make([]any, 0, 3)
	if searchAddress.Governorate != "" {
		findQuery = append(findQuery, "LOWER(governorate) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(searchAddress.Governorate))
	}
	if searchAddress.Suburb != "" {
		findQuery = append(findQuery, "LOWER(suburb) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(searchAddress.Suburb))
	}
	if searchAddress.Street != "" {
		findQuery = append(findQuery, "LOWER(street) LIKE LOWER(?)")
		findArgs = append(findArgs, likeArg(searchAddress.Street))
	}

	var addresses []models.Address

	err := tryWrapDbError(
		r.client.
			Model(new(models.Address)).
			Where(strings.Join(findQuery, " AND "), findArgs...).
			Find(&addresses).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return nil, &app.ErrNotFound{
			ResourceName: "address",
		}
	}
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *Repository) DeleteAddress(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Address)).
			Delete(&models.Address{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "address",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateJointEvaluation(je models.JointsEvaluation) (models.JointsEvaluation, error) {
	je.CreatedAt = time.Now().UTC()
	je.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.JointsEvaluation)).
			Create(&je).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.JointsEvaluation{}, &app.ErrExists{
			ResourceName: "joints_evaluation",
		}
	}
	if err != nil {
		return models.JointsEvaluation{}, err
	}

	return je, nil
}

func (r *Repository) ListJointEvaluationsForPatient(patientId uint) ([]models.JointsEvaluation, error) {
	var jes []models.JointsEvaluation

	err := tryWrapDbError(
		r.client.
			Model(new(models.JointsEvaluation)).
			Where("patient_id = ?", patientId).
			Find(&jes).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return jes, nil
}

func (r *Repository) CreateDiagnosis(diagnosis models.Diagnosis) (models.Diagnosis, error) {
	diagnosis.CreatedAt = time.Now().UTC()
	diagnosis.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.Diagnosis)).
			Create(&diagnosis).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.Diagnosis{}, &app.ErrExists{
			ResourceName: "diagnosis",
		}
	}
	if err != nil {
		return models.Diagnosis{}, err
	}

	return diagnosis, nil
}

func (r *Repository) DeleteDiagnisis(id uint) error {
	err := tryWrapDbError(
		r.client.
			Model(new(models.Diagnosis)).
			Delete(&models.Diagnosis{Id: id}, "id = ?", id).
			Error,
	)
	if _, ok := err.(*ErrRecordNotFound); ok {
		return &app.ErrNotFound{
			ResourceName: "diagnosis",
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListAllDiagnoses() ([]models.Diagnosis, error) {
	var diagnoses []models.Diagnosis

	err := tryWrapDbError(
		r.client.
			Model(new(models.Diagnosis)).
			Find(&diagnoses).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return diagnoses, nil
}

func (r *Repository) CreateDiagnosisResult(diagnosis models.DiagnosisResult) (models.DiagnosisResult, error) {
	diagnosis.CreatedAt = time.Now().UTC()
	diagnosis.UpdatedAt = time.Now().UTC()

	err := tryWrapDbError(
		r.client.
			Model(new(models.DiagnosisResult)).
			Create(&diagnosis).
			Error,
	)
	if _, ok := err.(*ErrRecordExists); ok {
		return models.DiagnosisResult{}, &app.ErrExists{
			ResourceName: "diagnosis_result",
		}
	}
	if err != nil {
		return models.DiagnosisResult{}, err
	}

	return diagnosis, nil
}

func (r *Repository) ListPatientDiagnosisResults(patientId uint) ([]models.DiagnosisResult, error) {
	var diagnoses []models.DiagnosisResult

	err := tryWrapDbError(
		r.client.
			Model(new(models.DiagnosisResult)).
			Where("patient_id = ?", patientId).
			Find(&diagnoses).
			Error,
	)
	if err != nil {
		return nil, err
	}

	return diagnoses, nil
}

func likeArg(arg string) string {
	return fmt.Sprintf("%%%s%%", arg)
}
