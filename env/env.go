package env

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

func (key EnvKey) Get() string {
	return os.Getenv(string(key))
}

func LoadEnv() {
	loclog := "[env.LoadEnv]"
	err := godotenv.Load()
	if err != nil {
		slog.Error(loclog, "FATAL", "error loading .env file", "error", err.Error())
		os.Exit(1)
	}
	slog.Info(loclog, "info", ".env loaded")

	// log
	vars := []EnvKey{
		DevHost,
		DevPort,
		Host,
		Port,
		DevMode,
		TLSCertPath,
		TLSKeyPath,
		AllowedOrigins,
		AllowedMethods,
		DBPath,
		HealthCheckToken,
		RateLimit,
		RateBurst,
	}
	for _, v := range vars {
		slog.Info(loclog, "info", "envvar", "key", v, "value", v.Get())
	}
}

const (
	DevHost          EnvKey = "DEV_HOST"
	DevPort          EnvKey = "DEV_PORT"
	Host             EnvKey = "HOST"
	Port             EnvKey = "PORT"
	DevMode          EnvKey = "DEV"
	TLSCertPath      EnvKey = "TLS_CERT_PATH"
	TLSKeyPath       EnvKey = "TLS_KEY_PATH"
	AllowedOrigins   EnvKey = "ALLOWED_ORIGINS"
	AllowedMethods   EnvKey = "ALLOWED_METHODS"
	DBPath           EnvKey = "DB_PATH"
	HealthCheckToken EnvKey = "HC_TOKEN"
	RateLimit        EnvKey = "RATE_LIMIT"
	RateBurst        EnvKey = "RATE_BURST"
)
