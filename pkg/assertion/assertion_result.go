package assertion

import (
	"fmt"

	"github.com/bdun1013/helm-snapshot/pkg/printer"
)

// AssertionResult result return by Assertion.Assert
type AssertionResult struct {
	Index      int
	FailInfo   []string
	Passed     bool
	AssertType string
	Not        bool
	CustomInfo string
}

func (ar AssertionResult) Print(printer *printer.Printer, verbosity int) {
	if ar.Passed {
		return
	}
	var title string
	if ar.CustomInfo != "" {
		title = ar.CustomInfo
	} else {
		var notAnnotation string
		if ar.Not {
			notAnnotation = " NOT"
		}
		title = fmt.Sprintf("- asserts[%d]%s `%s` fail", ar.Index, notAnnotation, ar.AssertType)
	}
	printer.Println(printer.Danger(title+"\n"), 2)
	for _, infoLine := range ar.FailInfo {
		printer.Println(infoLine, 3)
	}
	printer.Println("", 0)
}
