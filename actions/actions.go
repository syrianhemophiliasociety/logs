package actions

import "shs/app"

type Actions struct {
	app   *app.App
	cache Cache
	jwt   JwtManager[TokenPayload]
}

func New(
	app *app.App,
	cache Cache,
	jwt JwtManager[TokenPayload],
) *Actions {
	return &Actions{
		app:   app,
		cache: cache,
		jwt:   jwt,
	}
}
