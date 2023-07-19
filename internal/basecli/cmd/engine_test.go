package cmd

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type ArgumentDummy struct {
	ConfigID  string
	ShortName string
	LongName  string
	Message   string
	Default   interface{}
	IsHidden  bool
}

func (a ArgumentDummy) GetConfigID() string {
	return a.ConfigID
}

func (a ArgumentDummy) GetLongName() string {
	return a.LongName
}

func (a ArgumentDummy) GetShortName() string {
	return a.ShortName
}

func (a ArgumentDummy) GetMessage() string {
	return a.Message
}

func (a ArgumentDummy) GetDefault() interface{} {
	return a.Default
}

func (a ArgumentDummy) GetIsHidden() bool {
	return a.IsHidden
}

type ArgumentsArr []ArgumentDummy

func (a ArgumentsArr) Get(i int) Argument {
	return a[i]
}

func (a ArgumentsArr) Len() int {
	return len(a)
}

var root = &cobra.Command{
	Use:   "application_test",
	Short: "my application is doing xyz",
	Long:  "my application is achieving xyz by doing abc and then def",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Ran successfully")
	},
}

func TestAddArgumentsToGlobal(t *testing.T) {

	c := New(logrus.New(), root)
	args := ArgumentsArr([]ArgumentDummy{
		{
			ConfigID: "conf",
			LongName: "conf",
			Message:  "config path",
			Default:  "",
		},
		{
			ConfigID: "verb",
			LongName: "verb",
			Message:  "verbose log",
			Default:  false,
		},
		{
			ConfigID: "count",
			LongName: "count",
			Message:  "count of something",
			Default:  3,
		},
	})
	c.AddGlobalArguments(args)

	_, err := c.GetCommand().PersistentFlags().GetString("b")
	assert.Error(t, err, "invalid commands should return error")

	_, err = c.GetCommand().PersistentFlags().GetString(args[0].GetLongName())
	assert.NoError(t, err, "config should be added configID: %v", args[0].GetConfigID())

	_, err = c.GetCommand().PersistentFlags().GetBool(args[1].GetLongName())
	assert.NoError(t, err, "verbose should be added configID: %v", args[1].GetConfigID())

	_, err = c.GetCommand().PersistentFlags().GetInt(args[2].GetLongName())
	assert.NoError(t, err, "count should be present configID: %v", args[2].GetConfigID())

}

var versionString = `
Application Cli Version: %s
Base Cli Version       : v1.0
`

func GetMockCommand(version string) cobra.Command {
	ver := fmt.Sprintf(versionString, version)
	var mockCmd = cobra.Command{
		Use:   "mock",
		Short: version,
		Long:  ver,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ver)
		},
	}
	return mockCmd
}

func TestAddCommandAndArguments(t *testing.T) {

	t.Run("Command should be added", func(t *testing.T) {
		var root = &cobra.Command{
			Use:   "application_test",
			Short: "my application is doing xyz",
			Long:  "my application is achieving xyz by doing abc and then def",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Ran successfully")
			},
		}
		mockCmd := GetMockCommand("s1.1")
		found := false
    c := New(logrus.New(), root)
		c.AddCommandAndArguments(&mockCmd, ArgumentsArr([]ArgumentDummy{}))
		for _, cm := range c.GetCommand().Commands() {
			if cm.Use == "mock" {

				found = true
			}
		}
		assert.Equal(t, true, found, "command should be added")
	})
	testCases := []struct {
		Name        string
		CommandArgs []string
		Arguments   []ArgumentDummy
		Expected    []interface{}
		Err         error
	}{
		{
			Name: "string argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: ""},
			},
			CommandArgs: []string{"mock", "-c=helloworld"},
			Expected:    []interface{}{"helloworld"},
		},
		{
			Name: "string slice argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: []string{""}},
			},
			CommandArgs: []string{"mock", "-c=helloworld"},
			Expected:    []interface{}{[]string{"helloworld"}},
		},
		{
			Name: "int argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: 0},
			},
			CommandArgs: []string{"mock", "-c=4"},
			Expected:    []interface{}{4},
		},

		{
			Name: "uint argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: uint(0)},
			},
			CommandArgs: []string{"mock", "-c=4"},
			Expected:    []interface{}{uint(4)},
		},
		{
			Name: "bool argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: false},
			},
			CommandArgs: []string{"mock", "--config"},
			Expected:    []interface{}{true},
		},
		{
			Name: "struct argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: struct{}{}},
			},
			CommandArgs: []string{"mock", "-ccc"},
			Expected:    []interface{}{3},
		},
		{
			Name: "multiple argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: false},
				{LongName: "verbose", ConfigID: "verbose", ShortName: "v", Message: "some verbose", Default: struct{}{}},
			},
			CommandArgs: []string{"mock", "-vvv", "-c"},
			Expected:    []interface{}{true, 3},
		},
		{
			Name: "should not accept arguments that do not have supported type",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: false},
				{LongName: "verbose", ConfigID: "verbose", ShortName: "v", Message: "some verbose", Default: uint16(43)},
			},
			CommandArgs: []string{"mock"},
			Expected:    []interface{}{},
			Err:         fmt.Errorf("uint16 is not supported"),
		},
	}
	t.Run("Arguments should be set and pass", func(t *testing.T) {
		for _, test := range testCases {
			t.Run(test.Name, func(t *testing.T) {
				var root = &cobra.Command{
					Use:   "application_test",
					Short: "my application is doing xyz",
					Long:  "my application is achieving xyz by doing abc and then def",
					Run: func(cmd *cobra.Command, args []string) {
						fmt.Println("Ran successfully")
					},
				}
        mockCmd := GetMockCommand("s11")
        c := New(logrus.New(), root)
				err := c.AddCommandAndArguments(&mockCmd, ArgumentsArr(test.Arguments))
				if test.Err != nil {
					assert.Error(t, err, "error was expected for %v reason", test.Err)
				}
        root.SetArgs(test.CommandArgs)
				root.Execute()
				for i, arg := range test.Arguments {
					var err error
					var ans interface{}
					switch arg.Default.(type) {
					case []string:
						ans, err = mockCmd.Flags().GetStringSlice(arg.LongName)
					case uint:
						ans, err = mockCmd.Flags().GetUint(arg.LongName)
					case int:
						ans, err = mockCmd.Flags().GetInt(arg.LongName)
					case string:
						ans, err = mockCmd.Flags().GetString(arg.LongName)
					case bool:
						ans, err = mockCmd.Flags().GetBool(arg.LongName)
					case struct{}:
						ans, err = mockCmd.Flags().GetCount(arg.LongName)
					}
					assert.NoError(t, err, "error should not occur")
					if len(test.Expected) <= i {
						continue
					}
					assert.Equal(t, test.Expected[i], ans, "did not get the expected result")
				}
			})
		}
	})

}

func TestAddGlobalCommands(t *testing.T) {

	testCases := []struct {
		Name        string
		CommandArgs []string
		Arguments   []ArgumentDummy
		Expected    []interface{}
	}{
		{
			Name: "string argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: ""},
			},
			CommandArgs: []string{"mock", "-c=helloworld"},
			Expected:    []interface{}{"helloworld"},
		},
		{
			Name: "string slice argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: []string{""}},
			},
			CommandArgs: []string{"mock", "-c=helloworld"},
			Expected:    []interface{}{[]string{"helloworld"}},
		},
		{
			Name: "int argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: 0},
			},
			CommandArgs: []string{"mock", "-c=4"},
			Expected:    []interface{}{4},
		},

		{
			Name: "uint argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: uint(0)},
			},
			CommandArgs: []string{"mock", "-c=4"},
			Expected:    []interface{}{uint(4)},
		},
		{
			Name: "bool argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: false},
			},
			CommandArgs: []string{"mock", "--config"},
			Expected:    []interface{}{true},
		},
		{
			Name: "struct argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: struct{}{}},
			},
			CommandArgs: []string{"mock", "-ccc"},
			Expected:    []interface{}{3},
		},
		{
			Name: "multiple argument should work",
			Arguments: []ArgumentDummy{
				{LongName: "config", ConfigID: "config", ShortName: "c", Message: "some config", Default: false},
				{LongName: "verbose", ConfigID: "verbose", ShortName: "v", Message: "some verbose", Default: struct{}{}},
			},
			CommandArgs: []string{"mock", "-vvv", "-c"},
			Expected:    []interface{}{true, 3},
		},
	}
	t.Run("Global Arguments should be set and pass", func(t *testing.T) {
		for _, test := range testCases {
			t.Run(test.Name, func(t *testing.T) {
				var root = &cobra.Command{
					Use:   "application_test",
					Short: "my application is doing xyz",
					Long:  "my application is achieving xyz by doing abc and then def",
					Run: func(cmd *cobra.Command, args []string) {
						fmt.Println("Ran successfully")
					},
				}
				c := New(logrus.New(), root)
				ver := GetMockCommand("df")
				c.AddGlobalArguments(ArgumentsArr(test.Arguments))
				c.AddCommandAndArguments(&ver, ArgumentsArr([]ArgumentDummy{}))
				root.SetArgs(test.CommandArgs)
				root.Execute()
				for i, arg := range test.Arguments {
					var err error
					var ans interface{}
					switch arg.Default.(type) {
					case []string:
						ans, err = ver.Flags().GetStringSlice(arg.LongName)
					case uint:
						ans, err = ver.Flags().GetUint(arg.LongName)
					case int:
						ans, err = ver.Flags().GetInt(arg.LongName)
					case string:
						ans, err = ver.Flags().GetString(arg.LongName)
					case bool:
						ans, err = ver.Flags().GetBool(arg.LongName)
					case struct{}:
						ans, err = ver.Flags().GetCount(arg.LongName)
					}
					assert.NoError(t, err, "error should not occur")
					assert.Equal(t, test.Expected[i], ans, "did not get the expected result")
				}
			})
		}
	})
}
