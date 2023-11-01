package apiserver

type Config struct {
	BindAddress string `toml:"bind_address"`
	DatabaseUrl string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddress: ":8080",
	}
}
