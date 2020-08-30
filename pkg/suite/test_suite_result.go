package suite

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bpdunni/helm-snapshot/pkg/printer"
	"github.com/bpdunni/helm-snapshot/pkg/snapshot"
)

// TestSuiteResult result return by TestSuite.Run
type TestSuiteResult struct {
	DisplayName      string
	FilePath         string
	Passed           bool
	ExecError        error
	TestsResult      []*TestJobResult
	SnapshotCounting struct {
		Total    uint
		Failed   uint
		Created  uint
		Vanished uint
	}
}

func (tsr TestSuiteResult) Print(printer *printer.Printer, verbosity int) {
	tsr.printTitle(printer)
	if tsr.ExecError != nil {
		printer.Println(printer.Highlight("- Execution Error: "), 1)
		printer.Println(tsr.ExecError.Error()+"\n", 2)
		return
	}

	for _, result := range tsr.TestsResult {
		result.print(printer, verbosity)
	}
}

func (tsr TestSuiteResult) printTitle(printer *printer.Printer) {
	var label string
	if tsr.Passed {
		label = printer.SuccessLabel(" PASS ")
	} else {
		label = printer.DangerLabel(" FAIL ")
	}
	var pathToPrint string
	if tsr.FilePath != "" {
		pathToPrint = printer.Faint(filepath.Dir(tsr.FilePath)+string(os.PathSeparator)) +
			filepath.Base(tsr.FilePath)
	}
	name := printer.Highlight(tsr.DisplayName)
	printer.Println(
		fmt.Sprintf("%s %s\t%s", label, name, pathToPrint),
		0,
	)
}

func (tsr *TestSuiteResult) countSnapshot(cache *snapshot.Cache) {
	tsr.SnapshotCounting.Created = cache.InsertedCount()
	tsr.SnapshotCounting.Failed = cache.FailedCount()
	tsr.SnapshotCounting.Total = cache.CurrentCount()
	tsr.SnapshotCounting.Vanished = cache.VanishedCount()
}
