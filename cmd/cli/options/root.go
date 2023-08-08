package options

import (
	"github.com/spf13/cobra"
)

// RootOptions defines the options for the `root` command
type RootOptions struct {
	ConfigPath string
	Verbose    int
}

// AddFlags adds command line flags for the RootOptions struct
func (ro *RootOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&ro.ConfigPath, "config", "c", "", "path to config file")
	cmd.PersistentFlags().CountVarP(&ro.Verbose, "verbose", "v", "log verbosity level (-v=info, -vv=debug, -vvv=trace)")
}
