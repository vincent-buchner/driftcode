package views

import (
	"strings"

	"github.com/rivo/tview"
	"github.com/vincent-buchner/leetcode-framer/internal/constants"
	"github.com/vincent-buchner/leetcode-framer/internal/controllers"
	"github.com/vincent-buchner/leetcode-framer/pkg/utils"
)

func ShowChooseCompany(pages *tview.Pages) *tview.Frame {
	form := tview.NewForm()

	form.AddTextArea("Enter Company: ", "", 0, 1, 0, func(text string) {
		if form.GetFormItemCount() < 2 {
			form.AddTextView("Results: ", utils.FuzzySearch(text, constants.COMPANIES, 15), 0, 0, true, false)
		} else {
			form.RemoveFormItem(1)
			form.AddTextView("Results: ", utils.FuzzySearch(text, constants.COMPANIES, 15), 0, 0, true, false)
		}
	})
	form.AddButton("SUBMIT", func() {

		// Have the results been populated? If not, exit
		if form.GetFormItemCount() < 2 {
			return
		}

		// We grab the text from the results because it'll be easier to query the db again
		result_text := form.GetFormItem(1).(*tview.TextView).GetText(false)
		
		// Parse the string back to list
		companies := strings.Split(result_text, ",")

		// Trim the white space
		queryCompany := strings.TrimSpace(companies[0])

		queried_companies, err := controllers.GetQuestionsFromCompany(queryCompany)
		if err != nil {
			pages.AddAndSwitchToPage("error", ShowError(pages, err), true)
		} else {
			pages.AddAndSwitchToPage("company_question", ShowCompanyQuestion(pages, queryCompany, &queried_companies), true)
		}
	})
	form.AddButton("BACK", func() {
		pages.SwitchToPage("menu")
	})

	frame := AddFrameWrapper(form, "Enter the name of the company you'd like a random question from: (use `TAB` to navigate)", "(i.e Apple, Microsoft, Google)", "")
	return frame

}