package readme

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/bom-squad/go-cli/internal/basecli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

const (
	CONFIG_DOC          = "configuration.md"
	CONFIG_COMMAND_PATH = "command"
)

type Readme struct {
	applicationName string
	appConfig       config.ApplicationConfig
	command         *cobra.Command
}

func New(applicationName string, appConfig config.ApplicationConfig) Readme {
	return Readme{
		applicationName: applicationName,
		appConfig:       appConfig,
	}
}

func (r *Readme) GenTree(cmd *cobra.Command) error {
	configDocPath := r.GetPreferredFolder()

	err := doc.GenMarkdownTree(cmd, configDocPath)
	fmt.Printf("generated markdown usage, Err: %s\n", err)

	err = r.GenConfiguration()
	fmt.Printf("generated configuration, Err: %s\n", err)

	return err
}

func (r *Readme) GenConfiguration() error {
	data, err := yaml.Marshal(r.appConfig)
	if err != nil {
		return err
	}

	dataStr := string(data)
	readmeConfigStr := ReplacOutputDirHash(dataStr, r.applicationName, "")

	configDocPath := filepath.Join("docs", CONFIG_DOC)
	templateDocPath := filepath.Join(embbededDir, CONFIG_DOC)
	err = CreateFileFromTemplateFile(templateDocPath, configDocPath, ConfigurationTemplateData{
		ApplicationName: r.applicationName,
		DefaultConfig:   readmeConfigStr,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Readme) SetCommand(command *cobra.Command) {
	r.command = command
}

func (r *Readme) AddArgument(flat *pflag.Flag) {

}

func (r *Readme) Generate(writer io.Writer) error {
	return GenMarkdownCustom(r.command, writer, func(s string) string { return s })
}

func (r *Readme) GetPreferredFolder() string {
	return "docs/command"
}

func (r *Readme) GetPreferredFileName() string {
	bracketsRemoved := strings.ReplaceAll(r.command.Use, "[", "")
	bracketsRemoved = strings.ReplaceAll(bracketsRemoved, "]", "")
	cmdSplit := strings.Split(bracketsRemoved, " ")
	cmdstr := ""
	if len(cmdSplit) > 0 && r.applicationName != cmdSplit[0] {
		cmdstr = "_" + cmdSplit[0]
	}
	return r.applicationName + cmdstr + ".md"
}

func (r *Readme) IsIgnoreHidden() bool {
	return true
}
