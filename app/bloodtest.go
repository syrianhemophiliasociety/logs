package app

import (
	"shs/app/models"
	"time"
)

func (a *App) CreateBloodTest(bt models.BloodTest) (models.BloodTest, error) {
	return a.repo.CreateBloodTest(bt)
}

func (a *App) GetBloodTest(id uint) (models.BloodTest, error) {
	return a.repo.GetBloodTest(id)
}

func (a *App) DeleteBloodTest(id uint) error {
	return a.repo.DeleteBloodTest(id)
}

func (a *App) ListAllBloodTests() ([]models.BloodTest, error) {
	return a.repo.ListAllBloodTests()
}

func (a *App) CreateBloodTestResult(btr models.BloodTestResult) (models.BloodTestResult, error) {
	return a.repo.CreateBloodTestResult(btr)
}

func (a *App) ListPatientBloodTestResults(patientId uint) ([]models.BloodTestResult, error) {
	return a.repo.ListPatientBloodTestResults(patientId)
}

func (a *App) UpdatePatientPendingBloodTestResultFields(btrId uint, fields []models.BloodTestFilledField) error {
	err := a.repo.SetBloodTestResultPending(btrId, false)
	if err != nil {
		return err
	}

	err = a.repo.UpdateBloodTestResultCreatedAt(btrId, time.Now().UTC())
	if err != nil {
		return err
	}

	for i := range fields {
		fields[i].BloodTestResultId = btrId
	}

	return a.repo.CreateBloodTestResultFilledFields(fields)
}
