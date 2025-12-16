package app

import (
	"shs/app/models"
	"time"
)

type Repository interface {
	GetAccount(id uint) (models.Account, error)
	GetAccountByUsername(username string) (models.Account, error)
	CreateAccount(account models.Account) (models.Account, error)
	ListAllAccounts() ([]models.Account, error)
	DeleteAccount(id uint) error
	UpdateAccountPermissions(id uint, permissions models.AccountPermissions) error
	UpdateAccountDisplayName(id uint, name string) error
	UpdateAccountPassword(id uint, password string) error
	UpdateAccountUsername(id uint, username string) error

	CreateBloodTest(bt models.BloodTest) (models.BloodTest, error)
	DeleteBloodTest(id uint) error
	GetBloodTest(id uint) (models.BloodTest, error)
	UpdateBloodTest(id uint, bt models.BloodTest) (models.BloodTest, error)
	ListAllBloodTests() ([]models.BloodTest, error)

	CreateBloodTestResult(btResult models.BloodTestResult) (models.BloodTestResult, error)
	ListPatientBloodTestResults(patientId uint) ([]models.BloodTestResult, error)
	SetBloodTestResultPending(id uint, pending bool) error
	CreateBloodTestResultFilledFields(filledFields []models.BloodTestFilledField) error
	UpdateBloodTestResultCreatedAt(id uint, ts time.Time) error

	CreateVirus(virus models.Virus) (models.Virus, error)
	DeleteVirus(id uint) error
	ListAllViruses() ([]models.Virus, error)
	ListVirusesForPatient(patientId uint) ([]models.Virus, error)

	CreateMedicine(medicine models.Medicine) (models.Medicine, error)
	DeleteMedicine(id uint) error
	ListAllMedicines() ([]models.Medicine, error)
	ListMedicinesByIds(ids []uint) ([]models.Medicine, error)
	UpdateMedicineAmount(id uint, newAmount int) error
	DecrementMedicineAmount(id uint, amount int) error
	GetMedicine(id uint) (models.Medicine, error)

	CreatePatient(patient models.Patient) (models.Patient, error)
	GetPatientById(id uint) (models.Patient, error)
	GetPatientByPublicId(publicId string) (models.Patient, error)
	FindPatientsByVisitDateRange(from, to time.Time) ([]models.Patient, error)
	FindPatientsByFields(patientIndexFields models.PatientIndexFields) ([]models.Patient, error)
	ListLastPatients(limit int) ([]models.Patient, error)
	DeletePatient(id uint) error

	CreatePatientVisit(visit models.Visit) (models.Visit, error)
	ListPatientVisits(patientId uint) ([]models.Visit, error)
	GetPatientVisit(visitId uint) (models.Visit, error)
	CreatePrescribedMedicine(pm models.PrescribedMedicine) (models.PrescribedMedicine, error)
	GetPatientLastVisit(patientId uint) (models.Visit, error)
	ListPatientVisitPrescribedMedicine(visitId uint) ([]models.PrescribedMedicine, error)
	UseMedicineForVisit(prescribedMedicineId, visitId uint) error

	CreateAddress(address models.Address) (models.Address, error)
	GetAllAddresses() ([]models.Address, error)
	GetAllAddressesALike(searchAddress models.Address) ([]models.Address, error)
	DeleteAddress(id uint) error

	CreateJointEvaluation(je models.JointsEvaluation) (models.JointsEvaluation, error)
	ListJointEvaluationsForPatient(patientId uint) ([]models.JointsEvaluation, error)

	CreateDiagnosis(d models.Diagnosis) (models.Diagnosis, error)
	DeleteDiagnisis(id uint) error
	ListAllDiagnoses() ([]models.Diagnosis, error)

	CreateDiagnosisResult(dr models.DiagnosisResult) (models.DiagnosisResult, error)
	ListPatientDiagnosisResults(patientId uint) ([]models.DiagnosisResult, error)
}
