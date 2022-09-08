package config

var AppConfig config

type config struct {
	Database struct {
		Dsn     string
		EchoSQL bool
	}
	Slack struct {
		WebhookURL string
	}
	Debug bool
}
