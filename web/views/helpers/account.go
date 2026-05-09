package helpers

import (
	"context"
	"shs/actions"
	"shs/handlers/middlewares/webauth"
)

func AccountCtx(ctx context.Context) actions.Account {
	account, ok := ctx.Value(webauth.CtxAccountKey).(actions.Account)
	if !ok {
		return actions.Account{
			DisplayName: "N/A",
			Username:    "N/A",
			Type:        "N/A",
		}
	}

	return account
}

type accountType struct {
	t string
}

func (a accountType) Admin() bool {
	return a.t == "admin"
}

func (a accountType) SuperAdmin() bool {
	return a.t == "superadmin"
}

func (a accountType) Secritary() bool {
	return a.t == "secritary"
}

func (a accountType) Patient() bool {
	return a.t == "patient"
}

func AccountTypeCtx(ctx context.Context) accountType {
	t, ok := ctx.Value(webauth.CtxAccountTypeKey).(string)
	if !ok {
		return accountType{"N/A"}
	}

	return accountType{t}
}
