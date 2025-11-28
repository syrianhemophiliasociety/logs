package actions

type Actions struct {
	cache Cache
}

func New(cache Cache) *Actions {
	return &Actions{
		cache: cache,
	}
}

type RequestContext struct {
	SessionToken string
	Account      Account
}
