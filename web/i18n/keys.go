package i18n

import (
	"context"
	"shs-web/handlers/middlewares/locale"
	"time"
)

type Keys struct {
	Title       string
	Description string
	Hello       string

	ErrorSomethingWentWrong            string
	ErrorPermissionDenied              string
	ErrorInsufficientMedicineAmountFmt func(medicineName string, exceedingAmount, leftPackages int) string
	MessageSuccess                     string
	MessageDeleteConfirmFmt            func(resourceType, resourceName string) string
	MessageEmptyListFmt                func(resourceType string) string
	ChooseTheme                        string
	DarkTheme                          string
	LightTheme                         string
	ChooseLanguage                     string

	Yes    string
	No     string
	And    string
	Or     string
	With   string
	For    string
	On     string
	Denied string

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
	NavAccount    string
	NavStatistics string
	NavDiagnoses  string

	TabsList    string
	TabsSearch  string
	TabsCreate  string
	TabsCheckup string
	TabsVisits  string
	TabsJoints  string

	FormsSubmit   string
	FormsDelete   string
	FormsNewField string
	FormsFind     string
	FormsUpdate   string

	Virus               string
	EnterVirusName      string
	VirusBloodTest      string
	EnterVirusBloodTest string

	Medicine               string
	EnterMedicineName      string
	MedicineDose           string
	EnterMedicineDose      string
	MedicineUnit           string
	EnterMedicineUnit      string
	MedicineAmount         string
	EnterMedicineAmount    string
	EnterPrescribedAmount  string
	MedicinePackageLeftFmt func(n int) string

	BloodTest                         string
	BloodTestResult                   string
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
	BloodTestDoLater                  string
	BloodTestPending                  string
	BloodTestPendingInfo              string
	BloodTestDate                     string
	PrimaryProphylaxis                string
	SecondaryProphylaxis              string
	Surgery                           string
	JointEvaluation                   string
	JointInjection                    string
	Hemelibra                         string
	TreatmentAtHome                   string
	ActiveBleeding                    string

	Account                  string
	Accounts                 string
	AccountUsername          string
	EnterAccountUsername     string
	AccountDisplayName       string
	EnterAccountDisplayName  string
	AccountPassword          string
	AccountPasswordUnchanged string
	EnterAccountPassword     string
	AccountPermissions       string
	EnterAccountPermissions  string
	PermissionRead           string
	PermissionWrite          string
	AccountType              string
	EnterAccountType         string
	AccountTypeSecritary     string
	AccountTypeAdmin         string
	AccountDelete            string

	Patient                       string
	PatientFirstName              string
	EnterPatientFirstName         string
	PatientLastName               string
	EnterPatientLastName          string
	PatientFatherName             string
	EnterPatientFatherName        string
	PatientMotherName             string
	EnterPatientMotherName        string
	NationalId                    string
	EnterNationalId               string
	Nationality                   string
	EnterNationality              string
	PlaceOfBirth                  string
	Governorate                   string
	EnterGovernorate              string
	Suburb                        string
	EnterSuburb                   string
	Street                        string
	EnterStreet                   string
	DateOfBirth                   string
	EnterDateOfBirth              string
	Residency                     string
	Gender                        string
	EnterGender                   string
	GenderMale                    string
	GenderFemale                  string
	PhoneNumber                   string
	EnterPhoneNumber              string
	Diagnosis                     string
	PatientSonOf                  string
	PatientDaughterOf             string
	FindPatients                  string
	CheckupOnPatients             string
	PatientProfile                string
	PatientAddBloodTest           string
	PatientAddVirus               string
	PatientId                     string
	EnterPatientId                string
	EnterBloodTest                string
	FamilyHistoryExists           string
	FirstVisitReason              string
	EnterFirstVisitReason         string
	FirstVisitReasonFamilyHistory string
	FirstVisitReasonBleeding      string
	FirstVisitReasonReferral      string
	JointsEvaluations             string
	JointsRightAnkle              string
	JointsLeftAnkle               string
	JointsRightKnee               string
	JointsLeftKnee                string
	JointsRightElbow              string
	JointsLeftElbow               string
	VisitTitleFmt                 func(date time.Time) string

	CheckUpVisitReason              string
	EnterCheckUpVisitReason         string
	CheckUpPrescribedMedicines      string
	EnterCheckUpPrescribedMedicines string
	CheckUpAddPrescribedMedicine    string
	CheckUpExtraReason              string
	EnterCheckUpExtraReason         string
	PrescribedAmount                string
	UsedAmount                      string
	PrescribedMedicineUsedAt        string

	NationalitySyrian      string
	NationalityPalestinian string
	NationalityIraqi       string
	NationalityEgyptian    string
	NationalityLebanese    string
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
