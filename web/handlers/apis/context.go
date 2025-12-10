package apis

import (
	"context"
	"net/http"
	"shs-web/actions"
	"shs-web/errors"
	"shs-web/handlers/middlewares/auth"
)

func parseContext(ctx context.Context) (actions.RequestContext, error) {
	sessionToken, sessionTokenCorrect := ctx.Value(auth.CtxSessionTokenKey).(string)
	if !sessionTokenCorrect {
		return actions.RequestContext{}, errors.ErrInvalidSessionToken
	}

	return actions.RequestContext{
		SessionToken: sessionToken,
	}, nil
}

func writeRawTextResponse(w http.ResponseWriter, msg string) error {
	w.Header().Set("HX-Trigger", `{"respDetails": "`+msg+`"}`)
	w.Write([]byte(msg))
	return nil
}
