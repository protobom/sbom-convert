package readme

import (
	"testing"

	"github.com/bom-squad/go-cli/internal/basecli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestReadmeGetPreferredFileName(t *testing.T) {
	testcases := []struct {
		name             string
		cmd              *cobra.Command
		expectedFileName string
	}{
		{
			name: "normal command should be without underscore",
			cmd: &cobra.Command{
				Use: "mockapp",
			},
			expectedFileName: "mockapp.md",
		},
		{
			name: "subcommand readme names recommendation should be correct",
			cmd: &cobra.Command{
				Use: "mock [target]",
			},
			expectedFileName: "mockapp_mock.md",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			rg := New("mockapp", &basecli.BaseApplication{})
			rg.SetCommand(test.cmd)
			assert.Equal(t, test.expectedFileName, rg.GetPreferredFileName(), "preferred name should match")
		})
	}
}
