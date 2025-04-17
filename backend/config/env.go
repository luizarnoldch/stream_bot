package config

import (
	"fmt"
	"os"

	log "github.com/gofiber/fiber/v2/log"

	"github.com/joho/godotenv"
)

func sanityCheck() {
	requiredEnvVars := []string{
		"API_HOST",
		"API_PORT",

		"PSQL_HOST",
		"PSQL_PORT",
		"PSQL_USER",
		"PSQL_PASS",
		"PSQL_SCHEMA",

		"TWILIO_ACCOUNT_SID",
		"TWILIO_AUTH_TOKEN",
		"TWILIO_WHATSAPP_NUMBER",

		"DEEPSEEK_API_KEY",
	}

	for _, envVar := range requiredEnvVars {
		if value := os.Getenv(envVar); value == "" {
			log.Fatalf("Environment variable %s not defined. Terminating application...", envVar)
		}
	}
}

func LoadConfig() (*CONFIG, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
		return nil, fmt.Errorf("error loading .env: %v", err)
	}

	sanityCheck()

	return &CONFIG{
		MICRO: MICRO{
			API: API{
				HOST: os.Getenv("API_HOST"),
				PORT: os.Getenv("API_PORT"),
			},
			DB: DB{
				PSQL: PSQL{
					HOST:   os.Getenv("PSQL_HOST"),
					PORT:   os.Getenv("PSQL_PORT"),
					USER:   os.Getenv("PSQL_USER"),
					PASS:   os.Getenv("PSQL_PASS"),
					SCHEMA: os.Getenv("PSQL_SCHEMA"),
				},
			},
			TWILIO: TWILIO{
				ACCOUNT_SID:     os.Getenv("TWILIO_ACCOUNT_SID"),
				AUTH_TOKEN:      os.Getenv("TWILIO_AUTH_TOKEN"),
				WHATSAPP_NUMBER: os.Getenv("TWILIO_WHATSAPP_NUMBER"),
			},
			DEEPSEEK_API_KEY: os.Getenv("DEEPSEEK_API_KEY"),
		},
	}, nil
}
