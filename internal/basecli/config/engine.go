package config

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/adrg/xdg"
	"github.com/bom-squad/go-cli/internal/basecli/logger"

	"github.com/mitchellh/go-homedir"

	"github.com/spf13/viper"
)

var ErrApplicationConfigNotFound = fmt.Errorf("application config not found")

type ApplicationConfig interface {
	GetConfigPath() string
}

// config seperate3e
type BaseApplication struct {
	ConfigPath string        `yaml:"config_path,omitempty" json:"config_path,omitempty" mapstructure:"config_path"`
	Logger     logger.Config `yaml:"logger,omitempty" json:"logger,omitempty"  mapstructure:"logger"`
}

func (b *BaseApplication) GetConfigPath() string {
	return b.ConfigPath
}

func (b *BaseApplication) GetLoggingConfig() logger.Config {
	return b.Logger
}

type Engine struct {
	*viper.Viper
	applicationName string
	config          BaseApplication
	log             logger.Logger
}

func New(log logger.Logger, applicationName string) Engine {
	return Engine{
		Viper:           viper.New(),
		applicationName: applicationName,
		log:             log,
	}
}

func (c *Engine) GetApplicationName() string {
	return c.applicationName
}

func (c *Engine) SetLogger(log logger.Logger) {
	c.log = log
}

func (c *Engine) GetConfig() BaseApplication {
	return c.config
}

func (c *Engine) LoadApplicationConfig(applicationConfig ApplicationConfig) error {
	if err := c.Unmarshal(&c.config); err != nil {
		return fmt.Errorf("unable to parse basic config: %v", err)
	}

	if err := c.readConfig(&c.config); err != nil && !errors.Is(err, ErrApplicationConfigNotFound) {
		return err
	}

	if err := c.Unmarshal(&c.config); err != nil {
		return fmt.Errorf("unable to parse basic config: %v", err)
	}

	if err := c.Unmarshal(applicationConfig); err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}

	return nil
}

func (c *Engine) GetLoggingConfig() logger.Config {
	return c.config.Logger
}

// readConfig attempts to read the given config path from disk or discover an alternate store location
// nolint:funlen
func (c *Engine) readConfig(applicationConfig ApplicationConfig) error {
	var err error
	configPath := applicationConfig.GetConfigPath()
	c.AutomaticEnv()
	c.SetEnvPrefix(c.applicationName)
	// allow for nested options to be specified via environment variables
	// e.g. pod.context = APPNAME_POD_CONTEXT
	c.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	// use explicitly the given user config
	if configPath != "" {
		c.SetConfigFile(configPath)
		if err := c.ReadInConfig(); err != nil {
			return fmt.Errorf("unable to read application config=%q : %w", configPath, err)
		}
		// don't fall through to other options if the config path was explicitly provided
		return nil
	}

	// start searching for valid configs in order...

	// 1. look for .<appname>.yaml (in the current directory)
	c.AddConfigPath(".")
	c.SetConfigName("." + c.applicationName)
	if err = c.ReadInConfig(); err == nil {
		return nil
	} else if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return fmt.Errorf("unable to parse config=%q: %w", c.ConfigFileUsed(), err)
	}

	// 2. look for .<appname>/config.yaml (in the current directory)
	c.AddConfigPath("." + c.applicationName)
	c.SetConfigName(c.applicationName)
	if err = c.ReadInConfig(); err == nil {
		return nil
	} else if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return fmt.Errorf("unable to parse config=%q: %w", c.ConfigFileUsed(), err)
	}

	// 3. look for ~/.<appname>.yaml
	home, err := homedir.Dir()
	if err == nil {
		c.AddConfigPath(home)
		c.SetConfigName("." + c.applicationName)
		if err = c.ReadInConfig(); err == nil {
			return nil
		} else if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return fmt.Errorf("unable to parse config=%q: %w", c.ConfigFileUsed(), err)
		}
	}

	// 4. look for <appname>/config.yaml in xdg locations (starting with xdg home config dir, then moving upwards)
	c.AddConfigPath(path.Join(xdg.ConfigHome, c.applicationName))
	for _, dir := range xdg.ConfigDirs {
		c.AddConfigPath(path.Join(dir, c.applicationName))
	}
	c.SetConfigName(c.applicationName)
	if err = c.ReadInConfig(); err == nil {
		return nil
	} else if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return fmt.Errorf("unable to parse config=%q: %w", c.ConfigFileUsed(), err)
	}
	return ErrApplicationConfigNotFound
}
