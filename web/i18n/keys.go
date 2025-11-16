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

	And  string
	Or   string
	With string

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
	NavPatient    string
	NavBloodTests string
	NavMedicine   string
	NavViruses    string
	NavManagement string

	TabsList    string
	TabsSearch  string
	TabsCreate  string
	TabsCheckup string

	FormsSubmit   string
	FormsDelete   string
	FormsNewField string
	FormsFind     string

	Virus          string
	EnterVirusName string

	Medicine          string
	EnterMedicineName string
	MedicineDose      string
	EnterMedicineDose string
	MedicineUnit      string
	EnterMedicineUnit string

	BloodTest                         string
	BloodTestDetails                  string
	BloodTestName                     string
	EnterBloodTestName                string
	BloodTestFields                   string
	BloodTestFieldName                string
	EnterBloodTestFieldName           string
	BloodTestFieldUnit                string
	EnterBloodTestFieldUnit           string
	EnterBloodTestResultFieldValueFmt func(unit string) string
	RemoveBloodTest                   string

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

	Patient                string
	PatientFirstName       string
	EnterPatientFirstName  string
	PatientLastName        string
	EnterPatientLastName   string
	PatientFatherName      string
	EnterPatientFatherName string
	PatientMotherName      string
	EnterPatientMotherName string
	NationalId             string
	EnterNationalId        string
	Nationality            string
	EnterNationality       string
	PlaceOfBirth           string
	Governorate            string
	EnterGovernorate       string
	Suburb                 string
	EnterSuburb            string
	Street                 string
	EnterStreet            string
	DateOfBirth            string
	EnterDateOfBirth       string
	Residency              string
	Gender                 string
	EnterGender            string
	GenderMale             string
	GenderFemale           string
	PhoneNumber            string
	EnterPhoneNumber       string
	Diagnosis              string
	PatientSonOf           string
	PatientDaughterOf      string

	NationalitySyrian      string
	NationalityPalestinian string
	NationalityIraqi       string
	NationalityEgyptian    string
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
