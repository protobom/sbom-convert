package github

import (
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"

	"github.com/bom-squad/go-cli/internal/basecli/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"
)

const (
	ImageReg = "bom-squad"
)

type actionInput struct {
	Description        string      `yaml:"description"`
	Required           bool        `yaml:"required,omitempty"`
	Default            interface{} `yaml:"default,omitempty"`
	DeprecationMessage string      `yaml:"deprecationMessage"`
}

type actionRun struct {
	Using string   `yaml:"using"`
	Image string   `yaml:"image"`
	Args  []string `yaml:"args"`
}

type branding struct {
	Icon  string
	Color string
}

type Action struct {
	Name        string `yaml:"name"`
	fileName    string
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	// Output
	Outputs *yaml.Node `yaml:"outputs,omitempty"`
	// Inputs the inputs to the command line argument
	Inputs    *yaml.Node `yaml:"inputs"`
	inputsMap *orderedmap.OrderedMap[string, actionInput]

	Runs     actionRun `yaml:"runs"`
	Branding branding
}

func New(name string) Action {
	return Action{
		Name:      name,
		inputsMap: orderedmap.New[string, actionInput](),
		Runs: actionRun{
			Using: "docker",
			Image: fmt.Sprintf("docker://%s/%s:latest", ImageReg, name),
			Args:  make([]string, 0),
		},
		Branding: branding{
			Icon:  "shield",
			Color: "green",
		},
	}
}

func (a *Action) IsIgnoreHidden() bool {
	return true
}

func (a *Action) SetCommand(command *cobra.Command) {
	name := command.Use
	description := command.Example
	a.Description = command.Long
	bracketsRemoved := strings.ReplaceAll(name, "[", "")
	bracketsRemoved = strings.ReplaceAll(bracketsRemoved, "]", "")
	cmdSplit := strings.Split(bracketsRemoved, " ")
	for i, cmd := range cmdSplit {
		if i == 0 {
			a.Runs.Args = append(a.Runs.Args, cmd)
			continue
		}
		cmd = strings.ToLower(cmd)
		a.inputsMap.Set(cmd, actionInput{
			Description: extractDescription(description, cmd),
			Required:    true,
		})
		cmdName := strings.ReplaceAll(cmd, ".", "-")
		a.addSubCommand(cmdName)
	}
	if len(cmdSplit) > 0 {
		a.Name = a.Name + "_" + cmdSplit[0]
	}
	a.fileName = a.Name
	if val, ok := CommandMetaData[a.Name]; ok {
		a.Name = val.Name
		a.Author = val.Author
		if len(val.Outputs) > 0 {
			a.Outputs = getOutputNode(val.Outputs)
		}
	} else {
		logger.Warnf("You did not register %s in your action_constants.go", a.Name)
	}
}

func extractDescription(description, key string) string {
	key = "<" + key + ">"
	lines := strings.Split(description, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), key) {
			return strings.TrimSpace(strings.Replace(line, key, "", 1))
		}
	}
	return ""
}

func getOutputNode(output []actionOutput) *yaml.Node {
	mNode := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{},
	}
	for _, val := range output {
		mNode.Content = append(mNode.Content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: val.Name,
		})
		mNode.Content = append(mNode.Content, &yaml.Node{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{
					Kind:  yaml.ScalarNode,
					Value: "description",
				},
				{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: val.Description,
				},
			},
		})
	}
	return mNode
}

func (a *Action) GetPreferredFolder() string {
	return path.Join("ci_generate", "github", a.fileName)
}

func (a *Action) GetPreferredFileName() string {
	return "action.yaml"
}

func (a *Action) addSubCommand(name string) {
	a.Runs.Args = append(a.Runs.Args, fmt.Sprintf("${{ inputs.%s}}", name))
}

func (a *Action) Generate(writer io.Writer) error {
	mNode := getMapNode()
	for v := a.inputsMap.Oldest(); v != nil; v = v.Next() {

		mNode.Content = append(mNode.Content, getActionInput(v.Key, v.Value)...)
	}
	a.Inputs = mNode
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(1)
	err := encoder.Encode(a)
	if err != nil {
		return err
	}
	return nil
}

func (a *Action) AddArgument(flag *pflag.Flag) {
	name := flag.Name
	description := flag.Usage

	inputName := strings.ReplaceAll(name, ".", "-")
	var defaultValue interface{}

	a.inputsMap.Set(inputName, actionInput{
		Description: description,
		Required:    false,
		Default:     defaultValue,
	})
	a.Runs.Args = append(a.Runs.Args, fmt.Sprintf("${{ inputs.%s && '--%s=' }}${{ inputs.%s }}", inputName, name, inputName))
}

func (a *Action) SetImage(imageName string) {
	a.Runs.Image = imageName
}

func (a *Action) GetName() string {
	return a.Name
}

func getActionInput(key string, value actionInput) []*yaml.Node {
	requiredNode := []*yaml.Node{}
	if value.Required || value.DeprecationMessage != "" {
		requiredNode = append(requiredNode, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "required",
		})
		requiredNode = append(requiredNode, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!bool",
			Value: strconv.FormatBool(value.Required),
		})
	}

	defaultNode := []*yaml.Node{}
	switch fmt.Sprintf("%T", value.Default) {
	case "string":
		defaultNode = append(defaultNode, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "default",
		})
		defaultNode = append(defaultNode, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: value.Default.(string),
		})
		if value.Default.(string) == "" {
			defaultNode = []*yaml.Node{}
		}
	case "uint64", "int":
		defaultNode = append(defaultNode, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "default",
		})
		defaultNode = append(defaultNode, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!int",
			Value: strconv.Itoa(value.Default.(int)),
		})
	case "bool":
		defaultNode = append(defaultNode, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "default",
		})
		defaultNode = append(defaultNode, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!bool",
			Value: strconv.FormatBool(value.Default.(bool)),
		})
	}

	node := []*yaml.Node{
		{
			Kind:  yaml.ScalarNode,
			Value: key,
		},
		{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{
					Kind:  yaml.ScalarNode,
					Value: "description",
				},
				{
					Kind:  yaml.ScalarNode,
					Value: value.Description,
				},
			},
		},
	}

	if value.DeprecationMessage != "" {
		deprecationMessage := []*yaml.Node{}

		deprecationMessage = append(deprecationMessage, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "deprecationMessage",
		})
		deprecationMessage = append(deprecationMessage, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: value.DeprecationMessage,
		})
		node[1].Content = append(node[1].Content, deprecationMessage...)
	}

	node[1].Content = append(node[1].Content, requiredNode...)
	node[1].Content = append(node[1].Content, defaultNode...)

	return node
}

func getMapNode() *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{},
	}
}
