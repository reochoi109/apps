package logger

import (
	"os"
	"testing"
)

func TestCreateConfig(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		expect Config
	}{
		{
			name:   "sample_1",
			config: PresetProd("sample_service_name_1"),
			expect: Config{
				Service:      "sample_service_name_1",
				Format:       FormatJSON,
				Level:        "info",
				Output:       os.Stdout,
				ReportCaller: false,
			},
		},
		{
			name:   "sample_2",
			config: PresetDev("sample_service_name_2"),
			expect: Config{
				Service:      "sample_service_name_2",
				Format:       FormatText,
				Level:        "debug",
				Output:       os.Stdout,
				ReportCaller: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.config
			if c.Service != tc.expect.Service ||
				c.Format != tc.expect.Format ||
				c.Level != tc.expect.Level ||
				c.ReportCaller != tc.expect.ReportCaller {

				t.Errorf("Config mismatch!\nGot: %+v\nWant: %+v", c, tc.expect)
			}
		})
	}
}
