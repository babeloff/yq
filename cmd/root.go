package cmd

import (
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/thediveo/enumflag"
	"os"

	"github.com/spf13/cobra"
	logging "gopkg.in/op/go-logging.v1"
)

func New() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "yq",
		Short: "yq is a lightweight and portable command-line YAML processor.",
		Long:  `yq is a lightweight and portable command-line YAML processor. It aims to be the jq or sed of yaml files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				cmd.Print(GetVersionDisplay())
				return nil
			}
			cmd.Println(cmd.UsageString())
			return nil

		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetOut(cmd.OutOrStdout())
			var format = logging.MustStringFormatter(
				`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
			)
			var backend = logging.AddModuleLevel(
				logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))

			if verbose {
				backend.SetLevel(logging.DEBUG, "")
			} else {
				backend.SetLevel(logging.ERROR, "")
			}

			logging.SetBackend(backend)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")

	rootCmd.PersistentFlags().VarP(
		enumflag.New(&inputType, "input-type", yqlib.InputModeIds, enumflag.EnumCaseInsensitive),
		"from-type", "f",
		"type of the input; can be 'yaml', 'json', or 'props' [default: yaml]")
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&outputType, "output-type", yqlib.OutputModeIds, enumflag.EnumCaseInsensitive),
		"to-type", "t",
		"type of the output; can be 'yaml', 'json', or 'props' [default: yaml]")

	rootCmd.PersistentFlags().BoolVarP(&nullInput, "null-input", "n", false, "Don't read input, simply evaluate the expression given. Useful for creating yaml docs from scratch.")
	rootCmd.PersistentFlags().BoolVarP(&noDocSeparators, "no-doc", "N", false, "Don't print document separators (---)")

	rootCmd.PersistentFlags().IntVarP(&indent, "indent", "I", 2, "sets indent level for output")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information and quit")
	rootCmd.PersistentFlags().BoolVarP(&writeInplace, "inplace", "i", false, "update the yaml file inplace of first yaml file given.")
	rootCmd.PersistentFlags().BoolVarP(&unwrapScalar, "unwrapScalar", "", true, "unwrap scalar, print the value with no quotes, colors or comments")
	rootCmd.PersistentFlags().BoolVarP(&prettyPrint, "prettyPrint", "P", false, "pretty print, shorthand for '... style = \"\"'")
	rootCmd.PersistentFlags().BoolVarP(&exitStatus, "exit-status", "e", false, "set exit status if there are no matches or null or false is returned")

	rootCmd.PersistentFlags().BoolVarP(&forceColor, "colors", "C", false, "force print with colors")
	rootCmd.PersistentFlags().BoolVarP(&forceNoColor, "no-colors", "M", false, "force print with no colors")
	rootCmd.AddCommand(
		createEvaluateSequenceCommand(),
		createEvaluateAllCommand(),
		completionCmd,
	)
	return rootCmd
}
