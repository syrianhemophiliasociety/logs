package context

import (
	"context"
	"shs/actions"
	"shs/handlers/middlewares/webauth"
	"shs/handlers/web/errors"
	"shs/log"
)

func Parse(ctx context.Context) (actions.ActionContext, error) {
	sessionToken, sessionTokenCorrect := ctx.Value(webauth.CtxSessionTokenKey).(string)
	if !sessionTokenCorrect {
		log.Errorln("KURWA MISSING SESSION TOKEN")
		return actions.ActionContext{}, errors.ErrUnauthorized{}
	}
	account, accountOk := ctx.Value(webauth.CtxAccountKey).(actions.Account)
	if !accountOk {
		log.Errorln("KURWA MISSING ACCOUNT")
		return actions.ActionContext{}, errors.ErrUnauthorized{}
	}

	return actions.ActionContext{
		SessionToken: sessionToken,
		Account:      account,
	}, nil
}
