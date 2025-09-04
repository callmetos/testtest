package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	App struct {
		Timezone string
	}

	Google struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{}
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.Database.DBName = getEnv("DB_NAME", "navmate")
	cfg.App.Timezone = getEnv("APP_TIMEZONE", "Asia/Bangkok")
	cfg.Google.ClientID = getEnv("GOOGLE_CLIENT_ID", "")
	cfg.Google.ClientSecret = getEnv("GOOGLE_CLIENT_SECRET", "")
	cfg.Google.RedirectURL = getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback")

	return cfg
}
