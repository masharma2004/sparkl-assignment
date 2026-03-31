package config

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                         string
	DBHost                       string
	DBPort                       string
	DBUser                       string
	DBPass                       string
	DBName                       string
	DBSSLMode                    string
	JWTSecret                    string
	JWTIssuer                    string
	AccessTokenTTLMinutes        int
	RefreshTokenTTLDays          int
	MaintenanceIntervalMinutes   int
	RevokedRefreshRetentionHours int
	AllowedOrigins               []string
	AccessCookieName             string
	RefreshCookieName            string
	CookieDomain                 string
	CookieSecure                 bool
}

type AuthCookieConfig struct {
	AccessCookieName     string
	RefreshCookieName    string
	Domain               string
	Secure               bool
	AccessMaxAgeSeconds  int
	RefreshMaxAgeSeconds int
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func requireEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s is required but not set", key)
}

func parseAllowedOrigins(raw string) []string {
	values := strings.Split(raw, ",")
	origins := make([]string, 0, len(values))

	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" || slices.Contains(origins, trimmed) {
			continue
		}

		origins = append(origins, trimmed)
	}

	return origins
}

func parseBoolEnv(key string, fallback bool) bool {
	rawValue, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(rawValue) == "" {
		return fallback
	}

	value := strings.ToLower(strings.TrimSpace(rawValue))
	return value == "1" || value == "true" || value == "yes"
}

func parseIntEnv(key string, fallback int) int {
	rawValue, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(rawValue) == "" {
		return fallback
	}

	parsedValue, err := strconv.Atoi(strings.TrimSpace(rawValue))
	if err != nil || parsedValue <= 0 {
		return fallback
	}

	return parsedValue
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	required := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "JWT_SECRET"}
	values := make(map[string]string, len(required))

	for _, key := range required {
		value, err := requireEnv(key)
		if err != nil {
			return nil, err
		}
		values[key] = value
	}

	return &Config{
		Port:                         getEnv("PORT", "8080"),
		DBHost:                       values["DB_HOST"],
		DBPort:                       values["DB_PORT"],
		DBUser:                       values["DB_USER"],
		DBPass:                       values["DB_PASSWORD"],
		DBName:                       values["DB_NAME"],
		DBSSLMode:                    getEnv("DB_SSLMODE", "disable"),
		JWTSecret:                    values["JWT_SECRET"],
		JWTIssuer:                    getEnv("JWT_ISSUER", "sparkl-edventure"),
		AccessTokenTTLMinutes:        parseIntEnv("ACCESS_TOKEN_TTL_MINUTES", 15),
		RefreshTokenTTLDays:          parseIntEnv("REFRESH_TOKEN_TTL_DAYS", 7),
		MaintenanceIntervalMinutes:   parseIntEnv("MAINTENANCE_INTERVAL_MINUTES", 30),
		RevokedRefreshRetentionHours: parseIntEnv("REVOKED_REFRESH_RETENTION_HOURS", 24),
		AllowedOrigins:               parseAllowedOrigins(getEnv("FRONTEND_ORIGINS", "http://localhost:5173")),
		AccessCookieName:             getEnv("ACCESS_COOKIE_NAME", getEnv("SESSION_COOKIE_NAME", "sparkl_access")),
		RefreshCookieName:            getEnv("REFRESH_COOKIE_NAME", "sparkl_refresh"),
		CookieDomain:                 strings.TrimSpace(getEnv("COOKIE_DOMAIN", getEnv("SESSION_COOKIE_DOMAIN", ""))),
		CookieSecure:                 parseBoolEnv("COOKIE_SECURE", parseBoolEnv("SESSION_COOKIE_SECURE", false)),
	}, nil
}

func (cfg *Config) AuthCookieConfig() AuthCookieConfig {
	return AuthCookieConfig{
		AccessCookieName:     cfg.AccessCookieName,
		RefreshCookieName:    cfg.RefreshCookieName,
		Domain:               cfg.CookieDomain,
		Secure:               cfg.CookieSecure,
		AccessMaxAgeSeconds:  cfg.AccessTokenTTLMinutes * 60,
		RefreshMaxAgeSeconds: cfg.RefreshTokenTTLDays * 24 * 60 * 60,
	}
}
