package i18n

import (
	"fmt"
	"time"
)

var arabic = Keys{
	Title:                   "الجمعية السورية لمرضى الناعور",
	Description:             "سجل الاضطرابات النزفية للجمعية السورية لمرضى الناعور",
	Hello:                   "مرحبا",
	ErrorSomethingWentWrong: "حدث خطأ ما",
	ErrorPermissionDenied:   "صلاحية غير كافية",
	ErrorInsufficientMedicineAmountFmt: func(medicineName string, exceedingAmount, leftPackages int) string {
		switch leftPackages {
		case 0:
			return fmt.Sprintf("لا يمكن وصف %d عبوة من %s، حيث لم يتبق منها أي عبوات.", exceedingAmount, medicineName)
		case 1:
			return fmt.Sprintf("لا يمكن وصف %d عبوة من %s، حيث لم يتبق منها سوى عبوة واحدة.", exceedingAmount, medicineName)
		default:
			return fmt.Sprintf("لا يمكن وصف %d عبوة من %s، حيث لا يزال هناك %d عبوة متبقية منه", exceedingAmount, medicineName, leftPackages)
		}
	},

	MessageSuccess: "تمت العملية بنجاح",
	MessageDeleteConfirmFmt: func(resourceType, resourceName string) string {
		return fmt.Sprintf("هل انت متأكد من حذف '%s' %s؟", resourceName, resourceType)
	},

	MessageEmptyListFmt: func(resourceType string) string {
		return fmt.Sprintf("لا يوجد %s الان", resourceType)
	},

	ChooseTheme:           "اختر لون الواجهة",
	DarkTheme:             "واجهة سوداء",
	LightTheme:            "واجهة بيضاء",
	ChooseLanguage:        "اختر اللغة",
	Yes:                   "نعم",
	No:                    "لا",
	And:                   "و",
	Or:                    "او",
	With:                  "مع",
	For:                   "ل",
	On:                    "على",
	Denied:                "مرفوضة",
	LoginUsername:         "اسم المستخدم",
	LoginPassword:         "كلمة المرور",
	LoginEnterUsername:    "أدخل اسم المستخدم",
	LoginEnterPassword:    "أدخل كلمة المرور",
	Login:                 "تسجيل الدخول",
	Logout:                "تسجيل الخروج",
	NavHome:               "الرئيسية",
	NavAbout:              "حول",
	NavPrivacy:            "الخصوصية",
	NavLogin:              "تسجيل الدخول",
	NavPatients:           "المرضى",
	NavPatient:            "بيانات المريض",
	NavBloodTests:         "التحاليل الدموية",
	NavMedicine:           "الأدوية",
	NavViruses:            "الفيروسات",
	NavManagement:         "الإدارة",
	NavAccount:            "حساب المستخدم",
	NavStatistics:         "الإحصائيات",
	NavDiagnoses:          "التشخيصات",
	TabsList:              "قائمة",
	TabsSearch:            "البحث",
	TabsCreate:            "أنشئ",
	TabsCheckup:           "المراجعة",
	TabsVisits:            "الزيارات",
	TabsJoints:            "المفاصل",
	FormsSubmit:           "أرسل المحتوى",
	FormsDelete:           "احذف المحتوى",
	FormsNewField:         "حقل جديد",
	FormsFind:             "ابحث",
	FormsUpdate:           "تعديل المحتوى",
	Virus:                 "الفيروس",
	EnterVirusName:        "أدخل اسم الفيروس",
	VirusBloodTest:        "تحليل الدم المشخص",
	EnterVirusBloodTest:   "اختر تحليل الدم المشخص",
	Medicine:              "الدواء",
	EnterMedicineName:     "أدخل اسم الدواء",
	MedicineDose:          "الجرعة",
	EnterMedicineDose:     "أدخل جرعة الدواء",
	MedicineUnit:          "الواحدة",
	EnterMedicineUnit:     "اختر واحدة الدواء",
	MedicineAmount:        "كمية الدواء",
	EnterMedicineAmount:   "ادخل كمية الدواء",
	EnterPrescribedAmount: "ادخل الكمية الموصوفة",
	MedicinePackageLeftFmt: func(n int) string {
		switch n {
		case 0:
			return "لم يتبقى اي عبوة!"
		case 1:
			return "متبقي عبوة واحدة"
		case 2:
			return "متبقي عبوتان"
		default:
			return fmt.Sprintf("متبقي %d عبوة", n)
		}
	},

	BloodTest:                   "التحليل الدموي",
	BloodTestResult:             "نتيجة التحليل الدموي",
	BloodTestDetails:            "تفاصيل التحليل الدموي",
	BloodTestName:               "اسم التحليل الدموي",
	EnterBloodTestName:          "أدخل اسم التحليل الدموي",
	BloodTestFields:             "أجزاء التحليل",
	BloodTestFieldName:          "اسم الجزء من التحليل",
	EnterBloodTestFieldName:     "أدخل اسم الجزء من التحليل",
	BloodTestFieldUnit:          "واحدة قياس الجزء من التحليل",
	EnterBloodTestFieldUnit:     "اختر واحدة قياس الجزء من التحليل",
	BloodTestFieldMinValue:      "القيمة الصغرى للمجال المرجعي",
	EnterBloodTestFieldMinValue: "أدخل القيمة الصغرى للمجال المرجعي",
	BloodTestFieldMaxValue:      "القيمة الكبرى للمجال المرجعي",
	EnterBloodTestFieldMaxValue: "أدخل القيمة الكبرى للمجال المرجعي",
	EnterBloodTestResultFieldValueFmt: func(unit string) string {
		return fmt.Sprintf("أدخل قيمة النتيجة في %s", unit)
	},

	RemoveBloodTest:               "إزالة التحليل",
	BloodTestDoLater:              "القيام بالتحليل لاحقا",
	BloodTestPending:              "التحليل غير مكتمل",
	BloodTestPendingInfo:          "لم يكتمل تحليل الدم، يرجى ملء الحقول أدناه لإكماله",
	BloodTestDate:                 "تاريخ التحليل",
	PrimaryProphylaxis:            "العلاج الوقائي الأولي",
	SecondaryProphylaxis:          "العلاج الوقائي الثانوي",
	Surgery:                       "عمل جراحي",
	JointEvaluation:               "التقييم المفصلي",
	JointInjection:                "الحقن المفصلي",
	Hemelibra:                     "هيمليبرا",
	TreatmentAtHome:               "العلاج المنزلي",
	ActiveBleeding:                "نزف حالي",
	ValueInsideRange:              "القيمة داخل المجال المرجعي",
	ValueOutsideRange:             "القيمة خارج المجال المرجعي",
	Account:                       "الاسم",
	Accounts:                      "الحسابات",
	AccountUsername:               "اسم المستخدم",
	EnterAccountUsername:          "أدخل اسم المستخدم",
	AccountDisplayName:            "الاسم",
	EnterAccountDisplayName:       "أدخل الاسم",
	AccountPassword:               "كلمة المرور",
	AccountPasswordUnchanged:      "لم يتم تغير كلمة المرور",
	EnterAccountPassword:          "أدخل كلمة المرور",
	AccountPermissions:            "الصلاحيات",
	EnterAccountPermissions:       "اختر الصلاحيات",
	PermissionRead:                "صلاحية قراءة",
	PermissionWrite:               "صلاحية كتابة",
	AccountType:                   "نوع الحساب",
	EnterAccountType:              "اختر نوع الحساب",
	AccountTypeSecritary:          "محاسبة",
	AccountTypeAdmin:              "مدير",
	AccountDelete:                 "حذف الحساب",
	Patient:                       "المريض",
	PatientFirstName:              "الاسم",
	EnterPatientFirstName:         "أدخل الاسم",
	PatientLastName:               "النسبة",
	EnterPatientLastName:          "أدخل النسبة",
	PatientFatherName:             "اسم الأب",
	EnterPatientFatherName:        "أدخل اسم الأب",
	PatientMotherName:             "اسم الأم",
	EnterPatientMotherName:        "أدخل اسم الأم",
	NationalId:                    "الرقم الوطني",
	EnterNationalId:               "أدخل الرقم الوطني",
	Nationality:                   "الجنسية",
	EnterNationality:              "أدخل الجنسية",
	PlaceOfBirth:                  "مكان الولادة",
	Governorate:                   "المحافظة",
	EnterGovernorate:              "أدخل المحافظة",
	Suburb:                        "الحي",
	EnterSuburb:                   "أدخل الحي",
	Street:                        "الشارع",
	EnterStreet:                   "أدخل الشارع",
	DateOfBirth:                   "تاريخ الولادة",
	EnterDateOfBirth:              "أدخل تاريخ الولادة",
	Residency:                     "الإقامة",
	Gender:                        "الجنس",
	EnterGender:                   "احتر الجنس",
	GenderMale:                    "ذكر",
	GenderFemale:                  "أنثى",
	PhoneNumber:                   "رقم الهاتف",
	EnterPhoneNumber:              "أدخل رقم الهاتف",
	Diagnosis:                     "التّشخيص",
	PatientSonOf:                  "ابن",
	PatientDaughterOf:             "بنت",
	FindPatients:                  "ابحث عن المرضى",
	CheckupOnPatients:             "مراجعة المرضى",
	PatientProfile:                "الملف الشخصي",
	PatientAddBloodTest:           "تسجيل نتيجة تحليل الدم",
	PatientAddVirus:               "تسجيل فيروس",
	PatientId:                     "الرقم التعريفي",
	EnterPatientId:                "الرقم التعريفي",
	EnterBloodTest:                "اختر تحليل الدم",
	FamilyHistoryExists:           "هل يوجد قصة عاىلية",
	FirstVisitReason:              "سبب اول زيارة",
	EnterFirstVisitReason:         "اختر سبب اول زيارة",
	FirstVisitReasonFamilyHistory: "قصة عاىلية",
	FirstVisitReasonBleeding:      "نزف",
	FirstVisitReasonReferral:      "احالة",
	JointsEvaluations:             "تقاييم المفاصل",
	JointsRightAnkle:              "الكاحل الأيمن",
	JointsLeftAnkle:               "الكاحل الأيسر",
	JointsRightKnee:               "الركبة اليمنى",
	JointsLeftKnee:                "الركبة اليسرى",
	JointsRightElbow:              "المرفق الأيمن",
	JointsLeftElbow:               "المرفق الأيسر",
	VisitTitleFmt: func(date time.Time) string {
		return fmt.Sprintf("الزيارة بتاريخ %s", date.Format(time.DateOnly))
	},

	UseOnePrescribedMedicineFmt: func(medicineName string) string {
		return fmt.Sprintf("استخدم واحد من %s", medicineName)
	},

	UseMedicineParagraph:            "هذه هي الأدوية الموصوفة من زيارتك الأخيرة لنا، يرجى اختيار الدواء الذي ستستخدمه والنقر على \"استخدم واحد من ...\".",
	CheckUpVisitReason:              "سبب الزيارة",
	EnterCheckUpVisitReason:         "اختر سبب الزيارة",
	CheckUpPrescribedMedicines:      "الادوية الموصوفة",
	EnterCheckUpPrescribedMedicines: "اختر الادوية الموصوفة",
	CheckUpAddPrescribedMedicine:    "اضف دواء",
	CheckUpExtraReason:              "ملاحظة اضافية",
	EnterCheckUpExtraReason:         "ادخل الملاحظة الاضافية",
	CheckUpPatientWeight:            "وزن المريض",
	EnterCheckUpPatientWeight:       "ادخل وزن المريض",
	CheckUpPatientHeight:            "طول المريض",
	EnterCheckUpPatientHeight:       "ادخل طول المريض",
	PrescribedAmount:                "الكمية الموصوفة",
	UsedAmount:                      "الكمية المستخدمة",
	PrescribedMedicineUsedAt:        "استخدم في",
	NationalitySyrian:               "سوري",
	NationalityPalestinian:          "فلسطيني",
	NationalityIraqi:                "عراقي",
	NationalityEgyptian:             "مصري",
	NationalityLebanese:             "لبناني",
}
