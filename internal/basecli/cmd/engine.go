package cmd

import (
	"fmt"
	"reflect"

	"github.com/bom-squad/go-cli/internal/basecli/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Engine struct {
	root *cobra.Command
	log  logger.Logger
}

func New(l logger.Logger, root *cobra.Command) Engine {
	return Engine{
		log:  l,
		root: root,
	}
}

func (c *Engine) SetLogger(log logger.Logger) {
	c.log = log
}

func (c *Engine) GetVersion() string {
	return c.root.Version
}

func (c *Engine) AddGlobalArguments(argument Arguments) error {
	return AddArgumentsToCommand(c.root.PersistentFlags(), argument)
}

func (c *Engine) AddArguments(cmd *cobra.Command, arguments Arguments) error {
	if err := AddArgumentsToCommand(cmd.Flags(), arguments); err != nil {
		return err
	}

	return nil
}

func (c *Engine) AddCommandAndArguments(cmd *cobra.Command, arguments Arguments) error {
	if err := AddArgumentsToCommand(cmd.Flags(), arguments); err != nil {
		return err
	}

	c.root.AddCommand(cmd)
	return nil
}

func (c *Engine) GetCommand() *cobra.Command {
	return c.root
}

func AddArgumentsToCommand(flag *pflag.FlagSet, arguments Arguments) error {
	for i := 0; i < arguments.Len(); i++ {
		arg := arguments.Get(i)
		switch arg.GetDefault().(type) {
		case []string:
			flag.StringSliceP(
				arg.GetLongName(), arg.GetShortName(), arg.GetDefault().([]string),
				arg.GetMessage(),
			)
		case uint:
			flag.UintP(arg.GetLongName(), arg.GetShortName(), arg.GetDefault().(uint),
				arg.GetMessage(),
			)
		case int:
			flag.IntP(
				arg.GetLongName(), arg.GetShortName(), arg.GetDefault().(int),
				arg.GetMessage(),
			)
		case string:
			flag.StringP(
				arg.GetLongName(), arg.GetShortName(), arg.GetDefault().(string),
				arg.GetMessage(),
			)
		case bool:
			flag.BoolP(
				arg.GetLongName(), arg.GetShortName(), arg.GetDefault().(bool),
				arg.GetMessage(),
			)
		case struct{}:
			flag.CountP(
				arg.GetLongName(), arg.GetShortName(), arg.GetMessage(),
			)
		case map[string]string:
			fd := make(map[string]string)
			flag.StringToStringVar(&fd, arg.GetLongName(), arg.GetDefault().(map[string]string), arg.GetMessage())
		default:
			return fmt.Errorf("%v type is not supported yet", reflect.TypeOf(arg.GetDefault()))
		}
		flag.Lookup(arg.GetLongName()).Hidden = arg.GetIsHidden()
	}
	return nil
}
