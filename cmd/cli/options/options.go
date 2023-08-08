package options

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Interface interface {
	// AddFlags adds this options' flags to the cobra command.
	AddFlags(cmd *cobra.Command)
}

// BindConfig binds the config file to cobra flags
func BindConfig(v *viper.Viper, cmd *cobra.Command) error {
	cp := cmd.Flag("config").Value.String()
	if cp != "" {
		v.AddConfigPath(filepath.Dir(cp))
		v.SetConfigName(strings.TrimSuffix(filepath.Base(cp), filepath.Ext(cp)))
		v.SetConfigType(filepath.Ext(cp)[1:])
		v.AutomaticEnv()
		if err := v.ReadInConfig(); err != nil {
			if errors.As(err, &viper.ConfigFileNotFoundError{}) {
				zap.S().Debugf("no config file found, using defaults")
			} else {
				return fmt.Errorf("unable to load config: %w", err)
			}
		}

		bindFlags(cmd, v)
	}
	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			_ = v.BindEnv(f.Name, flagToEnvVar(f.Name))
		}
		if !f.Changed && v.IsSet((f.Name)) {
			val := v.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func flagToEnvVar(f string) string {
	f = strings.ToUpper(f)
	return strings.ReplaceAll(f, "-", "_")
}
