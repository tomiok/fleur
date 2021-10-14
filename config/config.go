package config

import "os"

type Config struct {
	Port string
}

func Bind() *Config {
	var port = os.Getenv("PORT")

	if port == "" {
		port = "5566"
	}

	return &Config{Port: port}
}
