package actions

import (
	"errors"
	"shs/app/models"
)

type BloodTestField struct {
	Id       uint                 `json:"id"`
	Name     string               `json:"name"`
	Unit     models.BlootTestUnit `json:"unit"`
	MinValue uint                 `json:"min_value"`
	MaxValue uint                 `json:"max_value"`
}

type BloodTest struct {
	Id     uint             `json:"id"`
	Name   string           `json:"name"`
	Fields []BloodTestField `json:"fields"`
}

type CreateBloodTestParams struct {
	ActionContext
	BloodTest BloodTest `json:"new_blood_test"`
}

type CreateBloodTestPayload struct {
}

func (a *Actions) CreateBloodTest(params CreateBloodTestParams) (CreateBloodTestPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteBloodTest) {
		return CreateBloodTestPayload{}, ErrPermissionDenied{}
	}

	bloodTestFields := make([]models.BloodTestField, 0, len(params.BloodTest.Fields))
	for _, field := range params.BloodTest.Fields {
		bloodTestFields = append(bloodTestFields, models.BloodTestField{
			Name:     field.Name,
			Unit:     field.Unit,
			MinValue: field.MinValue,
			MaxValue: field.MaxValue,
		})
	}
	bloodTest := models.BloodTest{
		Name:   params.BloodTest.Name,
		Fields: bloodTestFields,
	}

	_, err := a.app.CreateBloodTest(bloodTest)
	if err != nil {
		return CreateBloodTestPayload{}, err
	}

	return CreateBloodTestPayload{}, nil
}

type UpdateBloodTestParams struct {
	ActionContext
}

type UpdateBloodTestPayload struct {
}

func (a *Actions) UpdateBloodTest(params UpdateBloodTestParams) (UpdateBloodTestPayload, error) {
	return UpdateBloodTestPayload{}, errors.New("not implemented")
}

type DeleteBloodTestParams struct {
	ActionContext
	BloodTestId uint `json:"blood_test_id"`
}

type DeleteBloodTestPayload struct {
}

func (a *Actions) DeleteBloodTest(params DeleteBloodTestParams) (DeleteBloodTestPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionWriteBloodTest) {
		return DeleteBloodTestPayload{}, ErrPermissionDenied{}
	}

	err := a.app.DeleteBloodTest(params.BloodTestId)
	if err != nil {
		return DeleteBloodTestPayload{}, err
	}

	return DeleteBloodTestPayload{}, nil
}

type GetBloodTestParams struct {
	ActionContext
	BloodTestId uint `json:"blood_test_id"`
}

type GetBloodTestPayload struct {
	Data BloodTest `json:"data"`
}

func (a *Actions) GetBloodTest(params GetBloodTestParams) (GetBloodTestPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadBloodTest) {
		return GetBloodTestPayload{}, ErrPermissionDenied{}
	}

	bt, err := a.app.GetBloodTest(params.BloodTestId)
	if err != nil {
		return GetBloodTestPayload{}, err
	}

	return GetBloodTestPayload{
		Data: mapModelBloodTest(bt),
	}, nil
}

type ListAllBloodTestsParams struct {
	ActionContext
}

type ListAllBloodTestsPayload struct {
	Data []BloodTest `json:"data"`
}

func (a *Actions) ListAllBloodTests(params ListAllBloodTestsParams) (ListAllBloodTestsPayload, error) {
	if !params.Account.HasPermission(models.AccountPermissionReadBloodTest) {
		return ListAllBloodTestsPayload{}, ErrPermissionDenied{}
	}

	bloodTests, err := a.app.ListAllBloodTests()
	if err != nil {
		return ListAllBloodTestsPayload{}, err
	}

	outBloodTests := make([]BloodTest, 0, len(bloodTests))
	for _, bt := range bloodTests {
		outBloodTests = append(outBloodTests, mapModelBloodTest(bt))
	}

	return ListAllBloodTestsPayload{
		Data: outBloodTests,
	}, nil
}

func mapModelBloodTest(bt models.BloodTest) BloodTest {
	btFields := make([]BloodTestField, 0, len(bt.Fields))
	for _, field := range bt.Fields {
		btFields = append(btFields, BloodTestField{
			Id:       field.Id,
			Name:     field.Name,
			Unit:     field.Unit,
			MinValue: field.MinValue,
			MaxValue: field.MaxValue,
		})
	}

	return BloodTest{
		Id:     bt.Id,
		Name:   bt.Name,
		Fields: btFields,
	}
}
