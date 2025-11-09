package actions

import (
	"errors"
	"shs/app/models"
)

type BloodTestField struct {
	Name     string               `json:"name"`
	Unit     models.BlootTestUnit `json:"unit"`
	MinValue uint                 `json:"min_value"`
	MaxValue uint                 `json:"max_value"`
}

type BloodTest struct {
	Name   string           `json:"name"`
	Fields []BloodTestField `json:"fields"`
}

type CreateBloodTestParams struct {
	ActionContext
	BloodTest BloodTest `json:"blood_test"`
}

type CreateBloodTestPayload struct {
}

func (a *Actions) CreateBloodTest(params CreateBloodTestParams) (CreateBloodTestPayload, error) {
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
	if err != nil {
		return CreateBloodTestPayload{}, err
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

	_, err = a.app.CreateBloodTest(bloodTest)
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
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
	if err != nil {
		return DeleteBloodTestPayload{}, err
	}

	err = a.app.DeleteBloodTest(params.BloodTestId)
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
	err := checkAccountType(params.Account, models.AccountTypeAdmin, models.AccountTypeSuperAdmin)
	if err != nil {
		return GetBloodTestPayload{}, err
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
			Name:     field.Name,
			Unit:     field.Unit,
			MinValue: field.MinValue,
			MaxValue: field.MaxValue,
		})
	}

	return BloodTest{
		Name:   bt.Name,
		Fields: btFields,
	}
}
