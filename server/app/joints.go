package app

import "shs/app/models"

func (a *App) CreateJointsEvaluation(je models.JointsEvaluation) (models.JointsEvaluation, error) {
	return a.repo.CreateJointEvaluation(je)
}

func (a *App) ListPatientJointsEvaluations(patientId uint) ([]models.JointsEvaluation, error) {
	return a.repo.ListJointEvaluationsForPatient(patientId)
}
