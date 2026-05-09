package app

import "shs/app/models"

func (a *App) CreateDiagnosis(d models.Diagnosis) (models.Diagnosis, error) {
	return a.repo.CreateDiagnosis(d)
}

func (a *App) DeleteDiagnosis(id uint) error {
	return a.repo.DeleteDiagnisis(id)
}

func (a *App) ListAllDiagnoses() ([]models.Diagnosis, error) {
	return a.repo.ListAllDiagnoses()
}

func (a *App) CreateDiagnosisResult(dr models.DiagnosisResult) (models.DiagnosisResult, error) {
	return a.repo.CreateDiagnosisResult(dr)
}

func (a *App) ListPatientDiagnosisResults(patientId uint) ([]models.DiagnosisResult, error) {
	return a.repo.ListPatientDiagnosisResults(patientId)
}
