package v2

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const DEFAULT_ENV_FILE_NAME = ".env"

type Parsable interface {
	int | int64 | float32 | float64 | bool
}

// LoadEnv loads environment variables from the specified file.
func LoadEnv(filename string) error {
	if filename == "" {
		filename = DEFAULT_ENV_FILE_NAME
	}
	return godotenv.Load(filename)
}

// ------------- Default
func GetEnvString(key, defaultValue string) string {
	return getEnvGeneric(key, defaultValue, func(s string) (string, error) { return s, nil })
}

func GetEnvInt(key string, defaultValue int) int {
	return getEnvGeneric(key, defaultValue, strconv.Atoi)
}

func GetEnvBool(key string, defaultValue bool) bool {
	return getEnvGeneric(key, defaultValue, strconv.ParseBool)
}
func getEnvGeneric[T Parsable | string](key string, defaultValue T, parser func(string) (T, error)) T {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	val, err := parser(v)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s invalid: %v", key, err))
	}
	return val
}

// ------------- Must
func MustEnvString(key string) string {
	return mustEnvGeneric(key, func(s string) (string, error) { return s, nil })
}

func MustEnvInt(key string) int {
	return mustEnvGeneric(key, strconv.Atoi)
}

func MustEnvInt64(key string) int64 {
	return mustEnvGeneric(key, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	})
}

func MustEnvFloat32(key string) float32 {
	return mustEnvGeneric(key, func(s string) (float32, error) {
		val, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return 0, err
		}
		return float32(val), nil
	})
}

func MustEnvFloat64(key string) float64 {
	return mustEnvGeneric(key, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})
}

func MustEnvBool(key string) bool {
	return mustEnvGeneric(key, strconv.ParseBool)
}

func mustEnvGeneric[T Parsable | string](key string, parser func(string) (T, error)) T {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("environment variable %s is required", key))
	}
	val, err := parser(v)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s invalid: %v", key, err))
	}
	return val
}
