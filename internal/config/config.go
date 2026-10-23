package config

import (
	"errors"
	"os"
)

type Config struct {
	Port          string
	DBPath        string
	UploadsDir    string
	PublicURL     string
	JWTSecret     string
	AdminEmail    string
	AdminPassword string
	SeedDemo      bool
	DemoPassword  string
}

func Load() (Config, error) {
	cfg := Config{
		Port:          getenv("PORT", "3000"),
		DBPath:        getenv("DB_PATH", "database/database.db"),
		UploadsDir:    getenv("UPLOADS_DIR", "uploads"),
		PublicURL:     getenv("PUBLIC_URL", "http://localhost:3000"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		SeedDemo:      os.Getenv("SEED_DEMO") == "true",
		DemoPassword:  os.Getenv("DEMO_PASSWORD"),
	}

	if cfg.JWTSecret == "" {
		return cfg, errors.New("JWT_SECRET is not set")
	}
	if cfg.SeedDemo && len(cfg.DemoPassword) < 8 {
		return cfg, errors.New("DEMO_PASSWORD must be at least 8 characters when SEED_DEMO=true")
	}

	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
