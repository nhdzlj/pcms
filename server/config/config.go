package config

import (
	"os"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
	UploadDir  string
}

var AppConfig *Config

func Load() *Config {
	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "pcms"),
		DBPassword: getEnv("DB_PASSWORD", "pcms123"),
		DBName:     getEnv("DB_NAME", "pcms"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "pcms-jwt-secret-key-2024"),
		UploadDir:  getEnv("UPLOAD_DIR", "./uploads"),
	}

	// 确保上传目录存在
	os.MkdirAll(cfg.UploadDir, 0755)

	AppConfig = cfg
	return cfg
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode +
		" TimeZone=Asia/Shanghai"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
