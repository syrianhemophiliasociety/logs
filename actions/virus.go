package actions

import "shs/app/models"

type Virus struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	BloodTestIds []uint `json:"blood_test_ids"`
	// TODO: expose blood tests as a whole
}

func (v Virus) IntoModel() models.Virus {
	identifyingBloodTests := make([]models.BloodTest, 0, len(v.BloodTestIds))
	for _, btId := range v.BloodTestIds {
		identifyingBloodTests = append(identifyingBloodTests, models.BloodTest{
			Id: btId,
		})
	}

	return models.Virus{
		Name:                  v.Name,
		IdentifyingBloodTests: identifyingBloodTests,
	}
}

func (v *Virus) FromModel(virus models.Virus) {
	(*v) = Virus{
		Id:   virus.Id,
		Name: virus.Name,
	}
}

type CreateVirusParams struct {
	ActionContext
	NewVirus Virus `json:"new_virus"`
}

type CreateVirusPayload struct {
}

func (a *Actions) CreateVirus(params CreateVirusParams) (CreateVirusPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteVirus) {
		return CreateVirusPayload{}, ErrPermissionDenied{}
	}

	_, err := a.app.CreateVirus(params.NewVirus.IntoModel())
	if err != nil {
		return CreateVirusPayload{}, err
	}

	return CreateVirusPayload{}, nil
}

type DeleteVirusParams struct {
	ActionContext
	VirusId uint
}

type DeleteVirusPayload struct {
}

func (a *Actions) DeleteVirus(params DeleteVirusParams) (DeleteVirusPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteVirus) {
		return DeleteVirusPayload{}, ErrPermissionDenied{}
	}

	err := a.app.DeleteVirus(params.VirusId)
	if err != nil {
		return DeleteVirusPayload{}, err
	}

	return DeleteVirusPayload{}, nil
}

type ListAllVirusesParams struct {
	ActionContext
	NewVirus Virus `json:"new_virus"`
}

type ListAllVirusesPayload struct {
	Data []Virus `json:"data"`
}

func (a *Actions) ListAllViruses(params ListAllVirusesParams) (ListAllVirusesPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadVirus) {
		return ListAllVirusesPayload{}, ErrPermissionDenied{}
	}

	viruses, err := a.app.ListAllViruses()
	if err != nil {
		return ListAllVirusesPayload{}, err
	}

	outViruses := make([]Virus, 0, len(viruses))
	for _, virus := range viruses {
		outVirus := new(Virus)
		outVirus.FromModel(virus)
		outViruses = append(outViruses, *outVirus)
	}

	return ListAllVirusesPayload{
		Data: outViruses,
	}, nil
}
