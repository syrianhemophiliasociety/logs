package helpers

import (
	"context"
	"shs-web/actions"
	"shs-web/handlers/middlewares/auth"
)

func AccountCtx(ctx context.Context) actions.Account {
	account, ok := ctx.Value(auth.CtxAccountKey).(actions.Account)
	if !ok {
		return actions.Account{
			DisplayName: "N/A",
			Username:    "N/A",
			Type:        "N/A",
		}
	}

	return account
}

func AccountTypeCtx(ctx context.Context) string {
	accountType, ok := ctx.Value(auth.CtxAccountTypeKey).(string)
	if !ok {
		return "N/A"
	}

	return accountType
}
