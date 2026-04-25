package pages

import (
	"context"
	"shs-web/actions"
	"shs-web/errors"
	"shs-web/handlers/middlewares/auth"
)

func parseContext(ctx context.Context) (actions.RequestContext, error) {
	sessionToken, sessionTokenCorrect := ctx.Value(auth.CtxSessionTokenKey).(string)
	if !sessionTokenCorrect {
		return actions.RequestContext{}, errors.ErrInvalidSessionToken
	}
	account, accountOk := ctx.Value(auth.CtxAccountKey).(actions.Account)
	if !accountOk {
		return actions.RequestContext{}, errors.ErrInvalidSessionToken
	}

	return actions.RequestContext{
		SessionToken: sessionToken,
		Account:      account,
	}, nil
}
