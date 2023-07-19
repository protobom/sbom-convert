//go:build dev
// +build dev

package basecli

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/bom-squad/go-cli/internal/basecli/generator"
	"github.com/bom-squad/go-cli/internal/basecli/generator/plugin/github"
	"github.com/bom-squad/go-cli/internal/basecli/generator/plugin/readme"
	"github.com/spf13/cobra"
)

////go:build dev
//// +build dev

var (
	GithubCommand = "github"
	ReadmeCommand = "readme"
	AllCommands   = []string{
		GithubCommand,
		ReadmeCommand,
	}
)

func (b *Engine) addGenerateCommands() {
	b.AddCommandAndArguments(&cobra.Command{
		Use:    fmt.Sprintf("generate %s", AllCommands),
		Hidden: true,
		Args:   cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case GithubCommand:
				b.githubCommand()
			case ReadmeCommand:
				b.readmeCommand()
			case "all":
				b.githubCommand()
				b.readmeCommand()
			default:
				fmt.Printf("Unknown subcommand, supported: %s\n", AllCommands)
				return
			}
		},
	}, Arguments{})
}

func (b *Engine) readmeCommand() {
	readmeGenerator := generator.New(b.logger, func() generator.Plugin {
		p := readme.New(b.config.GetApplicationName(), b.appConfig)
		p.GenConfiguration()
		return &p
	})
	root := b.cmd.GetCommand()
	root.Hidden = false
	generateAction(root, readmeGenerator)
	root.Hidden = true

	b.logger.Infof("Generating readme docs")
}

func (b *Engine) githubCommand() {
	githubGenerator := generator.New(b.logger, func() generator.Plugin {
		p := github.New(b.config.GetApplicationName())
		return &p
	})
	generateAction(b.cmd.GetCommand(), githubGenerator)
	b.logger.Infof("Generated github actions")
}

func generateAction(command *cobra.Command, gen generator.Generator) error {
	plugins := gen.Generate(command)
	for _, a := range plugins {
		os.MkdirAll(a.GetPreferredFolder(), os.ModePerm)
		outputPath := path.Join(a.GetPreferredFolder(), a.GetPreferredFileName())
		if _, err := os.Stat(outputPath); !errors.Is(err, os.ErrNotExist) {
			os.Remove(outputPath)
		}
		f, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		err = a.Generate(f)
		if err != nil {
			return err
		}
	}
	return nil
}
