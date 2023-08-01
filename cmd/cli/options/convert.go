package options

import (
	"github.com/spf13/cobra"
)

// ConvertOptions defines the options for the `convert` command
type ConvertOptions struct {
	Format     string
	Encoding   string
	OutputPath string
}

// AddFlags adds command line flags for the ConvertOptions struct
func (o *ConvertOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Format, "format", "f", "", "format string")
	cmd.Flags().StringVarP(&o.Encoding, "encoding", "e", "", "output encoding")
	cmd.Flags().StringVarP(&o.OutputPath, "output", "o", "", "path to write the converted SBOM")
}
