package config

import (
	"bytes"
	"errors"
	"testing"
)

func TestMode_IsValid(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
		want bool
	}{
		{name: "Development", mode: Development, want: true},
		{name: "Production", mode: Development, want: true},
		{name: "DefaultEnvPath", mode: Mode("Stage"), want: false},
		{name: "empty", mode: Mode(""), want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.mode.IsValid()
			if got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestParseMode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Mode
		wantErr bool
	}{
		{"empty means default", "", Development, false},
		{"dev", "dev", Development, false},
		{"prod", "prod", Production, false},
		{"invalid", "local", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseMode(tc.input)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if tc.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}

		})
	}
}

func TestOption_ApplyDefaults(t *testing.T) {
	t.Run("set default env when empty", func(t *testing.T) {
		opt := Option{
			Mode: Development,
			Env:  "",
		}
		opt.ApplyDefaults()
	})

	t.Run("Keep exists env", func(t *testing.T) {
		opt := Option{
			Mode: Development,
			Env:  "/custom/path/.env",
		}

		opt.ApplyDefaults()

		if opt.Env != "/custom/path/.env" {
			t.Fatalf("expected env to stay unchanged, got %v", opt.Env)
		}
	})
}

func TestOptionValidate(t *testing.T) {
	tests := []struct {
		name    string
		opt     Option
		wantErr bool
	}{
		{"valid dev", Option{Mode: Development, Env: DefaultEnvPath}, false},
		{"valid prod", Option{Mode: Production, Env: DefaultEnvPath}, false},
		{"invalid mode", Option{Mode: Mode("local"), Env: DefaultEnvPath}, true},
		{"empty mode", Option{Mode: "", Env: DefaultEnvPath}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opt.Validate()

			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if tc.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			if tc.wantErr && !errors.Is(err, ErrInvalidMode) {
				t.Fatalf("expected ErrInvalidMode, got %v", err)
			}
		})
	}
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		want       Option
		wantErr    bool
		checkErrIs error
	}{
		{
			name: "defailt",
			args: []string{},
			want: Option{
				Mode: Development,
				Env:  DefaultEnvPath,
			},
			wantErr: false,
		},

		{
			name: "valid dev with env",
			args: []string{"-m", "dev", "-e", "/tmp/dev.env"},
			want: Option{
				Mode: Development,
				Env:  "/tmp/dev.env",
			},
			wantErr: false,
		},

		{
			name: "valid prod with env",
			args: []string{"-m", "prod", "-e", "/tmp/prod.env"},
			want: Option{
				Mode: Production,
				Env:  "/tmp/prod.env",
			},
			wantErr: false,
		},
		{
			name:       "invalid mode",
			args:       []string{"-m", "local"},
			wantErr:    true,
			checkErrIs: nil, // ParseMode는 ErrInvalidMode로 감싸지 않아서 여기선 문자열 검사만 해도 됨
		},

		{
			name:       "positional args not allowed",
			args:       []string{"hello"},
			wantErr:    true,
			checkErrIs: nil,
		},

		{
			name:    "unknown flag",
			args:    []string{"-x"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			got, err := ParseFlags(buf, tc.args)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			if !tc.wantErr && got != tc.want {
				t.Fatalf("expected %+v, got %+v", tc.want, got)
			}

			if tc.name == "invalid mode" && err != nil {
				if err.Error() != `invalid mode: "local"` {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if tc.name == "positional args not allowed" && err != nil {
				if err.Error() != "positional args not allowed" {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
