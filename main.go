package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bdun1013/helm-snapshot/pkg/printer"
	"github.com/bdun1013/helm-snapshot/pkg/runner"
	"github.com/spf13/cobra"
)

var version string

func main() {
	Execute(version)
}

var testConfig = runner.TestConfig{}

var cmd = &cobra.Command{
	Use:   "snapshot [flags] CHART [...]",
	Short: "snapshot for helm charts",
	Long: `Running chart snapshot written in YAML.

This renders your charts locally (without tiller) and
validates the rendered output with the tests defined in
test suite files. Simplest test suite file looks like
below:

---
# CHART_PATH/tests/deployment_test.yaml
suite: test my deployment
templates:
  - deployment.yaml
tests:
  - it: should be a Deployment
    asserts:
      - isKind:
          of: Deployment
---

Put the test files in "tests" directory under your chart
with suffix "_test.yaml", and run:

$ helm snapshot my-chart

Or specify the suite files glob path pattern:

$ helm snapshot -f 'my-tests/*.yaml' my-chart

Check https://github.com/bdun1013/helm-snapshot for more
details about how to write tests.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, chartPaths []string) {
		var colored *bool
		if cmd.PersistentFlags().Changed("color") {
			colored = &testConfig.Colored
		}
		printer := printer.NewPrinter(os.Stdout, colored)
		runner := runner.TestRunner{Printer: printer, Config: testConfig}
		passed := runner.Run(chartPaths)

		if !passed {
			os.Exit(1)
		}
	},
}

// Execute execute snapshot command
func Execute(version string) {
	cmd.AddCommand(newVersionCommand(os.Stdout, version))
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newVersionCommand(out io.Writer, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version number of helm snapshot",
		Long:  "Prints the version number of helm snapshot",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(out, "helm snapshot plugin version", version, "built with go version", runtime.Version())
		},
	}

	return cmd
}

func init() {
	cmd.PersistentFlags().BoolVar(
		&testConfig.Colored, "color", false,
		"enforce printing colored output even stdout is not a tty. Set to false to disable color",
	)

	defaultFilePattern := filepath.Join("templates/tests/snapshot", "*_test.yaml")
	cmd.PersistentFlags().StringArrayVarP(
		&testConfig.TestFiles, "file", "f", []string{defaultFilePattern},
		"glob paths of test files location, default to "+defaultFilePattern,
	)

	cmd.PersistentFlags().BoolVarP(
		&testConfig.UpdateSnapshot, "update-snapshot", "u", false,
		"update the snapshot cached if needed, make sure you review the change before update",
	)

	cmd.PersistentFlags().BoolVarP(
		&testConfig.WithSubChart, "with-subchart", "s", true,
		"include tests of the subcharts within `charts` folder",
	)
}
