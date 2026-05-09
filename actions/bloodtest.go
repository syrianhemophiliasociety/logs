package actions

import (
	"errors"
	"shs/app/models"
)

type BloodTestField struct {
	Id             uint                 `json:"id"`
	Name           string               `json:"name"`
	Unit           models.BlootTestUnit `json:"unit"`
	MinValueNumber float64              `json:"min_value_number"`
	MinValueString string               `json:"min_value_string"`
	MaxValueNumber float64              `json:"max_value_number"`
	MaxValueString string               `json:"max_value_string"`
}

type BloodTest struct {
	Id     uint             `json:"id"`
	Name   string           `json:"name"`
	Fields []BloodTestField `json:"fields"`
}

func (bt BloodTest) IntoModel() models.BloodTest {
	bloodTestFields := make([]models.BloodTestField, 0, len(bt.Fields))
	for _, field := range bt.Fields {
		bloodTestFields = append(bloodTestFields, models.BloodTestField{
			Name:           field.Name,
			Unit:           field.Unit,
			MinValueNumber: field.MinValueNumber,
			MinValueString: field.MinValueString,
			MaxValueNumber: field.MaxValueNumber,
			MaxValueString: field.MaxValueString,
		})
	}
	return models.BloodTest{
		Name:   bt.Name,
		Fields: bloodTestFields,
	}
}

func (bt *BloodTest) FromModel(bloodTest models.BloodTest) {
	btFields := make([]BloodTestField, 0, len(bloodTest.Fields))
	for _, field := range bloodTest.Fields {
		btFields = append(btFields, BloodTestField{
			Id:             field.Id,
			Name:           field.Name,
			Unit:           field.Unit,
			MinValueNumber: field.MinValueNumber,
			MinValueString: field.MinValueString,
			MaxValueNumber: field.MaxValueNumber,
			MaxValueString: field.MaxValueString,
		})
	}

	(*bt) = BloodTest{
		Id:     bloodTest.Id,
		Name:   bloodTest.Name,
		Fields: btFields,
	}
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

	_, err := a.app.CreateBloodTest(params.BloodTest.IntoModel())
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

	outBt := new(BloodTest)
	outBt.FromModel(bt)

	return GetBloodTestPayload{
		Data: *outBt,
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
		outBt := new(BloodTest)
		outBt.FromModel(bt)
		outBloodTests = append(outBloodTests, *outBt)
	}

	return ListAllBloodTestsPayload{
		Data: outBloodTests,
	}, nil
}
