package config

type Config struct {
	Port        string
	DatabaseURL string
	APIKey      string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv,
	}
}
