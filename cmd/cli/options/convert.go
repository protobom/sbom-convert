package options

import (
	"github.com/spf13/cobra"
)

// ConvertOptions defines the options for the `convert` command
type ConvertOptions struct {
	SBOMName string
}

// AddFlags adds command line flags for the ConvertOptions struct
func (o *ConvertOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.SBOMName, "name", "n", "", "name of convertd SBOM document.")
}
