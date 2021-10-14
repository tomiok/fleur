package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
)

var c *Config

type Config struct {
	Port string `json:"port"`
}

func Bind() *Config {
	if c != nil {
		return c
	}

	var port = os.Getenv("PORT")

	if port == "" {
		port = "5566"
	}

	c = &Config{Port: port}
	return c
}

func (conf *Config) PrintConfigs() {
	if conf == nil {
		log.Print("configs are nil")
		return
	}
	b, err := json.MarshalIndent(conf, "", "    ")

	if err != nil {
		log.Printf("cannot marshal configs %s", err.Error())
	}

	log.Print(string(b))
}
