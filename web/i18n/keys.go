package i18n

import (
	"context"
	"shs-web/handlers/middlewares/locale"
)

type Keys struct {
	Title       string
	Description string

	ErrorSomethingWentWrong string
	MessageSuccess          string

	LoginUsername      string
	LoginPassword      string
	LoginEnterUsername string
	LoginEnterPassword string
	Login              string

	NavHome       string
	NavAbout      string
	NavPrivacy    string
	NavLogin      string
	NavPatients   string
	NavBloodTests string
	NavMedicine   string
	NavViruses    string
	NavManagement string

	FormsSubmit string
	FormsDelete string

	Virus          string
	EnterVirusName string

	Medicine          string
	EnterMedicineName string
	MedicineDose      string
	EnterMedicineDose string
	MedicineUnit      string
	EnterMedicineUnit string
}

var localeKeys = map[string]Keys{
	"en": english,
	"ar": arabic,
}

func Strings(localeKey string) Keys {
	if keys, ok := localeKeys[localeKey]; ok {
		return keys
	}
	return english
}

func StringsCtx(ctx context.Context) Keys {
	localeKey, ok := ctx.Value(locale.LocaleKey).(string)
	if !ok {
		return Strings("en")
	}
	return Strings(localeKey)
}
