package logger

import (
	"io"
	"os"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

type Config struct {
	Service      string
	Format       Format
	Level        string
	Output       io.Writer
	ReportCaller bool
}

// product mode
func PresetProd(service string) Config {
	return Config{
		Service:      service,
		Format:       FormatJSON,
		Level:        "info",
		Output:       os.Stdout,
		ReportCaller: false,
	}
}

// dev mode
func PresetDev(service string) Config {
	return Config{
		Service:      service,
		Format:       FormatText,
		Level:        "debug",
		Output:       os.Stdout,
		ReportCaller: true,
	}
}
