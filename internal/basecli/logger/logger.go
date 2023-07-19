package logger

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	createPerm = 0o755
)

type Config struct {
	Structured   bool   `yaml:"structured,omitempty" json:"structured,omitempty" mapstructure:"structured"`
	Quiet        bool   `yaml:"quiet,omitempty" json:"quiet,omitempty" mapstructure:"quiet"`
	Level        string `yaml:"level,omitempty" json:"level,omitempty" mapstructure:"level"`
	FileLocation string `yaml:"file,omitempty" json:"file-location,omitempty" mapstructure:"file"`
	Verbosity    int    `yaml:"verbose,omitempty" json:"verbose,omitempty"  mapstructure:"verbose"`
}

func (c Config) ParseLevel() (logrus.Level, error) {
	if c.Quiet {
		return logrus.PanicLevel, nil
	} else {
		if c.Level != "" {
			if c.Verbosity > 0 {
				return logrus.PanicLevel, fmt.Errorf("cannot explicitly set log level (cfg file or env var) and use -v flag together")
			}
			lvl, err := logrus.ParseLevel(strings.ToLower(c.Level))
			if err != nil {
				return logrus.PanicLevel, fmt.Errorf("bad log level configured (%q): %w", c.Level, err)
			}
			return lvl, nil
		} else {
			// set the log level implicitly
			switch v := c.Verbosity; {
			case v == 1:
				return logrus.InfoLevel, nil
			case v == 2:
				return logrus.DebugLevel, nil
			case v == 3:
				return logrus.TraceLevel, nil
			default:
				return logrus.InfoLevel, nil
			}
		}
	}
}

func NewDefaultLogger(cfg Config) (Logger, error) {
	logrusLevel, err := cfg.ParseLevel()
	if err != nil {
		return nil, err
	}

	logrusCfg := LogrusConfig{
		EnableConsole: (cfg.FileLocation == "" || cfg.Verbosity > 0) && !cfg.Quiet,
		EnableFile:    cfg.FileLocation != "",
		Level:         logrusLevel,
		Structured:    cfg.Structured,
		FileLocation:  cfg.FileLocation,
	}

	return NewLogrusLogger(logrusCfg)
}
