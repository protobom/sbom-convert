package config

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type Arg struct {
	ConfigID  string
	ShortName string
	LongName  string
	Message   string
	Default   interface{}
}

func (a Arg) GetConfigID() string {
	return a.ConfigID
}

func (a Arg) GetLongName() string {
	return a.LongName
}

var root = &cobra.Command{
	Use:   "application_test",
	Short: "my application is doing xyz",
	Long:  "my application is achieving xyz by doing abc and then def",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Ran successfully")
	},
}

func TestArgumentConfig(t *testing.T) {
	conf := New(logrus.New(), "sd")
	args := []Arg{
		{
			ConfigID: "someID",
			LongName: "some",
			Default:  "",
		},
		{
			ConfigID: "someID3",
			LongName: "some3",
			Default:  "",
		},
	}
	argI := make([]Argument, len(args))
	for i, arg := range args {
		root.Flags().String(arg.LongName, arg.Default.(string), "some id config")
		argI[i] = arg
	}
	err := conf.BindFlags(root.Flags(), argI)
	assert.NoError(t, err, "Should not return error")
}
