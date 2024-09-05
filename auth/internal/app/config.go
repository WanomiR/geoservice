package app

type (
	Config struct {
		Host        string `env-required:"true" env:"HOST"`
		Port        string `env-required:"true" env:"PORT"`
		MetricsPort string `env-required:"true" env:"METRICS_PORT"`
		AppVersion  string `env-required:"true" env:"APP_VERSION"`

		Log
	}

	Log struct {
		Level string `env:"LOG_LEVEL" envDefault:"debug"`
	}
)
