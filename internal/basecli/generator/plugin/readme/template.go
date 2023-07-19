package readme

import (
	"embed"
	"os"
	"text/template"
)

//go:embed templates
var EmbbededFs embed.FS

const embbededDir = "templates"

type ConfigurationTemplateData struct {
	ApplicationName string
	DefaultConfig   string
}

func CreateFileFromTemplateFile(srcTemplate, dstFile string, data interface{}) error {
	v, err := EmbbededFs.ReadFile(srcTemplate)
	tmpl := template.New(srcTemplate)
	tmpl, err = tmpl.Parse(string(v))
	if err != nil {
		return err
	}

	err = CreateFileUsingTemplate(tmpl, dstFile, data)
	if err != nil {
		return err
	}

	return nil
}

func CreateFileUsingTemplate(t *template.Template, filename string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}
