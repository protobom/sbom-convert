package generator

import (
	"io"

	"github.com/bom-squad/go-cli/internal/basecli/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Plugin interface {
	SetCommand(command *cobra.Command)
	AddArgument(flags *pflag.Flag)
	Generate(writer io.Writer) error
	GetPreferredFolder() string
	GetPreferredFileName() string
	IsIgnoreHidden() bool
}

type Generator struct {
	log                logger.Logger
	exceptionArguments []string
	exceptionCommands  []string
	newPlugin          func() Plugin
}

func New(log logger.Logger, newPlugin func() Plugin) Generator {
	return Generator{
		log:       log,
		newPlugin: newPlugin,
		exceptionArguments: []string{
			"quiet",
		},
		exceptionCommands: []string{},
	}
}

func (g Generator) Generate(root *cobra.Command) []Plugin {
	actions := make([]Plugin, 0)
	actions = g.appendAction(actions, root, root)
	for _, command := range root.Commands() {
		actions = g.appendAction(actions, root, command)
	}
	return actions
}

func (g Generator) appendAction(plugins []Plugin, root, command *cobra.Command) []Plugin {
	if command.Name() == "help" {
		return plugins
	}
	action := g.newPlugin()
	if (command.Hidden && action.IsIgnoreHidden() && !contains(g.exceptionCommands, command.Name())) || (!(command.Hidden && action.IsIgnoreHidden()) && contains(g.exceptionCommands, command.Name())) {
		return plugins
	}

	action.SetCommand(command)
	argumentMap := make(map[string]struct{})
	command.Flags().VisitAll(func(f *pflag.Flag) {
		if _, ok := argumentMap[f.Name]; ok || (f.Hidden && !contains(g.exceptionArguments, f.Name)) || (!f.Hidden && contains(g.exceptionArguments, f.Name)) {
			if ok {
				g.log.Tracef("[gen] %s appeared multiple times in flags", f.Name)
			}
			return
		}
		argumentMap[f.Name] = struct{}{}
		action.AddArgument(f)
	})
	root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if _, ok := argumentMap[f.Name]; ok || (f.Hidden && !contains(g.exceptionArguments, f.Name)) || (!f.Hidden && contains(g.exceptionArguments, f.Name)) {
			if ok {
				g.log.Tracef("[gen] %s appeared multiple times in flags", f.Name)
			}
			return
		}
		argumentMap[f.Name] = struct{}{}
		action.AddArgument(f)
	})

	plugins = append(plugins, action)
	return plugins
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
