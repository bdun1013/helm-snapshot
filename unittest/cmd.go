package unittest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

// TestConfig stores config setup by user in command line
type TestConfig struct {
	Colored        bool
	UpdateSnapshot bool
	WithSubChart   bool
	TestFiles      []string
}

var testConfig = TestConfig{}

var cmd = &cobra.Command{
	Use:   "unittest [flags] CHART [...]",
	Short: "unittest for helm charts",
	Long: `Running chart unittest written in YAML.

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

$ helm unittest my-chart

Or specify the suite files glob path pattern:

$ helm unittest -f 'my-tests/*.yaml' my-chart

Check https://github.com/bpdunni/helm-unittest for more
details about how to write tests.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, chartPaths []string) {
		var colored *bool
		if cmd.PersistentFlags().Changed("color") {
			colored = &testConfig.Colored
		}
		printer := NewPrinter(os.Stdout, colored)
		runner := TestRunner{Printer: printer, Config: testConfig}
		passed := runner.Run(chartPaths)

		if !passed {
			os.Exit(1)
		}
	},
}

// Execute execute unittest command
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
		Short: "Prints the version number of helm unittest",
		Long:  "Prints the version number of helm unittest",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(out, "helm unittest plugin version", version, "built with go version", runtime.Version())
		},
	}

	return cmd
}

func init() {
	cmd.PersistentFlags().BoolVar(
		&testConfig.Colored, "color", false,
		"enforce printing colored output even stdout is not a tty. Set to false to disable color",
	)

	defaultFilePattern := filepath.Join("templates/tests/unit", "*_test.yaml")
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
