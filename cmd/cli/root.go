package cli

import (
	"fmt"
	"os"

	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
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
	rootCmd := &cobra.Command{
		Use:     "",
		Version: version,
		Short:   "",
		Long:    ``,
		Run:     func(cmd *cobra.Command, args []string) {},
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	// Commands
	rootCmd.AddCommand(ConvertCommand())

	// Manpages
	rootCmd.AddCommand(ManCommand(rootCmd))
	return rootCmd
}
