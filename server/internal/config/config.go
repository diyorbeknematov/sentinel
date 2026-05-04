package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	HTTPPort string

	DBHost     string
	DBName     string
	DBPort     string
	DBUser     string
	DBPassword string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string

	AccessToken  string
	RefreshToken string

	FrontendURL string

	MailerHost     string
	MailerPort     string
	MailerUsername string
	MailerPassword string
	MailerFrom     string

	KafkaBrokers         []string
	KafkaExternalBrokers []string
	KafkaTopic           string
	KafkaGroupID         string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		HTTPPort: getEnv("HTTP_PORT", "8081"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", ""),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", "rdpassword"),
		RedisDB:       getEnv("REDIS_DB", "0"),

		AccessToken:  getEnv("ACCESS_TOKEN", "key"),
		RefreshToken: getEnv("REFRESH_TOKEN", "key"),

		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),

		MailerHost:     getEnv("MAILER_HOST", "smtp.gmail.com"),
		MailerPort:     getEnv("MAILER_PORT", "587"),
		MailerUsername: getEnv("MAILER_USERNAME", ""),
		MailerPassword: getEnv("MAILER_PASSWORD", ""),
		MailerFrom:     getEnv("MAILER_FROM", "noreply@sentinel.com"),
		
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaExternalBrokers: strings.Split(getEnv("KAFKA_EXTERNAL_ROKERS", "localhost:9092"), ","),
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
