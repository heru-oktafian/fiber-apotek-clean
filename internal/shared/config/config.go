package config

import "os"

type Config struct {
	AppName      string
	ServerPort   string
	JWTSecret    string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPass       string
	DBName       string
	RedisHost    string
	RedisPort    string
	RedisPass    string
	RedisDB      string
	TimezoneName string
}

func Load() Config {
	return Config{
		AppName:      os.Getenv("APPNAME"),
		ServerPort:   envOr("SERVER_PORT", "1112"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPass:       os.Getenv("DB_PASS"),
		DBName:       os.Getenv("DB_NAME"),
		RedisHost:    os.Getenv("REDIS_HOST"),
		RedisPort:    os.Getenv("REDIS_PORT"),
		RedisPass:    os.Getenv("REDIS_PASS"),
		RedisDB:      envOr("REDIS_SHORT", "0"),
		TimezoneName: envOr("APP_TIMEZONE", "Asia/Jakarta"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
