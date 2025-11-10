package locale

import (
	"context"
	"net/http"
)

const LocaleKey = "locale"

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localeKey := "ar"
		cookie, err := r.Cookie(LocaleKey)
		if err == nil {
			localeKey = cookie.Value
		}

		ctx := context.WithValue(r.Context(), LocaleKey, localeKey)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
