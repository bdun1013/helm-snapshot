package suite

import (
	"github.com/bpdunni/helm-snapshot/pkg/assertion"
	"github.com/bpdunni/helm-snapshot/pkg/printer"
)

// TestJobResult result return by TestJob.Run
type TestJobResult struct {
	DisplayName   string
	Index         int
	Passed        bool
	ExecError     error
	AssertsResult []*assertion.AssertionResult
}

func (tjr TestJobResult) print(printer *printer.Printer, verbosity int) {
	if tjr.Passed {
		return
	}

	if tjr.ExecError != nil {
		printer.Println(printer.Highlight("- "+tjr.DisplayName), 1)
		printer.Println(
			printer.Highlight("Error: ")+
				tjr.ExecError.Error()+"\n",
			2,
		)
		return
	}

	printer.Println(printer.Danger("- "+tjr.DisplayName+"\n"), 1)
	for _, assertResult := range tjr.AssertsResult {
		assertResult.Print(printer, verbosity)
	}
}
