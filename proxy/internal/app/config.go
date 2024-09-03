package app

type (
	Config struct {
		Host       string `env-required:"true" env:"HOST"`
		Port       string `env-required:"true" env:"PORT"`
		AppVersion string `env-required:"true" env:"APP_VERSION"`

		Log
		Geo
	}

	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	Geo struct {
		Host string `env-required:"true" env:"GEO_HOST"`
		Port string `env-required:"true" env:"GEO_PORT"`
	}

	Auth struct {
		Host       string `env-required:"true" env:"AUTH_HOST"`
		Port       string `env-required:"true" env:"AUTH_PORT"`
		CookieName string `env-required:"true" env:"AUTH_COOKIE_NAME"`
	}

	//User struct {
	//	Host string `env-required:"true" env:"USER_HOST"`
	//	Port string `env-required:"true" env:"USER_PORT"`
	//}
)
