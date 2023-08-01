package cli

import (
	"fmt"
	"os"

	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/bom-squad/go-cli/cmd/cli/options"
	"github.com/bom-squad/go-cli/pkg/log"
)

var version = "dev"

func ManCommand(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mcobra.NewManPage(1, root)
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
			return err
		},
	}

	return cmd
}

func NewRootCmd() *cobra.Command {
	ro := &options.RootOptions{}
	rootCmd := &cobra.Command{
		Use:     "sbom-convert",
		Version: version,
		Short:   "",
		Long:    ``,
		Run:     func(cmd *cobra.Command, args []string) {},
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateRootOptions(ro); err != nil {
				return err
			}
			return setupLogger(ro)
		},
		SilenceErrors: true,
	}

	ro.AddFlags(rootCmd)

	// Commands
	rootCmd.AddCommand(ConvertCommand())

	// Manpages
	rootCmd.AddCommand(ManCommand(rootCmd))
	return rootCmd
}

func validateRootOptions(_ *options.RootOptions) error {
	return nil
}

func setupLogger(ro *options.RootOptions) error {
	level := zap.WarnLevel
	if ro.Verbose {
		level = zap.InfoLevel
	}

	if ro.Debug {
		level = zap.DebugLevel
	}

	log, err := log.NewLogger(
		log.WithLevel(level),
		log.WithGlobalLogger(),
	)
	if err != nil {
		return err
	}

	log.Debug("logger initialized")
	return nil
}
