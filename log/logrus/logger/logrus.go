package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// create logrus and settings
func Set(cfg Config) {
	_ = new(cfg)
}

func new(cfg Config) *logrus.Logger {
	l := logrus.New()

	l.SetOutput(getOutput(cfg.Output))
	lv, err := logrus.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		logrus.WithError(err).Errorf("invalid log level: %q", cfg.Level)
		lv = logrus.InfoLevel
	}
	l.SetLevel(lv)

	if cfg.Service != "" {
		l.AddHook(&serviceHook{service: cfg.Service})
	}
	l.SetReportCaller(cfg.ReportCaller)

	if err := setFormatter(l, cfg.Format); err != nil {
		logrus.WithError(err).Warnf("failed to set formatter, falling back to JSON")
		_ = setFormatter(l, FormatJSON)
	}
	return l
}

func getOutput(out io.Writer) io.Writer {
	if out == nil {
		return os.Stdout
	}
	return out
}

func customCallerPrettyfier(f *runtime.Frame) (string, string) {
	parts := strings.Split(f.Function, ".")
	fn := parts[len(parts)-1]
	return fn + "()", fmt.Sprintf("%s:%d", f.File, f.Line)
}

func setFormatter(l *logrus.Logger, format Format) error {
	switch format {
	case FormatJSON:
		l.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: customCallerPrettyfier,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "severity",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
			TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
		})
	case FormatText:
		l.SetFormatter(&logrus.TextFormatter{
			ForceColors:      true,
			FullTimestamp:    true,
			DisableColors:    false,
			CallerPrettyfier: customCallerPrettyfier,
		})
	default:
		return fmt.Errorf("invalid format type %q", format)
	}
	return nil
}

// logger service hook
type serviceHook struct{ service string }

func (sh *serviceHook) Levels() []logrus.Level { return logrus.AllLevels }
func (sh *serviceHook) Fire(e *logrus.Entry) error {
	if _, ok := e.Data["service"]; !ok {
		e.Data["service"] = sh.service
	}
	return nil
}
