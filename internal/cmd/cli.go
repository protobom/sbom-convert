package cmd

import (
	"fmt"
	"os"

	"github.com/bom-squad/go-cli/internal/basecli"
	"github.com/bom-squad/go-cli/internal/config"

	log "github.com/bom-squad/go-cli/internal/basecli/logger"
	"github.com/spf13/cobra"
)

var (
	ApplicationName = "protobom"
	cli             *basecli.Engine
	version         = "0.0.0"
	conf            config.Application
)

func NewCliInit() *basecli.Engine {
	cli := basecli.New(nil, ApplicationName, translateCmd)

	err := cli.AddBasicCommandsAndArguments(basecli.ARG_ALL_NO_CACHE)
	if err != nil {
		panic(err)
	}
	err = cli.AddArguments(translateCmd, config.TranslateCommandArguments)
	if err != nil {
		panic(err)
	}
	return cli
}

func init() {
	cli = NewCliInit()
	cobra.OnInitialize(LoadApplication)
}

func LoadApplication() {
	err := cli.LoadApplication(&conf)
	if err != nil {
		er(err)
	}
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func Execute() {
	err := translateCmd.Execute()
	if err != nil {
		log.Errorf(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}
