package config

type (
	CONFIG struct {
		MICRO MICRO
	}

	MICRO struct {
		DB               DB
		API              API
		TWILIO           TWILIO
		DEEPSEEK_API_KEY string
	}

	DB struct {
		PSQL PSQL
	}

	API struct {
		HOST string
		PORT string
	}

	PSQL struct {
		HOST   string
		PORT   string
		USER   string
		PASS   string
		SCHEMA string
	}

	TWILIO struct {
		ACCOUNT_SID     string
		AUTH_TOKEN      string
		WHATSAPP_NUMBER string
	}
)
