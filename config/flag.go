package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

type Mode string

const (
	Development    Mode = "dev"
	Production     Mode = "prod"
	DefaultEnvPath      = "service/default/path/env_name"
)

var (
	ErrPositionalArgs = errors.New("positional args not allowed")
	ErrInvalidMode    = errors.New("invalid mode")
)

type Option struct {
	Mode Mode   // dev , prod
	Env  string // path
}

func (m Mode) IsValid() bool {
	return m == Development || m == Production
}

func ParseMode(s string) (Mode, error) {
	if s == "" {
		return Development, nil
	}

	m := Mode(s)
	if !m.IsValid() {
		return "", fmt.Errorf("%w: %q", ErrInvalidMode, s)
	}
	return m, nil
}

func (o *Option) ApplyDefaults() {
	if o.Env == "" {
		o.Env = DefaultEnvPath
	}
}

func (o Option) Validate() error {
	if !o.Mode.IsValid() {
		return fmt.Errorf("%w: must be 'dev' or 'prod'", ErrInvalidMode)
	}
	return nil
}

func ParseFlags(w io.Writer, args []string) (Option, error) {
	var (
		opt     Option
		modeStr string
	)

	fs := flag.NewFlagSet("apps", flag.ContinueOnError)
	fs.SetOutput(w)

	fs.StringVar(&modeStr, "m", "", "Service mode: dev or prod")
	fs.StringVar(&opt.Env, "e", "", "Service env path")

	if err := fs.Parse(args); err != nil {
		return opt, err
	}

	if fs.NArg() != 0 {
		return opt, fmt.Errorf("positional args not allowed")
	}

	mode, err := ParseMode(modeStr)
	if err != nil {
		return opt, err
	}
	opt.Mode = mode
	opt.ApplyDefaults()

	if err := opt.Validate(); err != nil {
		return opt, err
	}

	return opt, nil
}
