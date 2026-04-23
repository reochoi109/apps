package v1

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const DEFAULT_ENV_FILE_NAME = ".env"

func LoadEnv(filename string) error {
	if filename == "" {
		filename = DEFAULT_ENV_FILE_NAME
	}
	return godotenv.Load(filename)
}

func EnvString(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func EnvInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}

func MustEnvString(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("environment variable %s is required", key))
	}
	return v
}

func MustEnvInt(key string) int {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("environment variable %s is required", key))
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s must be an integer, got: %s", key, v))
	}
	return i
}

func MustEnvInt64(key string) int64 {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("environment variable %s is required", key))
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s must be an int64, got: %s", key, v))
	}
	return i
}

func MustEnvBool(key string) bool {
	v := strings.ToLower(os.Getenv(key))
	if v == "" {
		panic(fmt.Sprintf("environment variable %s is required", key))
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s must be a boolean, got: %s", key, v))
	}
	return b
}

func MustEnvFloat32(key string) float32 {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("environment variable %s is required", key))
	}
	i, err := strconv.ParseFloat(v, 32)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s must be a float32, got: %v", key, err))
	}
	return float32(i)
}

func MustEnvFloat64(key string) float64 {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("environment variable %s is required", key))
	}
	i, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s must be a float64, got: %v", key, err))
	}
	return i
}
