package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/youssefM1999/report/pkg/filesystem"
)

func LoadEnvFile() error {
	rootDir, err := filesystem.FindModuleRoot()
	if err != nil {
		return fmt.Errorf("failed to find module root: %w", err)
	}
	envFile := filepath.Join(rootDir, ".env")
	if err := filesystem.CheckValidFile(envFile); err != nil {
		return fmt.Errorf("failed to check valid env file: %w", err)
	}
	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("failed to load env file: %w", err)
	}
	return nil
}

func GetString(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	duration, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return duration
}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}
