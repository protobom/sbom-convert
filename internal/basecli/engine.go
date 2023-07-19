package basecli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/bom-squad/go-cli/internal/basecli/cmd"
	"github.com/bom-squad/go-cli/internal/basecli/config"
	"github.com/bom-squad/go-cli/internal/basecli/logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ()

// Engine this is the structure that automates and makes the operation for command and config simple
// This struct can parse your config files, combine command flags with config
// It also helps you add some common arguments
type Engine struct {
	config          config.Engine
	cmd             cmd.Engine
	logger          logger.Logger
	globalArguments Arguments
	appConfig       config.ApplicationConfig
}

// New this functions creates a new instance of the basecli.Engine and it takes the application name and rootCmd of the cobra
func New(log logger.Logger, applicationName string, rootcmd *cobra.Command) *Engine {
	engine := &Engine{
		config:          config.New(log, applicationName),
		cmd:             cmd.New(log, rootcmd),
		logger:          log,
		globalArguments: make(Arguments, 0),
	}
	engine.addGenerateCommands()
	return engine
}

// AddCommandAndArguments adds command and arguments to the local command. It also binds with the config structure
func (b *Engine) AddCommandAndArguments(cmd *cobra.Command, arguments Arguments) error {
	if err := b.cmd.AddCommandAndArguments(cmd, arguments); err != nil {
		return err
	}
	cargs := make([]config.Argument, len(arguments))
	convertToConfigArguments(arguments, cargs)
	return b.config.BindFlags(cmd.Flags(), cargs)
}

// AddArguments adds arguments to command, It also binds with the config structure
func (b *Engine) AddArguments(cmd *cobra.Command, arguments Arguments) error {
	if err := b.cmd.AddArguments(cmd, arguments); err != nil {
		return err
	}
	cargs := make([]config.Argument, len(arguments))
	convertToConfigArguments(arguments, cargs)
	return b.config.BindFlags(cmd.Flags(), cargs)
}

// AddGlobalArguments adds arguments to global command and binds them to the config structure.
func (b *Engine) AddGlobalArguments(arguments Arguments) error {
	if err := b.cmd.AddGlobalArguments(arguments); err != nil {
		return err
	}
	cargs := make([]config.Argument, len(arguments))
	convertToConfigArguments(arguments, cargs)
	return b.config.BindFlags(b.cmd.GetCommand().PersistentFlags(), cargs)
}

// AddBasicCommandsAndArguments adds global commands and basic commands and arguments to the root command.
func (b *Engine) AddBasicCommandsAndArguments(options ArgOptions) error {
	args := getBasicGlobalArguments(b.config.GetApplicationName(), options)
	return b.AddGlobalArguments(args)
}

// LoadApplication loads the the configuration of the application and also initializes the logger
func (b *Engine) LoadApplication(applicationConfig config.ApplicationConfig) error {
	if err := b.config.LoadApplicationConfig(applicationConfig); err != nil {
		return err
	}

	b.appConfig = applicationConfig

	if b.logger == nil {
		_, err := b.SetDefaultLogger(true)
		if err != nil {
			return err
		}
	}

	b.DebugConfig(applicationConfig)
	return nil
}

func (b *Engine) DebugConfig(applicationConfig config.ApplicationConfig) {
	if b.config.GetConfig().Logger.Structured {
		data, err := json.Marshal(applicationConfig)
		if err != nil {
			b.logger.Warnf("[config] error marshaling application config %v", err)
		}
		fieldMap := make(map[string]interface{})
		err = json.Unmarshal(data, &fieldMap)
		if err != nil {
			b.logger.Warnf("[config] error unmarshaling application config %v", err)
		}

		b.logger.WithFields(fieldMap).Debugf("Application Configuration")
	} else {
		data, err := yaml.Marshal(applicationConfig)
		if err != nil {
			b.logger.Warnf("[config] error marshaling application config %v", err)
		}
		b.logger.Debugf("%v version: %v", b.config.GetApplicationName(), b.cmd.GetVersion())

		configDelimiter := strings.Repeat("-", 20)
		configLog := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", configDelimiter, "Application Configuration", configDelimiter, string(data), configDelimiter)
		b.logger.Debugf("%s", configLog)
	}
}

func (e *Engine) SetLogger(l logger.Logger) {
	e.config.SetLogger(l)
	e.cmd.SetLogger(l)
	e.logger = l
}

func (e *Engine) SetDefaultLogger(global bool) (logger.Logger, error) {
	loggerConfig := e.config.GetLoggingConfig()
	defaultLogger, err := logger.NewDefaultLogger(loggerConfig)
	if err != nil {
		return nil, err
	}

	e.SetLogger(defaultLogger)
	if global {
		logger.SetGlobalLogger(defaultLogger)
	}
	return defaultLogger, nil
}

func convertToConfigArguments(args Arguments, cargs []config.Argument) {
	for i, arg := range args {
		cargs[i] = arg
	}
}

func Tprintf(tmpl string, data map[string]interface{}) string {
	t := template.Must(template.New("").Parse(tmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return ""
	}
	return buf.String()
}
