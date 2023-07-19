package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bom-squad/go-cli/internal/basecli/logger"
	log "github.com/bom-squad/go-cli/internal/basecli/logger"
	"github.com/bom-squad/go-cli/internal/config"
	translate "github.com/bom-squad/go-cli/pkg"

	"github.com/bom-squad/go-cli/internal/basecli"

	"github.com/spf13/cobra"
)

const (
	TranslateUserExample = `
	{{.appName}} {{.command}} sbom.spdx.json -f cyclonedx           translate SPDX to CycloneDX
	{{.appName}} {{.command}} sbom.cdx.json  -f spdx                translate CycloneDX to SPDX
	{{.appName}} {{.command}} sbom.cdx.json  -f cyclonedx -V 1.5    select specific version
	{{.appName}} {{.command}} sbom.cdx.json -f spdx -E text         select specific encoding
	{{.appName}} {{.command}} sbom.cdx.json  -o sbom.spdx.json      output to file
	`
	CommandName               = "translate"
	TranslateLongDescription  = "Translate between SBOM formats, Bridging the gap between spdx and cyclonedx"
	TranslateShortDescription = "translate an SBOM into another format"
)

var translateCmd = &cobra.Command{
	Version: version, //For root command setup the version
	Long:    TranslateLongDescription,
	Use:     fmt.Sprintf("%s [path]", ApplicationName),
	Short:   TranslateShortDescription,
	Example: basecli.Tprintf(TranslateUserExample, map[string]interface{}{
		"appName": ApplicationName,
		"command": "",
	}),
	DisableAutoGenTag: true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Args:              cobra.MinimumNArgs(1), // Require SBOM path as argument
	Hidden:            false,
	SilenceUsage:      true,
	SilenceErrors:     false,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		sbom := args[0]
		conf.SetVersion(version)
		log.Info("Running protobom")
		return Translate(sbom, &conf, log.GetGlobalLogger())
	},
}

func Translate(sbom string, config *config.Application, l logger.Logger) error {

	if sbom == "" {
		return errors.New("Empty SBOM path")
	}

	// Open the SBOM for ingestion
	f, err := os.Open(sbom)
	if err != nil {
		return fmt.Errorf("opening SBOM file: %w", err)
	}

	// Determine where to write the translated SBOM
	var out io.WriteCloser
	if config.Translate.OutputPath == "" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.Create(config.Translate.OutputPath)
		if err != nil {
			return fmt.Errorf("opening output file: %w", err)
		}
		defer out.Close()
	}

	format, err := config.Translate.Format()
	if err != nil {
		return err
	}

	l.Infof("translating %s in to %s ...", sbom, format)

	// Create new translator object
	t := translate.NewTranslator()

	// Create the new options
	translateOpts := &translate.TranslationOptions{
		Format: format,
	}

	l.Debugf("Translation options, %+v", translateOpts)

	// Call the translation function:
	if err := t.TranslateWithOptions(translateOpts, f, out); err != nil {
		return fmt.Errorf("translating: %w", err)
	}

	if config.Translate.OutputPath != "" {
		l.Infof("Translation success, path: %s", config.Translate.OutputPath)
	}

	// Success !
	return nil

}
