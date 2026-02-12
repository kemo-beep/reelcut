package config

import (
	"github.com/joho/godotenv"
)

func init() {
	godotenvLoad = func() error { return godotenv.Load() }
}
