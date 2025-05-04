package env

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Env struct {
	AppPort        uint   `env:"APP_PORT"`
	DBName         string `env:"DB_NAME"`
	DBUsername     string `env:"DB_USERNAME"`
	DBPassword     string `env:"DB_PASSWORD"`
	DBHost         string `env:"DB_HOST"`
	DBPort         uint   `env:"DB_PORT"`
	JWTSecretKey   string `env:"JWT_SECRET_KEY"`
	JWTExpiredDays uint   `env:"JWT_EXPIRED_DAYS"`
	GoogleClientID string `env:"GOOGLE_CLIENT_ID"`
}

func New() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	envParsed := new(Env)
	err = env.Parse(envParsed)
	if err != nil {
		return nil, err
	}

	return envParsed, nil
}
