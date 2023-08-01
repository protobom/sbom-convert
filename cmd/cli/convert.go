package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/bom-squad/go-cli/cmd/cli/options"
	"github.com/bom-squad/go-cli/internal/convert"
	"github.com/bom-squad/go-cli/pkg/format"
)

var outputDirPermissions = 0o755

func ConvertCommand() *cobra.Command {
	co := &options.ConvertOptions{}
	c := &cobra.Command{
		Use:               "convert",
		Aliases:           []string{""},
		SuggestFor:        []string{},
		Short:             "",
		Long:              ``,
		Example:           ``,
		SilenceErrors:     true,
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Args:              cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConvert(cmd.Context(), co, args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Parent().PersistentPreRunE(cmd.Parent(), args); err != nil {
				return err
			}

			return validateConvertOptions(co, args)
		},
	}

	co.AddFlags(c)

	return c
}

func validateConvertOptions(co *options.ConvertOptions, args []string) error {
	_, err := format.ParseFormat(co.Format, co.Encoding)
	if err != nil {
		return err
	}
	return nil
}

// runConvert is the main entrypoint for the `convert` command
func runConvert(ctx context.Context, co *options.ConvertOptions, args []string) error {
	f, err := os.Open(args[0])
	if err != nil {
		return err
	}

	defer f.Close()

	frmt, err := parseFormat(co.Format, co.Encoding, f)
	if err != nil {
		return err
	}

	out, err := createOutputStream(co.OutputPath, frmt)
	if err != nil {
		return err
	}

	cs := convert.NewService(
		convert.WithFormat(frmt),
	)
	return cs.Convert(ctx, f, out)
}

// parseFormat parses the format string and returns target format
// if format string is empty, it will try to detect the format automatically and return the inverse
func parseFormat(f, e string, r io.ReadSeekCloser) (*format.Format, error) {
	if f == "" {
		df, err := format.DetectFormat(r)
		if err != nil {
			return nil, err
		}

		return df.Inverse()
	}

	format, err := format.ParseFormat(f, e)
	if err != nil {
		return nil, err
	}

	return format, nil
}

func createOutputStream(o string, frmt *format.Format) (io.WriteCloser, error) {
	log := zap.S()

	if o == "" {
		log.Debug("no output path specified, using stdout")
		return os.Stdout, nil
	}

	// NOTE: @mrsufgi: we are autofixing the output file name to match the format
	// this is very opinionated and might not be the best UX (-o ori.file -> ori.file.cdx.json)
	// but it's the simplest way to make sure outputs are always in the right format
	dir := filepath.Dir(o)
	ext := filepath.Ext(o)
	name := strings.TrimSuffix(filepath.Base(o), ext)

	if ext != fmt.Sprintf(".%s", frmt.Encoding()) {
		log.Debug("output path extension does not match format encoding, appending")
		name = fmt.Sprintf("%s%s", name, ext)
		ext = ""
	}
	var output string
	if frmt.Type() == format.CDX {
		output = fmt.Sprintf("%s.cdx", name)
	}

	if frmt.Type() == format.SPDX {
		output = fmt.Sprintf("%s.%s", name, format.SPDX)
	}

	if ext == "" {
		if frmt.Encoding() == format.JSONEncoding {
			log.Debug("output path does not contain a valid format extension, appending")
			output = fmt.Sprintf("%s.%s", output, frmt.Encoding())
		}
	} else {
		output = fmt.Sprintf("%s%s", output, ext)
	}

	log.Debugf("creating output directory: %s", dir)
	if err := os.MkdirAll(dir, os.FileMode(outputDirPermissions)); err != nil {
		return nil, err
	}

	output = filepath.Join(dir, output)

	out, err := os.Create(output)
	log.Debugf("creating output file: %s", output)
	if err != nil {
		return nil, err
	}

	return out, nil
}
