package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	HTTPPort     string
	DBHost       string
	DBName       string
	DBPort       string
	DBUser       string
	DBPassword   string
	CasbinDriver string
	AccessToken  string
	RefreshToken string

	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		HTTPPort:     getEnv("HTTP_PORT", "8081"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBName:       getEnv("DB_NAME", ""),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "password"),
		CasbinDriver: getEnv("DB_CASBIN_DRIVER", "postgres"),
		AccessToken:  getEnv("ACCESS_TOKEN", "key"),
		RefreshToken: getEnv("REFRESH_TOKEN", "key"),

		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "test-topic"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "test-group"),
	}
}

func getEnv(key string, defaultValue interface{}) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return cast.ToString(defaultValue)
	}
	return val
}
