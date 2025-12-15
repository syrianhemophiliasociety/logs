package mariadb

import (
	"shs/app/models"
	"shs/config"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Migrate() error {
	dbConn, err := dbConnector()
	if err != nil {
		return err
	}

	err = dbConn.Debug().AutoMigrate(
		new(models.Account),
		new(models.Virus),
		new(models.Medicine),
		new(models.Visit),
		new(models.BloodTest),
		new(models.BloodTestResult),
		new(models.BloodTestField),
		new(models.BloodTestFilledField),
		new(models.Address),
		new(models.Patient),
		new(models.PatientId),
		new(models.PatientUseMedicine),
		new(models.PrescribedMedicine),
		new(models.JointsEvaluation),
	)
	if err != nil {
		return err
	}

	for _, tableName := range []string{
		"accounts",
		"addresses",
		"blood_test_fields",
		"blood_test_filled_fields",
		"blood_test_results",
		"blood_tests",
		"did_blood_tests",
		"has_viri",
		"medicines",
		"patient_use_medicines",
		"patients",
		"prescribed_medicines",
		"viri",
		"visits",
		"joints_evaluations",
	} {
		err = dbConn.Exec("ALTER TABLE " + tableName + " CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error
		if err != nil {
			return err
		}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(config.Env().SuperAdmin.Password), bcrypt.DefaultCost)

	superMechman := models.Account{
		DisplayName: "Super Admin!",
		Username:    config.Env().SuperAdmin.Username,
		Password:    string(hashedPassword),
		Type:        models.AccountTypeSuperAdmin,
		Permissions: models.AccountPermissionReadAccounts | models.AccountPermissionWriteAccounts |
			models.AccountPermissionReadPatient | models.AccountPermissionWritePatient |
			models.AccountPermissionReadMedicine | models.AccountPermissionWriteMedicine |
			models.AccountPermissionReadVirus | models.AccountPermissionWriteVirus |
			models.AccountPermissionReadBloodTest | models.AccountPermissionWriteBloodTest |
			models.AccountPermissionReadOtherVisits | models.AccountPermissionWriteOtherVisits,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = dbConn.Create(&superMechman)

	return nil
}
