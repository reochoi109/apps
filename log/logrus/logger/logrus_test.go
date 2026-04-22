package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	t.Run("Success Configuration", func(t *testing.T) {
		cfg := Config{
			Service:      "test-service",
			Format:       FormatJSON,
			Level:        "debug",
			ReportCaller: true,
		}

		log := new(cfg)
		if log.GetLevel() != logrus.DebugLevel {
			t.Errorf("Level mismatch: got %v, want %v", log.GetLevel(), logrus.DebugLevel)
		}
	})

	t.Run("Invalid Level Fallback", func(t *testing.T) {
		cfg := Config{Level: "invalid-level"}
		log := new(cfg)
		if log.GetLevel() != logrus.InfoLevel {
			t.Errorf("Expected fallback to InfoLevel, got %v", log.GetLevel())
		}
	})
}
func TestGetOutput(t *testing.T) {
	var buf bytes.Buffer
	cfg := PresetDev("test")
	cfg.Output = io.MultiWriter(&buf, os.Stdout)

	log := new(cfg)

	log.Info("hello info log")
	if buf.Len() == 0 {
		t.Error("not print log")
	}
}
func TestCustomCallerPrettyfier(t *testing.T) {
	var buf bytes.Buffer
	cfg := PresetDev("test")
	cfg.ReportCaller = true
	cfg.Output = &buf

	log := new(cfg)
	log.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "TestCustomCallerPrettyfier()") {
		t.Errorf("invalid print to Caller. Got: %s", output)
	}
}
func TestSetFormatter(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		want   string
	}{
		{"JSON Format", FormatJSON, "*logrus.JSONFormatter"},
		{"Text Format", FormatText, "*logrus.TextFormatter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := PresetDev("test")
			cfg.Format = tt.format
			l := new(cfg)

			got := fmt.Sprintf("%T", l.Formatter)
			if got != tt.want {
				t.Errorf("Formatter mismatch: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceHook(t *testing.T) {
	serviceName := "prod-service"
	var buf bytes.Buffer
	cfg := PresetProd(serviceName)
	cfg.Output = &buf

	log := new(cfg)
	log.Info("hook test")

	output := buf.String()
	if !strings.Contains(output, serviceName) {
		t.Errorf("Service name missing in log. Expected to contain: %s, Got: %s", serviceName, output)
	}
}
