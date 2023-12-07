package config

type Config struct {
	BindAddress  string `toml:"bind_address"`
	DatabaseUrl  string `toml:"database_url"`
	SessionKey   string `toml:"session_key"`
	SmtpEmail    string `toml:"smtp_email"`
	SmtpPassword string `toml:"smtp_password"`
}

func NewConfig() *Config {
	return &Config{
		BindAddress: ":8080",
	}
}
