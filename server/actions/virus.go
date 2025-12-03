package actions

import "shs/app/models"

type Virus struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	BloodTestIds []uint `json:"blood_test_ids"`
	// TODO: expose blood tests as a whole
}

type CreateVirusParams struct {
	ActionContext
	NewVirus Virus `json:"new_virus"`
}

type CreateVirusPayload struct {
}

func (a *Actions) CreateVirus(params CreateVirusParams) (CreateVirusPayload, error) {
	err := params.Account.CheckType(models.AccountTypeAdmin)
	if err != nil {
		return CreateVirusPayload{}, err
	}

	identifyingBloodTests := make([]models.BloodTest, 0, len(params.NewVirus.BloodTestIds))
	for _, btId := range params.NewVirus.BloodTestIds {
		identifyingBloodTests = append(identifyingBloodTests, models.BloodTest{
			Id: btId,
		})
	}

	_, err = a.app.CreateVirus(models.Virus{
		Name:                  params.NewVirus.Name,
		IdentifyingBloodTests: identifyingBloodTests,
	})
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
	err := params.Account.CheckType(models.AccountTypeAdmin)
	if err != nil {
		return DeleteVirusPayload{}, err
	}

	err = a.app.DeleteVirus(params.VirusId)
	if err != nil {
		return DeleteVirusPayload{}, err
	}

	return DeleteVirusPayload{}, nil
}

type ListAllViriParams struct {
	ActionContext
	NewVirus Virus `json:"new_virus"`
}

type ListAllViriPayload struct {
	Data []Virus `json:"data"`
}

func (a *Actions) ListAllViri(params ListAllViriParams) (ListAllViriPayload, error) {
	err := params.Account.CheckType(models.AccountTypeAdmin, models.AccountTypeSecritary)
	if err != nil {
		return ListAllViriPayload{}, err
	}

	viri, err := a.app.ListAllViri()
	if err != nil {
		return ListAllViriPayload{}, err
	}

	outViri := make([]Virus, 0, len(viri))
	for _, virus := range viri {
		outViri = append(outViri, Virus{
			Id:   virus.Id,
			Name: virus.Name,
		})
	}

	return ListAllViriPayload{
		Data: outViri,
	}, nil
}
