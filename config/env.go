package config

import (
	"os"
	"shs/log"
)

var (
	_config = config{}
)

func initEnvVars() {
	_config = config{
		Port:      getEnv("PORT"),
		GoEnv:     GoEnv(getEnv("GO_ENV")),
		Hostname:  getEnv("HOST_NAME"),
		JwtSecret: getEnv("JWT_SECRET"),
		BlobsDir:  getEnv("BLOBS_DIR"),
		DB: struct {
			Name     string
			Host     string
			Username string
			Password string
		}{
			Name:     getEnv("DB_NAME"),
			Host:     getEnv("DB_HOST"),
			Username: getEnv("DB_USERNAME"),
			Password: getEnv("DB_PASSWORD"),
		},
		Cache: struct {
			Host     string
			Password string
		}{
			Host:     getEnv("CACHE_HOST"),
			Password: getEnv("CACHE_PASSWORD"),
		},
		SuperAdmin: struct {
			Username string
			Password string
		}{
			Username: getEnv("SUPERADMIN_USERNAME"),
			Password: getEnv("SUPERADMIN_PASSWORD"),
		},
	}
}

type GoEnv string

const (
	GoEnvProd GoEnv = "prod"
	GoEnvBeta GoEnv = "beta"
	GoEnvDev  GoEnv = "dev"
	GoEnvTest GoEnv = "test"
)

type config struct {
	Port      string
	GoEnv     GoEnv
	Hostname  string
	JwtSecret string
	BlobsDir  string
	DB        struct {
		Name     string
		Host     string
		Username string
		Password string
	}
	Cache struct {
		Host     string
		Password string
	}
	SuperAdmin struct {
		Username string
		Password string
	}
}

// Env returns the thing's config values :)
func Env() config {
	return _config
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalln("The \"" + key + "\" variable is missing.")
	}
	return value
}
