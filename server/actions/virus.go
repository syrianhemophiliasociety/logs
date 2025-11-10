package actions

import "shs/app/models"

type Virus struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateVirusParams struct {
	ActionContext
	NewVirus Virus `json:"new_virus"`
}

type CreateVirusPayload struct {
}

func (a *Actions) CreateVirus(params CreateVirusParams) (CreateVirusPayload, error) {
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
	if err != nil {
		return CreateVirusPayload{}, err
	}

	_, err = a.app.CreateVirus(models.Virus{
		Name: params.NewVirus.Name,
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
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
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
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
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
