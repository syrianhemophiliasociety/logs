package apis

import (
	"context"
	"shs/actions"
	"shs/handlers/middlewares/auth"
)

func parseContext(ctx context.Context) (actions.ActionContext, error) {
	account, accountCorrect := ctx.Value(auth.AccountKey).(actions.Account)
	if !accountCorrect {
		return actions.ActionContext{}, &ErrUnauthorized{}
	}

	return actions.ActionContext{
		Account: account,
	}, nil
}
