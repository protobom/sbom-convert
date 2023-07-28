package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/log"
	"sigs.k8s.io/release-utils/version"
)

type commandLineOptions struct {
	logLevel string
}

var commandLineOpts = &commandLineOptions{}

var rootCmd = &cobra.Command{
	Short: "A translation utility to demo the protobom libraries",
	// Long: ``,
	Use:               "protobom",
	SilenceUsage:      false,
	PersistentPreRunE: initLogging,
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&commandLineOpts.logLevel,
		"log-level",
		"info",
		fmt.Sprintf("the logging verbosity, either %s", log.LevelNames()),
	)
	rootCmd.AddCommand(version.WithFont("doom"))
	rootCmd.AddCommand(Translate())
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func initLogging(*cobra.Command, []string) error {
	return log.SetupGlobalLogger(commandLineOpts.logLevel)
}
