package apis

import (
	"encoding/json"
	"io"
	"net/http"
	"shs/actions"
	"shs/app/models"
	"shs/handlers/web/context"
	"shs/log"
	"shs/web/i18n"
	"shs/web/views/components"
	"strconv"
)

type RequestBloodTest struct {
	Id         uint     `json:"id"`
	Name       string   `json:"name"`
	FieldNames []string `json:"blood_test_field_name"`
	FieldUnits []string `json:"blood_test_field_unit"`
	MinValues  []string `json:"blood_test_field_min_value"`
	MaxValues  []string `json:"blood_test_field_max_value"`
}

type RequestBloodTestSingle struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	FieldName string `json:"blood_test_field_name"`
	FieldUnit string `json:"blood_test_field_unit"`
	MinValue  string `json:"blood_test_field_min_value"`
	MaxValue  string `json:"blood_test_field_max_value"`
}

func clusterFuckBloodTestsToActionsOne(btSingle RequestBloodTestSingle, btMulti RequestBloodTest) actions.BloodTest {
	var newBloodTest actions.BloodTest

	if btMulti.Name != "" {
		newBloodTest.Name = btMulti.Name
		for i := range len(btMulti.FieldNames) {
			minValue, _ := strconv.ParseFloat(btMulti.MinValues[i], 64)
			maxValue, _ := strconv.ParseFloat(btMulti.MaxValues[i], 64)

			newBloodTest.Fields = append(newBloodTest.Fields, actions.BloodTestField{
				Name:           btMulti.FieldNames[i],
				Unit:           models.BlootTestUnit(btMulti.FieldUnits[i]),
				MinValueString: btMulti.MinValues[i],
				MinValueNumber: minValue,
				MaxValueString: btMulti.MaxValues[i],
				MaxValueNumber: maxValue,
			})
		}
	}
	if btSingle.Name != "" {
		newBloodTest.Name = btSingle.Name
		minValue, _ := strconv.ParseFloat(btSingle.MinValue, 64)
		maxValue, _ := strconv.ParseFloat(btSingle.MaxValue, 64)
		newBloodTest.Fields = append(newBloodTest.Fields, actions.BloodTestField{
			Name:           btSingle.FieldName,
			Unit:           models.BlootTestUnit(btSingle.FieldUnit),
			MinValueString: btSingle.MinValue,
			MinValueNumber: minValue,
			MaxValueString: btSingle.MaxValue,
			MaxValueNumber: maxValue,
		})
	}

	return newBloodTest
}

///

type bloodTestApi struct {
	usecases *actions.Actions
}

func NewBloodTestApi(usecases *actions.Actions) *bloodTestApi {
	return &bloodTestApi{
		usecases: usecases,
	}
}

func (v *bloodTestApi) HandleCreateBloodTest(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody RequestBloodTest
	var reqBody2 RequestBloodTestSingle
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		err = json.Unmarshal(body, &reqBody2)
		if err != nil {
			components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
			log.Errorln(err)
			return
		}
	}

	_, err = v.usecases.CreateBloodTest(actions.CreateBloodTestParams{
		ActionContext: ctx,
		BloodTest:     clusterFuckBloodTestsToActionsOne(reqBody2, reqBody),
	})
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}

func (v *bloodTestApi) HandleDeleteBloodTest(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	id := r.PathValue("id")
	intId, _ := strconv.Atoi(id)

	_, err = v.usecases.DeleteBloodTest(actions.DeleteBloodTestParams{
		ActionContext: ctx,
		BloodTestId:   uint(intId),
	})
	if err != nil {
		writeRawTextResponse(w, i18n.Strings("en").ErrorSomethingWentWrong)
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	writeRawTextResponse(w, i18n.Strings("en").MessageSuccess)
}
