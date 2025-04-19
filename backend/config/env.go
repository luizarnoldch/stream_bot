package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	log "github.com/gofiber/fiber/v2/log"

	"github.com/joho/godotenv"
)

func sanityCheck() {
	requiredEnvVars := []string{
		"ENV",

		"API_HOST",
		"API_PORT",

		"PSQL_HOST",
		"PSQL_PORT",
		"PSQL_USER",
		"PSQL_PASS",
		"PSQL_SCHEMA",
		"PSQL_MAX_CONNS",

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

func (c *CONFIG) GetConnString() string {
	// Base connection string with common parameters
	base := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.MICRO.DB.PSQL.USER,
		c.MICRO.DB.PSQL.PASS,
		c.MICRO.DB.PSQL.HOST,
		c.MICRO.DB.PSQL.PORT,
		c.MICRO.DB.PSQL.SCHEMA,
	)

	// Parámetros adicionales
	params := url.Values{}

	if c.ENV == "dev" {
		params.Add("sslmode", "disable")
	} else {
		params.Add("sslmode", "require")
	}

	// Configuración del pool de conexiones
	params.Add("pool_max_conns", strconv.Itoa(c.MICRO.DB.PSQL.MAX_CONNS))

	return fmt.Sprintf("%s?%s", base, params.Encode())
}

func LoadConfig() (*CONFIG, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
		return nil, fmt.Errorf("error loading .env: %v", err)
	}

	sanityCheck()

	conns, _ := strconv.Atoi(os.Getenv("PSQL_MAX_CONNS"))

	return &CONFIG{
		MICRO: MICRO{
			API: API{
				HOST: os.Getenv("API_HOST"),
				PORT: os.Getenv("API_PORT"),
			},
			DB: DB{
				PSQL: PSQL{
					HOST:      os.Getenv("PSQL_HOST"),
					PORT:      os.Getenv("PSQL_PORT"),
					USER:      os.Getenv("PSQL_USER"),
					PASS:      os.Getenv("PSQL_PASS"),
					SCHEMA:    os.Getenv("PSQL_SCHEMA"),
					MAX_CONNS: conns,
				},
			},
			TWILIO: TWILIO{
				ACCOUNT_SID:     os.Getenv("TWILIO_ACCOUNT_SID"),
				AUTH_TOKEN:      os.Getenv("TWILIO_AUTH_TOKEN"),
				WHATSAPP_NUMBER: os.Getenv("TWILIO_WHATSAPP_NUMBER"),
			},
			DEEPSEEK_API_KEY: os.Getenv("DEEPSEEK_API_KEY"),
		},
		ENV: os.Getenv("ENV"),
	}, nil
}
