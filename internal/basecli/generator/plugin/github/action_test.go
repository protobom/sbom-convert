package github

import (
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// -- string Value
type stringValue string

func (s stringValue) Set(val string) error {
	return nil
}

func (s stringValue) Type() string {
	return "string"
}

func (s stringValue) String() string {
	return string(s)
}

// -- stringSlice Value
type stringSliceValue []string

func (s stringSliceValue) Set(val string) error {
	return nil
}

func (s stringSliceValue) Type() string {
	return "stringSlice"
}

func (s stringSliceValue) String() string {
	return "[" + strings.Join([]string(s), ",") + "]"
}

func (s stringSliceValue) Append(val string) error {
	return nil
}

func (s stringSliceValue) Replace(val []string) error {
	return nil
}

func (s stringSliceValue) GetSlice() []string {
	return nil
}

// -- string Value
type intValue int

func (s intValue) Set(val string) error {
	return nil
}

func (s intValue) Type() string {
	return "int"
}

func (s intValue) String() string {
	return strconv.Itoa(int(s))
}

type boolValue bool

func (s boolValue) Set(val string) error {
	return nil
}

func (s boolValue) Type() string {
	return "bool"
}

func (s boolValue) String() string {
	return strconv.FormatBool(bool(s))
}

type uintValue uint

func (s uintValue) Set(val string) error {
	return nil
}

func (s uintValue) Type() string {
	return "uint"
}

func (s uintValue) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

func TestAddArgument(t *testing.T) {
	ae := New("examplecli")
	testcases := []struct {
		testName    string
		name        string
		description string
		required    bool
		def         pflag.Value
		expected    string
		expectedDef interface{}
	}{
		{
			testName:    "argument with string default",
			name:        "arg_1",
			description: "",
			def:         stringValue("default"),
			expected:    "${{ inputs.arg_1 && '--arg_1=' }}${{ inputs.arg_1 }}",
			expectedDef: "default",
		},
		{
			testName:    "argument with string default",
			name:        "arg_2",
			description: "",
			def:         stringSliceValue([]string{"default1", "default2"}),
			expected:    "${{ inputs.arg_2 && '--arg_2=' }}${{ inputs.arg_2 }}",
			expectedDef: "default1,default2",
		},
		{
			testName:    "argument with int default",
			name:        "arg_3",
			description: "",
			def:         intValue(1),
			expected:    "${{ inputs.arg_3 && '--arg_3=' }}${{ inputs.arg_3 }}",
			expectedDef: 1,
		},
		{
			testName:    "argument with int default",
			name:        "arg_4",
			description: "",
			def:         boolValue(true),
			expected:    "${{ inputs.arg_4 && '--arg_4=' }}${{ inputs.arg_4 }}",
			expectedDef: true,
		},
		{
			testName:    "argument with uint default",
			name:        "arg_5",
			description: "",
			def:         uintValue(uint(1)),
			expected:    "${{ inputs.arg_5 && '--arg_5=' }}${{ inputs.arg_5 }}",
			expectedDef: uint64(1),
		},
	}

	for i, test := range testcases {
		t.Log(test.def)
		t.Run(test.testName, func(t *testing.T) {

			ae.AddArgument(&pflag.Flag{
				Name:  test.name,
				Usage: test.description,
				Value: test.def,
			})
			assert.Equal(t, test.expected, ae.Runs.Args[i])
			val, _ := ae.inputsMap.Get(strings.ReplaceAll(test.name, ".", "-"))
			assert.Equal(t, test.expectedDef, val.Default)
		})
	}
}

func TestSetSubCommand(t *testing.T) {
	testcases := []struct {
		testName    string
		commandName string
		expectedArg []string
	}{
		{
			testName:    "normal command should be working",
			commandName: "simple",
			expectedArg: []string{
				"simple",
			},
		},
		{
			testName:    "normal command with named arg should be working",
			commandName: "simple [target]",
			expectedArg: []string{
				"simple",
				"${{ inputs.target}}",
			},
		},
	}

	for _, test := range testcases {
		t.Run(test.testName, func(t *testing.T) {
			ae := New("mockcli")
			ae.SetCommand(&cobra.Command{
				Use:  test.commandName,
				Long: "some description",
			})
			assert.Equal(t, test.expectedArg, ae.Runs.Args)
		})
	}
}

func TestDescriptionExtraction(t *testing.T) {
	tests := []struct {
		Description         string
		Key                 string
		ExpectedDescription string
	}{
		{
			Description: `  {{.appName}} {{.command}} <target>

			<target> Target object name format=[<image:tag>, <dir path>, <git url>]

			{{.appName}} {{.command}} alpine:latest                         verify target against signed attestation of sbom
			{{.appName}} {{.command}} alpine:latest -i attest-slsa          verify target against signed attestation of SLSA provenance
			{{.appName}} {{.command}} alpine:latest -vv                     show verbose debug information`,
			Key:                 "target",
			ExpectedDescription: "Target object name format=[<image:tag>, <dir path>, <git url>]",
		},
		{
			Description: `  {{.appName}} {{.command}} <target>

			<target> Target object name format=[<image:tag>, <dir path>, <git url>, <target>]
		  
			{{.appName}} {{.command}} alpine:latest                         verify target against signed attestation of sbom
			{{.appName}} {{.command}} alpine:latest -i attest-slsa          verify target against signed attestation of SLSA provenance
			{{.appName}} {{.command}} alpine:latest -vv                     show verbose debug information`,
			Key:                 "target",
			ExpectedDescription: "Target object name format=[<image:tag>, <dir path>, <git url>, <target>]",
		},
	}
	for _, test := range tests {
		got := extractDescription(test.Description, test.Key)
		assert.Equal(t, test.ExpectedDescription, got)
	}
}
