package db

import (
	"fmt"
	"net/url"
	"os"
)

// DSNFromEnv builds a Postgres connection string from the environment.
//
// If DATABASE_URL is set it is used verbatim. Otherwise the DSN is composed
// from the POSTGRES_* variables, each falling back to a dev-friendly default:
//
//	POSTGRES_HOST     (localhost)
//	POSTGRES_PORT     (5432)
//	POSTGRES_USER     (haves)
//	POSTGRES_PASSWORD (haves)
//	POSTGRES_DB       (haves)
//	POSTGRES_SSLMODE  (disable)
func DSNFromEnv() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	u := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			env("POSTGRES_USER", "haves"),
			env("POSTGRES_PASSWORD", "haves"),
		),
		Host: fmt.Sprintf("%s:%s",
			env("POSTGRES_HOST", "localhost"),
			env("POSTGRES_PORT", "5432"),
		),
		Path:     env("POSTGRES_DB", "haves"),
		RawQuery: "sslmode=" + env("POSTGRES_SSLMODE", "disable"),
	}
	return u.String()
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
