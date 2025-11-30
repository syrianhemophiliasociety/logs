package app

import (
	"shs/app/models"
	"time"
)

func (a *App) CreatePatient(patient models.Patient) (models.Patient, error) {
	return a.repo.CreatePatient(patient)
}

func (a *App) GetPatientById(id uint) (models.Patient, error) {
	patient, err := a.repo.GetPatientById(id)
	if err != nil {
		return models.Patient{}, err
	}

	viri, err := a.repo.ListViriForPatient(patient.Id)
	if err != nil {
		return models.Patient{}, err
	}

	patient.Viri = viri

	bloodTests, err := a.repo.ListPatientBloodTestResults(patient.Id)
	if err != nil {
		return models.Patient{}, err
	}

	patient.BloodTests = bloodTests

	return patient, nil
}

func (a *App) GetMinimalPatientByPublicId(publicId string) (models.Patient, error) {
	patient, err := a.repo.GetPatientByPublicId(publicId)
	if err != nil {
		return models.Patient{}, err
	}

	return patient, nil
}

func (a *App) GetFullPatientByPublicId(publicId string) (models.Patient, error) {
	patient, err := a.repo.GetPatientByPublicId(publicId)
	if err != nil {
		return models.Patient{}, err
	}

	viri, err := a.repo.ListViriForPatient(patient.Id)
	if err != nil {
		return models.Patient{}, err
	}

	patient.Viri = viri

	bloodTests, err := a.repo.ListPatientBloodTestResults(patient.Id)
	if err != nil {
		return models.Patient{}, err
	}

	patient.BloodTests = bloodTests

	return patient, nil
}

func (a *App) FindPatientsByVisitRange(from, to time.Time) ([]models.Patient, error) {
	return a.repo.FindPatientsByVisitDateRange(from, to)
}

func (a *App) FindPatientsByIndexFields(fields models.PatientIndexFields) ([]models.Patient, error) {
	return a.repo.FindPatientsByFields(fields)
}

func (a *App) ListPatientVisitPrescribedMedicine(visitId uint) ([]models.PrescribedMedicine, error) {
	return a.repo.ListPatientVisitPrescribedMedicine(visitId)
}

func (a *App) UseMedicineForVisit(prescribedMedicineId, visitId uint) error {
	return a.repo.UseMedicineForVisit(prescribedMedicineId, visitId)
}

func (a *App) GetPatientLastVisit(patientId uint) (models.Visit, error) {
	return a.repo.GetPatientLastVisit(patientId)
}
