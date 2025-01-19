package views

import (
	"github.com/rivo/tview"
	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/controllers"
	stateservice "github.com/vincent-buchner/leetcode-framer/internal/services/state_service"
)

func ShowSpecificQuestion(pages *tview.Pages) *tview.Frame {
	form := tview.NewForm()

	form.AddTextArea("Enter Question ID: ", "", 0, 1, 0, nil)
	form.AddButton("SEARCH", func() {
		id_text := form.GetFormItem(0).(*tview.TextArea)
		question, err := controllers.GetSpecificQuestion(id_text.GetText())
		if err != nil {
			pages.AddAndSwitchToPage("error", ShowError(pages, err), true)
		} else {
			config.STATE.Question = question
			stateservice.StateToJSON(config.STATE, "./configs/config.json")
			pages.SwitchToPage("pick_lang")
		}
	})
	form.AddButton("BACK", func() {
		pages.SwitchToPage("menu")
	})

	frame := AddFrameWrapper(form, "Enter the question ID of the problem you'd like to solve! (use `TAB` to navigate)", "(Ex. 'Two Sum' is problem number 1)", "")
	return frame
}