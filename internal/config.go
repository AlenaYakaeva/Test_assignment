package internal

import (
	"cmp"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host  string
	Port  string
	DBDSN string
}

func ReadConfig() *Config {
	if err := godotenv.Load("config.env"); err != nil {
		//log.Printf("Предупреждение: файл config.env не найден, используются дефолтные значения или системный env: %w", err)
		panic(err)
	}

	var cfg Config

	cfg.Host = cmp.Or(os.Getenv("EXAMPLER_SERVICE_HOST"), "0.0.0.0")
	cfg.Port = cmp.Or(os.Getenv("EXAMPLER_SERVICE_PORT"), "8080")
	cfg.DBDSN = cmp.Or(os.Getenv("EXAMPLER_SERVICE_DB_DSN"), "postgres://postgres:postgresql_pass@localhost:5432/example_app?sslmode=disable")

	return &cfg
}
