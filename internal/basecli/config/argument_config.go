package config

import (
	"github.com/spf13/pflag"
)

type Argument interface {
	GetConfigID() string
	GetLongName() string
}

func (c Engine) BindFlags(flag *pflag.FlagSet, args []Argument) error {
	for _, arg := range args {
		if err := c.BindPFlag(arg.GetConfigID(), flag.Lookup(arg.GetLongName())); err != nil {
			return err
		}
	}
	return nil
}
