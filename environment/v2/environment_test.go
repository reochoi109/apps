package v2

import (
	"os"
	"strconv"
	"testing"
)

func setEnv(t *testing.T, key, val string) {
	t.Helper()
	os.Setenv(key, val)
	t.Cleanup(func() { os.Unsetenv(key) })
}

func assertPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("expected panic, but did not panic")
		}
	}()
	f()
}

func TestGetEnv(t *testing.T) {
	t.Run("DefaultValueWhenMissing", func(t *testing.T) {
		val := GetEnvInt("MISSING_KEY", 100)
		if val != 100 {
			t.Errorf("expected 100, got %d", val)
		}
	})

	t.Run("ValueWhenExists", func(t *testing.T) {
		setEnv(t, "PORT", "8080")
		val := GetEnvInt("PORT", 3000)
		if val != 8080 {
			t.Errorf("expected 8080, got %d", val)
		}
	})
}
func TestMustEnv(t *testing.T) {
	// 1. success
	t.Run("SuccessCases", func(t *testing.T) {
		setEnv(t, "K_INT", "10")
		setEnv(t, "K_F32", "3.14")
		setEnv(t, "K_BOOL", "true")

		if MustEnvInt("K_INT") != 10 {
			t.Error("MustEnvInt mismatch")
		}
		if MustEnvFloat32("K_F32") != 3.14 {
			t.Error("MustEnvFloat32 mismatch")
		}
		if MustEnvBool("K_BOOL") != true {
			t.Error("MustEnvBool mismatch")
		}
	})

	// 2. painc
	t.Run("PanicCases", func(t *testing.T) {
		// 필수 값이 없을 때
		assertPanic(t, func() { MustEnvInt("REQUIRED_KEY") })

		// 타입이 잘못되었을 때
		setEnv(t, "K_BAD", "not-a-number")
		assertPanic(t, func() { MustEnvInt("K_BAD") })
	})
}

func TestLoadEnv(t *testing.T) {
	filename := ".env.test"
	content := "TEST_NUMBER=117\nTEST_USER=John"
	os.WriteFile(filename, []byte(content), 0664)
	t.Cleanup(func() { os.Remove(filename) })

	if err := LoadEnv(filename); err != nil {
		t.Errorf("LoadEnv() failed: %v", err)
	}

	if os.Getenv("TEST_NUMBER") != "117" || os.Getenv("TEST_USER") != "John" {
		t.Error("environment variables not loaded correctly")
	}
}

// Default
func TestGetEnvGeneric(t *testing.T) {
	os.Unsetenv("TEST_INT")
	os.Unsetenv("TEST_BOOL")

	t.Run("ReturnDefaultValueWhenMissing", func(t *testing.T) {
		val := getEnvGeneric("NON_EXISTENT", 100, strconv.Atoi)
		if val != 100 {
			t.Errorf("expected 100, got %d", val)
		}
	})

	t.Run("ReturnParsedValueWhenEnvExists", func(t *testing.T) {
		os.Setenv("TEST_INT", "200")
		defer os.Unsetenv("TEST_INT")

		val := getEnvGeneric("TEST_INT", 100, strconv.Atoi)
		if val != 200 {
			t.Errorf("expected 200, got %d", val)
		}
	})

	t.Run("PanicsWhenParsingFails", func(t *testing.T) {
		os.Setenv("TEST_BOOL", "not-a-bool")
		defer os.Unsetenv("TEST_BOOL")

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic, but did not panic")
			}
		}()

		getEnvGeneric("TEST_BOOL", false, strconv.ParseBool)
	})
}
func TestGetEnvString(t *testing.T) {
	os.Setenv("HOST", "127.0.0.1")
	defer os.Unsetenv("HOST")

	val := GetEnvString("HOST", "localhost")
	if val != "127.0.0.1" {
		t.Errorf("expected 127.0.0.1, got %s", val)
	}

	val = GetEnvString("EMPTY", "localhost")
	if val != "localhost" {
		t.Errorf("expected localhost, got %s", val)
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	val := GetEnvInt("PORT", 3000)
	if val != 8080 {
		t.Errorf("expected 8080, got %d", val)
	}

	val = GetEnvInt("EMPTY", 3000)
	if val != 3000 {
		t.Errorf("expected 3000, got %d", val)
	}
}
func TestGetEnvBool(t *testing.T) {
	os.Setenv("Debug", "false")
	defer os.Unsetenv("Debug")

	val := GetEnvBool("Debug", true)
	if val {
		t.Errorf("expected false, got %v", val)
	}

	val = GetEnvBool("EMPTY", true)
	if !val {
		t.Errorf("expected false, got %v", val)
	}
}

// Must
func TestMustEnvParsers(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  string
		fn   func(string) any
		want any
	}{
		{"Int", "K1", "10", func(k string) any { return MustEnvInt(k) }, 10},
		{"Int64", "K2", "30", func(k string) any { return MustEnvInt64(k) }, int64(30)},
		{"Float32", "K3", "3.14", func(k string) any { return MustEnvFloat32(k) }, float32(3.14)},
		{"Float64", "K4", "3.141592", func(k string) any { return MustEnvFloat64(k) }, float64(3.141592)},
		{"bool", "K5", "true", func(k string) any { return MustEnvBool(k) }, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.key, test.val)
			defer os.Unsetenv(test.key)
			if got := test.fn(test.key); got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
