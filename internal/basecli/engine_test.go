package basecli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type DummyConfig struct {
	BaseApplication `mapstructure:",squash"`
	TestURL         string `mapstructure:"test-url"`
	TestUser        string `mapstructure:"test-user"`
	mockConfigPath  string
}

func (d DummyConfig) GetConfigPath() string {
	return d.mockConfigPath
}

func TestApplicationConfig(t *testing.T) {
	var (
		root = &cobra.Command{
			Use:   "application_test",
			Short: "my application is doing xyz",
			Long:  "my application is achieving xyz by doing abc and then def",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Ran successfully")
			},
		}
		cli = New(logrus.New(), "application_test", root)
	)
	testCases := []struct {
		Name   string
		Config string
	}{
		{
			Name: "Check if config works",
			Config: `


			`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			cli.config.ReadConfig(strings.NewReader(testCase.Config))
		})
	}

}

func TestApplicationCLI(t *testing.T) {
	var (
		root = &cobra.Command{
			Use:   "application_test",
			Short: "my application is doing xyz",
			Long:  "my application is achieving xyz by doing abc and then def",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Ran successfully")
			},
		}

		applicationPrint = &cobra.Command{
			Use:   "print",
			Short: "prints everything",
			Long:  "will print whatever you say",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(args)
				fmt.Println(cmd.Flags().GetString("argument"))
			},
		}
		cli = New(logrus.New(), "application_test", root)
	)
	cli.AddBasicCommandsAndArguments(ARG_ALL)
	cli.AddCommandAndArguments(applicationPrint, []Argument{
		{
			ConfigID:  "output",
			ShortName: "a",
			LongName:  "argument",
			Message:   "some basic string arg",
			Default:   "hello world",
		},
	})

	cli.cmd.GetCommand().SetArgs([]string{"print", "hello world", "--argument=fasas", "-vv"})
	cobra.OnInitialize(func() {
		conf := &DummyConfig{}
		cli.LoadApplication(conf)
		assert.Equal(t, "https://dummy.com", conf.TestURL, "configuration was set")
	})

	cli.cmd.GetCommand().Execute()
	t.Run("config should be set correctly", func(t *testing.T) {
		assert.Equal(t, "fasas", cli.config.GetString("output"), "output key is not set")
	})

}
