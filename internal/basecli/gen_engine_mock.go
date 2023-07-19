//go:build !dev
// +build !dev

package basecli

import "github.com/spf13/cobra"

func (b *Engine) commandAndArgsToActionArgs(cmd *cobra.Command, arguments Arguments) error {
	return nil
}

func (b *Engine) addGlobalToConfig(arguments Arguments) {
}

func (b *Engine) addGenerateCommands() {
}
