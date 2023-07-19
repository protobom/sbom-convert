package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	testCases := []struct {
		name          string
		cfg           Config
		expectedLevel logrus.Level
		err           bool
	}{
		{
			name:          "no level or verbose",
			cfg:           Config{},
			expectedLevel: logrus.InfoLevel,
		},
		{
			name: "level debug",
			cfg: Config{
				Level: "debug",
			},
			expectedLevel: logrus.DebugLevel,
		},
		{
			name: "level info",
			cfg: Config{
				Level: "info",
			},
			expectedLevel: logrus.InfoLevel,
		},
		{
			name: "level warn",
			cfg: Config{
				Level: "warn",
			},
			expectedLevel: logrus.WarnLevel,
		},
		{
			name: "level error",
			cfg: Config{
				Level: "error",
			},
			expectedLevel: logrus.ErrorLevel,
		},
		{
			name: "verbose 0",
			cfg: Config{
				Verbosity: 0,
			},
			expectedLevel: logrus.InfoLevel,
		},
		{
			name: "verbose 1",
			cfg: Config{
				Verbosity: 1,
			},
			expectedLevel: logrus.InfoLevel,
		},
		{
			name: "verbose 2",
			cfg: Config{
				Verbosity: 2,
			},
			expectedLevel: logrus.DebugLevel,
		},
		{
			name: "verbose 3",
			cfg: Config{
				Verbosity: 3,
			},
			expectedLevel: logrus.TraceLevel,
		},
		{
			name: "level and verbose",
			cfg: Config{
				Level:     "warn",
				Verbosity: 2,
			},
			err: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			l, err := NewDefaultLogger(test.cfg)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedLevel, l.GetLevel(), "logger level not set accurately")
			}
		})
	}
}
