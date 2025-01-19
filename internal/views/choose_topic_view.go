package views

import (
	"strings"

	"github.com/rivo/tview"
	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/constants"
	"github.com/vincent-buchner/leetcode-framer/internal/controllers"
	stateservice "github.com/vincent-buchner/leetcode-framer/internal/services/state_service"
	"github.com/vincent-buchner/leetcode-framer/pkg/utils"
)

func ShowChooseTopic(pages *tview.Pages) *tview.Frame {
	form := tview.NewForm()
	keys := func(m map[string]int) []string {
		k := make([]string, 0, len(m))
		for key := range m {
			k = append(k, key)
		}
		return k
	}(constants.TOPICS)

	form.AddTextArea("Enter Topic: ", "", 0, 1, 0, func(text string) {

		if form.GetFormItemCount() < 2 {
			form.AddTextView("Results: ", utils.FuzzySearch(text, keys, 15), 0, 0, true, false)
		} else {
			form.RemoveFormItem(1)
			form.AddTextView("Results: ", utils.FuzzySearch(text, keys, 15), 0, 0, true, false)
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
		narrowed_topics := strings.Split(result_text, ",")

		// Trim the white space
		queryTopic := strings.TrimSpace(narrowed_topics[0])

		// Lower and Replace space with `-`
		slugQueryTopic := strings.ToLower(strings.Replace(queryTopic, " ", "-", -1))

		// Get random topic question
		question, err := controllers.GetRandomTopicQuestion(slugQueryTopic, constants.TOPICS[narrowed_topics[0]])
		if err != nil {
			pages.AddAndSwitchToPage("error", ShowError(pages, err), true)
		} else {
			config.STATE.Question = question
			stateservice.StateToJSON(config.STATE, "./configs/config.json")
			pages.AddAndSwitchToPage("pick_lang", ShowSupportedLanguages(pages), true)
		}
	})
	form.AddButton("BACK", func() {
		pages.SwitchToPage("menu")
	})

	frame := AddFrameWrapper(form, "Enter the name of the topic you'd like a random question from: (use `TAB` to navigate)", "(i.e Array, Stack, Backtracking)", "")
	return frame

}