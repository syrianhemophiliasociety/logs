package helpers

import "context"

type TranslatedStringParams struct {
	CTX     context.Context
	English string
	Arabic  string
}

func TranslatedString(params TranslatedStringParams) string {
	str := params.English
	localeKey, _ := params.CTX.Value("locale").(string)
	if localeKey == "ar" {
		str = params.Arabic
	}

	return str
}
