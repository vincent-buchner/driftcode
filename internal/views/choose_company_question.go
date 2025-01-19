package views

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/controllers"
	"github.com/vincent-buchner/leetcode-framer/internal/models"
	stateservice "github.com/vincent-buchner/leetcode-framer/internal/services/state_service"
)

func ShowCompanyQuestion(pages *tview.Pages, companyName string, questions *[]models.QuestionDB) *tview.Frame {
	list := tview.NewList()

	for i, question := range *questions {
		list.AddItem(fmt.Sprintf("%d. %s", question.ID, question.Name), fmt.Sprintf("%s, %s", question.Difficulty, question.Link), rune(i + 1), func() {
			model_question, err := controllers.GetSpecificQuestion(strconv.Itoa(question.ID))
			if err != nil {
				pages.AddAndSwitchToPage("error", ShowError(pages, err), true)
			} else {
				config.STATE.Question = model_question
				stateservice.StateToJSON(config.STATE, "./configs/config.json")
				pages.AddAndSwitchToPage("pick_lang", ShowSupportedLanguages(pages), true)
			}
		})
	}

	frame := AddFrameWrapper(list, "Choose question: ", fmt.Sprintf("Company Selected: %s", companyName), "")
	return frame
}
