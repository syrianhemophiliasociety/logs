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
	ChooseTheme             string
	DarkTheme               string
	LightTheme              string
	ChooseLanguage          string

	LoginUsername      string
	LoginPassword      string
	LoginEnterUsername string
	LoginEnterPassword string
	Login              string
	Logout             string

	NavHome       string
	NavAbout      string
	NavPrivacy    string
	NavLogin      string
	NavPatients   string
	NavBloodTests string
	NavMedicine   string
	NavViruses    string
	NavManagement string

	TabsList   string
	TabsSearch string
	TabsCreate string

	FormsSubmit   string
	FormsDelete   string
	FormsNewField string

	Virus          string
	EnterVirusName string

	Medicine          string
	EnterMedicineName string
	MedicineDose      string
	EnterMedicineDose string
	MedicineUnit      string
	EnterMedicineUnit string

	BloodTest               string
	BloodTestDetails        string
	BloodTestName           string
	EnterBloodTestName      string
	BloodTestFields         string
	BloodTestFieldName      string
	EnterBloodTestFieldName string
	BloodTestFieldUnit      string
	EnterBloodTestFieldUnit string

	Account                 string
	Accounts                string
	AccountUsername         string
	EnterAccountUsername    string
	AccountDisplayName      string
	EnterAccountDisplayName string
	AccountPassword         string
	EnterAccountPassword    string
	AccountType             string
	EnterAccountType        string
	AccountTypeSecritary    string
	AccountTypeAdmin        string
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

type language struct {
	DisplayName string
	LocaleKey   string
}

func Languages() []language {
	return []language{
		{DisplayName: "العربية", LocaleKey: "ar"},
		{DisplayName: "English", LocaleKey: "en"},
	}
}
