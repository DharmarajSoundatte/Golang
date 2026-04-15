package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Server
	Port string
	Env  string

	// Database — either set DATABASE_URL for a full connection string,
	// or set the individual DB_* variables below.
	DatabaseURL string // e.g. postgres://user:pass@host:port/db?sslmode=require
	SSLRootCert string // path to CA cert file, e.g. ca.pem

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret      string
	JWTExpireHours int

	// CORS
	AllowedOrigins string
}

// Load reads the .env file (if present) and populates Config from env vars.
func Load() (*Config, error) {
	// .env is optional — in production env vars are injected by the platform.
	_ = godotenv.Load()

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		Env:            getEnv("APP_ENV", "development"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		SSLRootCert:    getEnv("SSL_ROOT_CERT", "ca.pem"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "golang_poc"),
		DBSSLMode:      getEnv("DB_SSL_MODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpireHours: 24,
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
	}

	if cfg.JWTSecret == "change-me-in-production" && cfg.Env == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return cfg, nil
}

// DSN returns the PostgreSQL keyword connection string for GORM.
// When DATABASE_URL is set it is parsed into keyword format so that special
// characters in the password are handled correctly without URL-encoding.
// sslrootcert is appended only when the cert file actually exists on disk.
func (c *Config) DSN() string {
	if c.DatabaseURL != "" {
		host, port, user, password, dbname, sslmode := parsePostgresURL(c.DatabaseURL)
		if host != "" {
			dsn := fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
				host, port, user, password, dbname, sslmode,
			)
			// Only attach the cert file when it actually exists on disk.
			if c.SSLRootCert != "" {
				if _, err := os.Stat(c.SSLRootCert); err == nil {
					dsn += " sslrootcert=" + c.SSLRootCert
				}
			}
			return dsn
		}
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// parsePostgresURL manually extracts connection components from a Postgres URL.
// It uses strings.LastIndex for "@" so passwords containing "@" are handled
// correctly, and avoids net/url which rejects unencoded special characters.
//
// Expected format: postgres://user:password@host:port/dbname?key=val&...
func parsePostgresURL(raw string) (host, port, user, password, dbname, sslmode string) {
	// Strip scheme.
	s := strings.TrimPrefix(raw, "postgres://")
	s = strings.TrimPrefix(s, "postgresql://")

	// Split userinfo from the rest using the LAST "@" so passwords with "@"
	// are handled correctly.
	atIdx := strings.LastIndex(s, "@")
	if atIdx < 0 {
		return
	}
	userinfo := s[:atIdx]
	hostAndRest := s[atIdx+1:]

	// Split user and password.
	if colonIdx := strings.Index(userinfo, ":"); colonIdx >= 0 {
		user = userinfo[:colonIdx]
		password = userinfo[colonIdx+1:]
	} else {
		user = userinfo
	}

	// Split host:port/dbname?params.
	slashIdx := strings.Index(hostAndRest, "/")
	if slashIdx < 0 {
		host = hostAndRest
		return
	}
	hostport := hostAndRest[:slashIdx]
	pathAndQuery := hostAndRest[slashIdx+1:]

	// Split host and port.
	if colonIdx := strings.LastIndex(hostport, ":"); colonIdx >= 0 {
		host = hostport[:colonIdx]
		port = hostport[colonIdx+1:]
	} else {
		host = hostport
	}

	// Split dbname and query params.
	if qIdx := strings.Index(pathAndQuery, "?"); qIdx >= 0 {
		dbname = pathAndQuery[:qIdx]
		for _, param := range strings.Split(pathAndQuery[qIdx+1:], "&") {
			kv := strings.SplitN(param, "=", 2)
			if len(kv) == 2 && kv[0] == "sslmode" {
				sslmode = kv[1]
			}
		}
	} else {
		dbname = pathAndQuery
	}

	if sslmode == "" {
		sslmode = "disable"
	}
	return
}

// IsDevelopment returns true when running in dev mode.
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// getEnv returns the env variable value or a fallback default.
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
