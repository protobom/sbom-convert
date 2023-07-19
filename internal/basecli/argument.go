package basecli

import "github.com/bom-squad/go-cli/internal/basecli/cmd"

// Argument Contains the configuration of command line argument
type Argument struct {
	// ConfigID this is used by config to set the value for the argument
	ConfigID string
	// ShortName this is the shorthand name of the command argument
	// Note: it should be a single charachter
	ShortName string

	// LongName this is the long name of the command argument.
	LongName string

	// Message this is the help message for the user
	Message string

	// Default this is the default value of the command if the command is not set by the user.
	// Note the type of the default is used to set the type of argument so if default is string the argument will be a string type and so on.
	// Only primitive types and slices are supported.
	Default interface{}

	// IsHidden this boolean is used if we want to hide the command from the help menu of the cobra
	IsHidden bool

	// IsCliOnly this boolean restricts the commandline argument from being added to the config. if set to true the command line argument will not be included in the config
	IsCliOnly bool
}

// GetConfigID returns the config id of the argument
func (a Argument) GetConfigID() string {
	return a.ConfigID
}

// GetLongName returns the long name of the command argument
func (a Argument) GetLongName() string {
	return a.LongName
}

// GetShortName returns the short name of the command argument
func (a Argument) GetShortName() string {
	return a.ShortName
}

// GetMessage returns the help message for the user
func (a Argument) GetMessage() string {
	return a.Message
}

// GetDefault returns the default value of the command line argument
func (a Argument) GetDefault() interface{} {
	return a.Default
}

// GetIsHidden returns if the cmd argument should be hidden from user
func (a Argument) GetIsHidden() bool {
	return a.IsHidden
}

// GetIsCliOnly returns if the argument should be cli only
func (a Argument) GetIsCliOnly() bool {
	return a.IsCliOnly
}

// Arguments this is array of argument
type Arguments []Argument

// Get returns the i th Argument
func (a Arguments) Get(i int) cmd.Argument {
	return a[i]
}

// Len returns the length of the argument
func (a Arguments) Len() int {
	return len(a)
}
