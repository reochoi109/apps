package v1

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, key, val string) {
	os.Setenv(key, val)
	t.Cleanup(func() { os.Unsetenv(key) })
}

func TestLoadEnv(t *testing.T) {
	filename := ".env.test"
	os.WriteFile(filename, []byte("TEST_NUMBER=117\nTEST_USER=John"), 0664)
	defer os.Remove(filename)

	if err := LoadEnv(filename); err != nil {
		t.Errorf("LoadEnv() failed : %v", err)
	}

	if os.Getenv("TEST_NUMBER") != "117" || os.Getenv("TEST_USER") != "John" {
		t.Error("environment variables not loaded correctly")
	}

	if err := LoadEnv("non_existent_file"); err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}
func TestEnvInt(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  string
		def  int
		want int
	}{
		{"Valid", "K1", "10", 0, 10},
		{"Missing", "K2", "", 50, 50},
		{"Invalid", "K3", "abc", 20, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val != "" {
				setEnv(t, tt.key, tt.val)
			}
			if got := EnvInt(tt.key, tt.def); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

// Must
func TestMustEnv(t *testing.T) {
	t.Run("ValidCases", func(t *testing.T) {
		setEnv(t, "K_INT", "10")
		setEnv(t, "K_BOOL", "true")

		if MustEnvInt("K_INT") != 10 {
			t.Error("MustEnvInt failed")
		}
		if MustEnvBool("K_BOOL") != true {
			t.Error("MustEnvBool failed")
		}
	})

	t.Run("PanicCases", func(t *testing.T) {
		assertPanic(t, func() { MustEnvInt("NON_EXISTENT") })

		setEnv(t, "K_BAD", "not-an-int")
		assertPanic(t, func() { MustEnvInt("K_BAD") })
	})

	t.Run("AdditionalTypes", func(t *testing.T) {
		setEnv(t, "K_STR", "hello")
		setEnv(t, "K_F32", "3.14")
		setEnv(t, "K_F64", "3.141592")

		if MustEnvString("K_STR") != "hello" {
			t.Error("MustEnvString failed")
		}

		if MustEnvFloat32("K_F32") != float32(3.14) {
			t.Error("MustEnvFloat32 failed")
		}
		if MustEnvFloat64("K_F64") != 3.141592 {
			t.Error("MustEnvFloat64 failed")
		}
	})
}

func assertPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic did not occur")
		}
	}()
	f()
}
