package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func AddFrameWrapper(primitive tview.Primitive, header string, subheader string, footer string) *tview.Frame {

	if footer == "" {
		footer = "`Ctrl + c` to exit program..."
	}

	frame := tview.NewFrame(primitive).
	AddText(header, true, tview.AlignLeft, tcell.Color231).
	AddText(subheader, true, tview.AlignLeft, tcell.Color101).
	AddText(footer, false, tview.AlignLeft, tcell.Color101)

	return frame
	
}