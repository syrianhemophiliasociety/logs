package actions

import (
	"shs/app/models"
	"time"
)

type JointsEvaluation struct {
	Id         uint      `json:"id"`
	RightAnkle int       `json:"right_ankle"`
	LeftAnkle  int       `json:"left_ankle"`
	RightKnee  int       `json:"right_knee"`
	LeftKnee   int       `json:"left_knee"`
	RightElbow int       `json:"right_elbow"`
	LeftElbow  int       `json:"left_elbow"`
	Result     int       `json:"result"`
	CreatedAt  time.Time `json:"created_at"`
}

func (j *JointsEvaluation) FromModel(je models.JointsEvaluation) {
	(*j).Id = je.Id
	(*j).RightAnkle = je.RightAnkle
	(*j).LeftAnkle = je.LeftAnkle
	(*j).RightKnee = je.RightKnee
	(*j).LeftKnee = je.LeftKnee
	(*j).RightElbow = je.RightElbow
	(*j).LeftElbow = je.LeftElbow
	(*j).CreatedAt = je.CreatedAt

	(*j).Result = je.RightAnkle + je.LeftAnkle +
		je.RightKnee + je.LeftKnee +
		je.RightElbow + je.LeftElbow
}

func (j JointsEvaluation) IntoModel() models.JointsEvaluation {
	return models.JointsEvaluation{
		RightAnkle: j.RightAnkle,
		LeftAnkle:  j.LeftAnkle,
		RightKnee:  j.RightKnee,
		LeftKnee:   j.LeftKnee,
		RightElbow: j.RightElbow,
		LeftElbow:  j.LeftElbow,
	}
}

type CreatePatientJointsEvaluationParams struct {
	ActionContext
	PatientId        string
	JointsEvaluation JointsEvaluation `json:"joints_evaluation"`
}

type CreatePatientJointsEvaluationPayload struct {
}

func (a *Actions) CreatePatientJointsEvaluation(params CreatePatientJointsEvaluationParams) (CreatePatientJointsEvaluationPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWritePatient) {
		return CreatePatientJointsEvaluationPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetMinimalPatientByPublicId(params.PatientId)
	if err != nil {
		return CreatePatientJointsEvaluationPayload{}, err
	}

	je := params.JointsEvaluation.IntoModel()
	je.PatientId = patient.Id

	_, err = a.app.CreateJointsEvaluation(je)
	if err != nil {
		return CreatePatientJointsEvaluationPayload{}, err
	}

	return CreatePatientJointsEvaluationPayload{}, nil
}

type ListPatientJointsEvaluationsParams struct {
	ActionContext
	PatientId string
}

type ListPatientJointsEvaluationsPayload struct {
	Data []JointsEvaluation `json:"data"`
}

func (a *Actions) ListPatientJointsEvaluations(params ListPatientJointsEvaluationsParams) (ListPatientJointsEvaluationsPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadPatient) {
		return ListPatientJointsEvaluationsPayload{}, ErrPermissionDenied{}
	}

	patient, err := a.app.GetMinimalPatientByPublicId(params.PatientId)
	if err != nil {
		return ListPatientJointsEvaluationsPayload{}, err
	}

	joints, err := a.app.ListPatientJointsEvaluations(patient.Id)
	if err != nil {
		return ListPatientJointsEvaluationsPayload{}, err
	}

	outJoints := make([]JointsEvaluation, 0, len(joints))
	for _, joint := range joints {
		outJoint := new(JointsEvaluation)
		outJoint.FromModel(joint)
		outJoints = append(outJoints, *outJoint)
	}

	return ListPatientJointsEvaluationsPayload{
		Data: outJoints,
	}, nil
}
