package yqlib

import "github.com/thediveo/enumflag"

type InputTypeEnum enumflag.Flag

const (
	FromYaml InputTypeEnum = iota
	FromJson
	FromProps
)

var InputModeIds = map[InputTypeEnum][]string{
	FromYaml:  {"from-yaml"},
	FromJson:  {"from-json"},
	FromProps: {"from-props"},
}

type OutputTypeEnum enumflag.Flag

const (
	ToYaml OutputTypeEnum = iota
	ToJson
	ToProps
)

var OutputModeIds = map[OutputTypeEnum][]string{
	ToYaml:  {"to-yaml"},
	ToJson:  {"to-json"},
	ToProps: {"to-props"},
}
