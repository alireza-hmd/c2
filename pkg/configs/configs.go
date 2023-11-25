package configs

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}
	return env
}

var Configs = make(map[string]string)

func Get(key string) string {
	return Configs[key]
}

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	Configs["DB_USER"] = GetEnv("DB_USER", "test")
	Configs["DB_PASS"] = GetEnv("DB_PASS", "test")
	Configs["DB_HOST"] = GetEnv("DB_HOST", "127.0.0.1")
	Configs["DB_PORT"] = GetEnv("DB_PORT", "3306")
	Configs["DB_DATABASE"] = GetEnv("DB_DATABASE", "test")
	return nil
}
