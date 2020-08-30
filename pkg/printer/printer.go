package printer

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// NewPrinter create a Printer with Writer to print and colored config
func NewPrinter(writer io.Writer, colored *bool) *Printer {
	p := &Printer{
		Writer:  writer,
		Colored: colored,
	}

	p.colors.Success = color.New(color.FgGreen)
	p.SetupColor(p.colors.Success)
	p.colors.SuccessBg = color.New(color.BgGreen, color.FgBlack)
	p.SetupColor(p.colors.SuccessBg)
	p.colors.Danger = color.New(color.FgRed)
	p.SetupColor(p.colors.Danger)
	p.colors.DangerBg = color.New(color.BgRed, color.FgWhite)
	p.SetupColor(p.colors.DangerBg)
	p.colors.Warning = color.New(color.FgYellow)
	p.SetupColor(p.colors.Warning)
	p.colors.WarningBg = color.New(color.BgYellow, color.FgBlack)
	p.SetupColor(p.colors.WarningBg)
	p.colors.Highlight = color.New(color.Bold)
	p.SetupColor(p.colors.Highlight)
	p.colors.Faint = color.New(color.Faint)
	p.SetupColor(p.colors.Faint)

	return p
}

// Printer simple printing implement
type Printer struct {
	Writer  io.Writer
	Colored *bool
	colors  struct {
		Success   *color.Color
		SuccessBg *color.Color
		Warning   *color.Color
		WarningBg *color.Color
		Danger    *color.Color
		DangerBg  *color.Color
		Highlight *color.Color
		Faint     *color.Color
	}
}

func (p *Printer) SetupColor(color *color.Color) {
	if p.Colored != nil {
		if *p.Colored {
			color.EnableColor()
		} else {
			color.DisableColor()
		}
	}
}

func (p *Printer) Println(content string, indentLevel int) {
	var indent string
	for i := 0; i < indentLevel; i++ {
		indent += "\t"
	}
	fmt.Fprintln(p.Writer, indent+content)
}

func (p *Printer) Success(format string, a ...interface{}) string {
	return p.colors.Success.Sprintf(format, a...)
}

func (p *Printer) SuccessLabel(format string, a ...interface{}) string {
	return p.colors.SuccessBg.Sprintf(format, a...)
}

func (p *Printer) Danger(format string, a ...interface{}) string {
	return p.colors.Danger.Sprintf(format, a...)
}

func (p *Printer) DangerLabel(format string, a ...interface{}) string {
	return p.colors.DangerBg.Sprintf(format, a...)
}

func (p *Printer) Warning(format string, a ...interface{}) string {
	return p.colors.Warning.Sprintf(format, a...)
}

func (p *Printer) WarningLabel(format string, a ...interface{}) string {
	return p.colors.WarningBg.Sprintf(format, a...)
}

func (p *Printer) Highlight(format string, a ...interface{}) string {
	return p.colors.Highlight.Sprintf(format, a...)
}

func (p *Printer) Faint(format string, a ...interface{}) string {
	return p.colors.Faint.Sprintf(format, a...)
}
