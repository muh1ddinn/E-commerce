package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const ()

// Config ...
type Config struct {
	Environment   string // develop, staging, production
	RedisHost     string
	RedisPort     int
	RedisPassword string

	CatalogServiceHost string
	CatalogServicePort string

	LogLevel string
	HTTPPort string
}

// Load loads environment vars and inflates Config
func Load() Config {

	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "prod"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":1239"))

	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "127.0.0.1"))
	c.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))
	c.RedisPassword = cast.ToString(getOrReturnDefault("REDIS_PASSWORD", "3EEdwhDOfx"))

	c.CatalogServiceHost = cast.ToString(getOrReturnDefault("CATALOG_SERVICE_HOST", "localhost"))
	c.CatalogServicePort = cast.ToString(getOrReturnDefault("CATALOG_GRPC_PORT", "8082"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	if os.Getenv(key) == "" {
		return defaultValue
	}

	return os.Getenv(key)
}
