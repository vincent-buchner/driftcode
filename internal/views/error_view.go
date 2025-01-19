package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	// "github.com/vincent-buchner/leetcode-framer/cmd/config"
)

func ShowError(pages *tview.Pages, err error) *tview.Frame {
	form := tview.NewForm()
	form.AddTextView("", fmt.Sprintf(`
	.-""""""-.
  .'          '.
 /   O      O   \
:           ''   :
|                |
:    .------.    :
 \  '        '  /
  '.          .'
    '-......-'
ERROR: %s
`, err), 0, 11, true, false)


	form.SetFieldTextColor(tcell.ColorDarkRed)
	form.AddButton("MAIN MENU", func() {
		pages.SwitchToPage("menu")
	})

	frame := AddFrameWrapper(form, "Oh no! There was an error within the program!", "", "")
	return frame

}
