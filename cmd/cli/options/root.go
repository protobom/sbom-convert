package options

import (
	"github.com/spf13/cobra"
)

// RootOptions defines the options for the `root` command
type RootOptions struct {
	ConfigPath string
	Verbose    bool
	Debug      bool
}

// AddFlags adds command line flags for the RootOptions struct
func (ro *RootOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&ro.ConfigPath, "config", "c", "", "path to config file")
	cmd.PersistentFlags().BoolVarP(&ro.Verbose, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().BoolVar(&ro.Debug, "debug", false, "debug output")
}
