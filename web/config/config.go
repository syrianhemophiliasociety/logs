package config

import (
	"os"
	"shs-web/log"
)

var (
	_config = config{}
)

func init() {
	_config = config{
		Port:          getEnv("WEB_PORT"),
		GoEnv:         getEnv("GO_ENV"),
		Hostname:      getEnv("HOST_NAME"),
		ServerAddress: getEnv("SERVER_ADDRESS"),
		Cache: struct {
			Host     string
			Password string
		}{
			Host:     getEnv("CACHE_HOST"),
			Password: getEnv("CACHE_PASSWORD"),
		},
	}
}

type config struct {
	Port          string
	GoEnv         string
	Hostname      string
	ServerAddress string
	Cache         struct {
		Host     string
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
		log.Fatalln(log.ErrorLevel, "The \""+key+"\" variable is missing.")
	}
	return value
}
