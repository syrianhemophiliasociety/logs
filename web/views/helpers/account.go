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
	t, ok := ctx.Value(auth.CtxAccountTypeKey).(string)
	if !ok {
		return accountType{"N/A"}
	}

	return accountType{t}
}
