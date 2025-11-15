package app

import "shs/app/models"

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
