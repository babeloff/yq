package cmd

import "github.com/mikefarah/yq/v4/pkg/yqlib"

var unwrapScalar = true

var inputType = yqlib.FromYaml
var outputType = yqlib.ToYaml

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
