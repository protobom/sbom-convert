package cmd

type Argument interface {
	GetConfigID() string
	GetLongName() string
	GetShortName() string
	GetMessage() string
	GetIsHidden() bool
	GetDefault() interface{}
}

type Arguments interface {
	Get(i int) Argument
	Len() int
}
