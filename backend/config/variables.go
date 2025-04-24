package config

type (
	CONFIG struct {
		MICRO MICRO
		ENV   string
	}

	MICRO struct {
		DB               DB
		API              API
		TWILIO           TWILIO
		DEEPSEEK_API_KEY string
		OPENAI_API_KEY   string
	}

	DB struct {
		PSQL PSQL
	}

	API struct {
		HOST string
		PORT string
	}

	PSQL struct {
		HOST      string
		PORT      string
		USER      string
		PASS      string
		DB        string
		MAX_CONNS int
	}

	TWILIO struct {
		ACCOUNT_SID     string
		AUTH_TOKEN      string
		WHATSAPP_NUMBER string
		MESSAGING_SERVICE_SID string
	}
)
