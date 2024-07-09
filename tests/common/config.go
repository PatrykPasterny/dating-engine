package common

import (
	"os"
)

type Config struct {
	DatabaseURI        string
	DatabaseName       string
	DatabaseCollection string
	BaseURL            string
}

func NewConfig() *Config {
	return &Config{
		DatabaseURI:        os.Getenv("DATABASE_URI"),
		DatabaseName:       os.Getenv("DATABASE_NAME"),
		DatabaseCollection: os.Getenv("DATABASE_COLLECTION"),
		BaseURL:            os.Getenv("BASE_URL"),
	}
}
