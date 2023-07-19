package basecli

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/bom-squad/go-cli/internal/basecli/config"

	"github.com/sirupsen/logrus"
)

// DefaultBaseConfig the default configuration of the global command line argument
type BaseApplication = config.BaseApplication

var DefaultBaseApplication config.BaseApplication = config.BaseApplication{}

var basicArgs = []Argument{
	{ConfigID: "config_path", LongName: "config", ShortName: "c", Message: "Configuration file path", Default: DefaultBaseApplication.ConfigPath},
}

var loggerArgs = []Argument{
	{ConfigID: "logger.level", LongName: "level", ShortName: "D", Message: fmt.Sprintf("Log depth level, options=%v", logrus.AllLevels), Default: ""},
	{ConfigID: "logger.verbose", LongName: "verbose", ShortName: "v", Message: "Log verbosity level [-v,--verbose=1] = info, [-vv,--verbose=2] = debug, [-vvv,--verbose=3] = trace", Default: struct{}{}, IsCliOnly: false},
	{ConfigID: "logger.quiet", LongName: "quiet", ShortName: "q", Message: "Suppress all logging output", Default: DefaultBaseApplication.Logger.Quiet},
	{ConfigID: "logger.structured", LongName: "structured", ShortName: "", Message: "Enable structured logger", Default: false, IsHidden: false},
}

type ArgOptions int

const (
	ARG_ALL          ArgOptions = 0b11111
	ARG_ALL_NO_CACHE ArgOptions = 0b01111
	ARG_BASIC        ArgOptions = 0b00001
	LOG_ARGS         ArgOptions = 0b01000
	CACHE_ARGS       ArgOptions = 0b10000
)

func getBasicGlobalArguments(applicationName string, options ArgOptions) []Argument {
	var DEFAULT_OUTPUT_DIRECTORY = filepath.Join(xdg.CacheHome, applicationName)
	args := make([]Argument, 0)
	if options&ARG_BASIC > 0 {
		args = append(args, basicArgs...)
	}
	if options&LOG_ARGS > 0 {
		args = append(args, loggerArgs...)
	}
	if options&CACHE_ARGS > 0 {
		args = append(args, Argument{ConfigID: "output_directory", LongName: "output-directory", ShortName: "d", Message: "Output directory path", Default: DEFAULT_OUTPUT_DIRECTORY})
	}
	return args
}
