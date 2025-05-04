package config

import (
	"future_today/internal/cerrors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbHost         string
	DbPort         string
	DbUser         string
	DbPass         string
	DbName         string
	AgifyURL       string
	GenderizeURL   string
	NationalizeURL string
	ServerPort     string
}

func GetConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, cerrors.ErrLoadEnv
	}
	return &Config{
		ServerPort:     os.Getenv("SERVER_PORT"),
		DbHost:         os.Getenv("DB_HOST"),
		DbPort:         os.Getenv("DB_PORT"),
		DbUser:         os.Getenv("DB_USER"),
		DbPass:         os.Getenv("DB_PASS"),
		DbName:         os.Getenv("DB_NAME"),
		AgifyURL:       os.Getenv("API_AGIFY_URL"),
		GenderizeURL:   os.Getenv("API_GENDERIZE_URL"),
		NationalizeURL: os.Getenv("API_NATIONALIZE_URL"),
	}, nil

}
