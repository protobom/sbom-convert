package cli

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/bom-squad/go-cli/cmd/cli/options"
)

func ConvertCommand() *cobra.Command {
	co := &options.ConvertOptions{}
	c := &cobra.Command{
		Use:               "",
		Aliases:           []string{""},
		SuggestFor:        []string{},
		Short:             "",
		Long:              ``,
		Example:           ``,
		SilenceErrors:     true,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Args:              cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgsEnd := len(args)
			if cmd.ArgsLenAtDash() > -1 {
				cmdArgsEnd = cmd.ArgsLenAtDash()
			}

			return convert(cmd.Context(), co, args[:cmdArgsEnd], args[cmdArgsEnd:])
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return validateConvertOptions(co, args)
		},
	}

	co.AddFlags(c)

	return c
}

func validateConvertOptions(co *options.ConvertOptions, args []string) error {
	log.Println(co, args)
	return nil
}

func convert(ctx context.Context, co *options.ConvertOptions, args, ptargs []string) error {
	log.Println(ctx, co, args, ptargs)
	return nil
}
