package views

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/controllers"
)

func ShowChoseDirectory(pages *tview.Pages) *tview.Frame {
	questionInfo := config.STATE.Question
	form := tview.NewForm()
	form.AddTextArea("Enter the path to the directory you'd like to create the template...", ".", 0, 1, 0, nil)
	form.AddTextView("Question Title", fmt.Sprintf("%s. %s", questionInfo.Id, questionInfo.QuestionTitle), 0, 1, true, false)
	form.AddTextView("Difficulty", questionInfo.Difficulty, 0, 1, true, false)
	form.AddTextView("Topics", questionInfo.Topics, 0, 1, true, false)
	form.AddTextView("Link", questionInfo.QuestionLink, 0, 1, true, false)

	form.AddButton("CREATE", func() {
			controllers.CreateQuestionProject(form.GetFormItem(0).(*tview.TextArea).GetText(), &questionInfo)
			config.APP.Stop()
		})
	form.AddButton("GO BACK", func() {
			pages.SwitchToPage("pick_lang")
		})

	frame := AddFrameWrapper(form, "Choose directory location...", "This will create a folder in this directory with the name of the leetcode problem (NOTE: use `tab` to navigate)", "")
	return frame
}
