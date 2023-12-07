package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"relay-backend/internal/apiserver"
	"relay-backend/internal/apiserver/config"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "config file path")
}

func main() {
	flag.Parse()

	cfg := config.NewConfig()

	_, err := toml.DecodeFile(configPath, cfg)

	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(cfg); err != nil {
		log.Fatal(err)
	}
}
