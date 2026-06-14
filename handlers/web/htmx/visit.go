package htmx

import (
	"encoding/json"
	"errors"
	"net/http"
	"shs/actions"
	"shs/app"
	"shs/handlers/web/context"
	"shs/log"
	"shs/web/i18n"
	"shs/web/views/components"
	"time"
)

type visitHtmx struct {
	usecases *actions.Actions
}

func NewVisitHtmx(usecases *actions.Actions) *visitHtmx {
	return &visitHtmx{
		usecases: usecases,
	}
}

type findVisitsRequest struct {
	actions.ListAllVisitsParams
}

func (dr *findVisitsRequest) UnmarshalJSON(data []byte) error {
	type Alias findVisitsRequest
	aux := &struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		*Alias
	}{
		Alias: (*Alias)(dr),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.StartDate != "" {
		t, err := time.Parse(time.DateOnly, aux.StartDate)
		if err != nil {
			return err
		}
		dr.StartDate = t
	}

	if aux.EndDate != "" {
		t, err := time.Parse(time.DateOnly, aux.EndDate)
		if err != nil {
			return err
		}
		dr.EndDate = t
	}

	return nil
}

func (p *visitHtmx) HandleFindVisits(w http.ResponseWriter, r *http.Request) {
	ctx, err := context.Parse(r.Context())
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	var reqBody findVisitsRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	payload, err := p.usecases.ListAllVisits(actions.ListAllVisitsParams{
		ActionContext:     ctx,
		StartDate:         reqBody.StartDate,
		EndDate:           reqBody.EndDate,
		SortByVisitReason: reqBody.SortByVisitReason,
	})
	if errors.Is(err, app.ErrNotFound{}) {
		components.NotFoundError(i18n.StringsCtx(r.Context()).NavVisits).Render(r.Context(), w)
		return
	}
	if err != nil {
		components.GenericError(i18n.StringsCtx(r.Context()).ErrorSomethingWentWrong).Render(r.Context(), w)
		log.Errorln(err)
		return
	}

	components.VisitsBrief(payload.Data).Render(r.Context(), w)
}
