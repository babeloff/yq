package cmd

import (
	"github.com/thediveo/enumflag"
)

var unwrapScalar = true

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

var inputType = FromYaml
var outputType = ToYaml

var writeInplace = false

var exitStatus = false
var forceColor = false
var forceNoColor = false
var colorsEnabled = false
var indent = 2
var noDocSeparators = false
var nullInput = false
var verbose = false
var version = false
var prettyPrint = false

var completedSuccessfully = false
