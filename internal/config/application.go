package config

import (
	"fmt"
	"strings"

	"github.com/bom-squad/go-cli/internal/basecli"
	translate "github.com/bom-squad/go-cli/pkg"
	"github.com/bom-squad/protobom/pkg/formats"
)

var DefaultConfig = Application{
	Translate: TranslateConfig{
		FormatString: "",
		OutputPath:   "",
		Encoding:     translate.DefaultEncoding,
		URI:          "",
	},
}

// Command line arguments
var TranslateCommandArguments = basecli.Arguments{
	basecli.Argument{LongName: "output", ConfigID: "translate.output", ShortName: "o", Message: "Output path", Default: ""},
	basecli.Argument{LongName: "format", ConfigID: "translate.format", ShortName: "f", Message: fmt.Sprintf("Select Formats, options=%v", formats.ListFormats), Default: ""},
	basecli.Argument{LongName: "ver", ConfigID: "translate.version", ShortName: "V", Message: fmt.Sprintf("Select Specific version, options=%v (default %v)", GetVersionGroups().VersionMap(), GetVersionGroups().DefaultVersions()), Default: ""},
	basecli.Argument{LongName: "encoding", ConfigID: "translate.encoding", ShortName: "E", Message: fmt.Sprintf("Select encoding, options=%v", GetVersionGroups().EncodingMap()), Default: DefaultConfig.Translate.Encoding},
	basecli.Argument{LongName: "uri", ConfigID: "translate.uri", ShortName: "", Message: fmt.Sprintf("Select uri, options=%v", formats.List), Default: "", IsHidden: true},
}

// Configuration structure
type Application struct {
	basecli.BaseApplication `yaml:",omitempty,inline" json:",omitempty,inline" mapstructure:",squash"`
	Translate               TranslateConfig `yaml:"translate,omitempty" json:"translate,omitempty" mapstructure:"translate"`
	version                 string          `yaml:"-" json:"-" mapstructure:"-"`
}

func (cfg *Application) SetVersion(version string) {
	cfg.version = version
}

func (a Application) GetConfigPath() string {
	return a.BaseApplication.ConfigPath
}

type TranslateConfig struct {
	FormatString string `yaml:"format,omitempty" json:"format,omitempty" mapstructure:"format"`
	OutputPath   string `yaml:"output,omitempty" json:"output,omitempty" mapstructure:"output"`
	Version      string `yaml:"version,omitempty" json:"version,omitempty" mapstructure:"version"`
	Encoding     string `yaml:"encoding,omitempty" json:"encoding,omitempty" mapstructure:"encoding"`
	URI          string `yaml:"uri,omitempty" json:"uri,omitempty" mapstructure:"uri"`
}

func (to *TranslateConfig) Format() (formats.Format, error) {
	if to.URI != "" {
		return formats.Format(to.URI), nil
	}

	if to.FormatString == "" {
		return "", fmt.Errorf("format type must be provided") //How do we connect this to the sniffer?
	}

	if to.FormatString != "" {
		for _, defaultFormat := range translate.DefaultsFormatsList {

			if defaultFormat.Type() == strings.ToLower(to.FormatString) {
				var encoding, version string
				if to.Encoding == "" && to.Version == "" {
					return defaultFormat, nil
				}

				if to.Encoding == "" {
					encoding = defaultFormat.Encoding()
				} else {
					encoding = to.Encoding
				}

				if to.Version == "" {
					version = defaultFormat.Version()
				} else {
					version = to.Version
				}

				return formats.Format(fmt.Sprintf("%s+%s;version=%s", defaultFormat.URI(), encoding, version)), nil
			}
		}
	}

	return "", fmt.Errorf("no format found")
}
