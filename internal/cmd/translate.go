package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	translate "github.com/bom-squad/go-cli/pkg"
	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/spf13/cobra"
)

type TranslateOptions struct {
	FormatString string
	OutputPath   string
}

func (to *TranslateOptions) Format() formats.Format {
	if to.FormatString == "" {
		return ""
	}

	s := to.FormatString
	if strings.ToLower(s) == formats.SPDXFORMAT {
		return translate.DefaultSPDXVersion
	} else if strings.ToLower(s) == formats.CDXFORMAT {
		return translate.DefaultCycloneDXVersion
	}
	return formats.Format(s)
}

func (o *TranslateOptions) AddFlags(command *cobra.Command) {
	command.Flags().StringVar(&o.FormatString, "format", "", "format string")
	command.Flags().StringVarP(&o.OutputPath, "output", "o", "", "path to write the ranslated SBOM")
}

func Translate() *cobra.Command {
	opts := &TranslateOptions{}
	command := &cobra.Command{
		Use:     "translate [flags] sbom.json",
		Short:   "translate an SBOM into another format",
		Long:    "translate an SBOM into another format",
		Example: "protobom translate sbom.spdx.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("no SBOM specified")
			}

			// Open the SBOM for ingestion
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("opening SBOM file: %w", err)
			}

			// Determine where to write the translated SBOM
			var out io.WriteCloser
			if opts.OutputPath == "" {
				out = os.Stdout
			} else {
				var err error
				out, err = os.Create(opts.OutputPath)
				if err != nil {
					return fmt.Errorf("opening output file: %w", err)
				}
				defer out.Close()
			}

			// Create new translator object
			t := translate.NewTranslator()

			// Create the new options
			opts := &translate.TranslationOptions{
				Format: opts.Format(),
			}

			// Call the translation function:
			if err := t.TranslateWithOptions(opts, f, out); err != nil {
				return fmt.Errorf("translating sbom: %w", err)
			}

			// Success !
			return nil
		},
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	opts.AddFlags(command)

	return command
}
