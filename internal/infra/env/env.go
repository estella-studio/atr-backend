package env

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Env struct {
	LimiterMax                     int    `env:"LIMITER_MAX"`
	LimiterExpirationMinutes       int    `env:"LIMITER_EXPIRATION_MINUTES"`
	PasswordChangeExpiryMinutes    int    `env:"PASSWORD_CHANGE_EXPIRY_MINUTES"`
	PasswordChangeCodeRetrySeconds int    `env:"PASSWORD_CHANGE_RETRY_SECONDS"`
	AppPort                        uint   `env:"APP_PORT"`
	DBName                         string `env:"DB_NAME"`
	DBUsername                     string `env:"DB_USERNAME"`
	DBPassword                     string `env:"DB_PASSWORD"`
	DBHost                         string `env:"DB_HOST"`
	DBPort                         uint   `env:"DB_PORT"`
	JWTSecretKey                   string `env:"JWT_SECRET_KEY"`
	JWTExpiredDays                 uint   `env:"JWT_EXPIRED_DAYS"`
	EmailFrom                      string `env:"EMAIL_FROM"`
	SMTPServer                     string `env:"SMTP_SERVER"`
	SMTPPort                       int    `env:"SMTP_PORT"`
	SMTPUsername                   string `env:"SMTP_USERNAME"`
	SMTPPassword                   string `env:"SMTP_PASSWORD"`
	SMTPFrom                       string `env:"SMTP_FROM"`
	MailtrapToken                  string `env:"MAILTRAP_TOKEN"`
	MailtrapURL                    string `env:"MAILTRAP_URL"`
	MailtrapTemplate               string `env:"MAILTRAP_TEMPLATE"`
	MailtrapCompanyInfoName        string `env:"MAILTRAP_COMPANY_INFO_NAME"`
	MailtrapCompanyInfoAddress     string `env:"MAILTRAP_COMPANY_INFO_ADDRESS"`
	MailtrapCompanyInfoCity        string `env:"MAILTRAP_COMPANY_INFO_CITY"`
	MailtrapCompanyInfoZipCode     string `env:"MAILTRAP_COMPANY_INFO_ZIP_CODE"`
	MailtrapCompanyInfoCountry     string `env:"MAILTRAP_COMPANY_INFO_COUNTRY"`
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
