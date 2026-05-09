package mariadb

import (
	"shs/app/models"
	"shs/config"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/schema"
)

var migratableModels = []schema.Tabler{
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
	new(models.Diagnosis),
	new(models.DiagnosisResult),
}

func Migrate() error {
	dbConn, err := dbConnector()
	if err != nil {
		return err
	}

	for _, table := range migratableModels {
		err = dbConn.Debug().AutoMigrate(table)
		if err != nil {
			return err
		}
	}

	for _, table := range migratableModels {
		err = dbConn.Exec("ALTER TABLE " + table.TableName() + " CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error
		if err != nil {
			return err
		}
	}

	_ = (&Repository{dbConn}).CreateSuperAdmin()

	return nil
}

func (r *Repository) CreateSuperAdmin() error {
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

	return r.client.Create(&superMechman).Error
}

func (r *Repository) DeleteAll() error {
	err := r.client.Exec("SET FOREIGN_KEY_CHECKS=0;").Error
	if err != nil {
		return err
	}

	for _, table := range migratableModels {
		var err error
		_, ok := table.(models.Account)
		if ok {
			err = r.client.Model(table).Where("username != ?", "b").Delete(nil).Error
		} else {
			err = r.client.Model(table).Where("true").Delete(nil).Error
		}
		if err != nil {
			return err
		}
	}

	return nil
}
