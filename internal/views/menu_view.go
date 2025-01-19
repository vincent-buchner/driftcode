package views

import (
	"github.com/rivo/tview"
	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/controllers"
	stateservice "github.com/vincent-buchner/leetcode-framer/internal/services/state_service"
)

func ShowMainMenu(pages *tview.Pages) *tview.Frame  {

	list := tview.NewList()
	list.AddItem("Daily Question", "Start with the daily question presented by Leetcode.", 'a', func () {
		question, err := controllers.GetDailyQuestion()
		if err != nil {
			pages.AddAndSwitchToPage("error", ShowError(pages, err), true)
		} else {
			config.STATE.Question = question
			stateservice.StateToJSON(config.STATE, "./configs/config.json")
			pages.SwitchToPage("pick_lang")
		}
	})
	list.AddItem("Select Question", "Select a specific question via the ID of the question.", 'b', func ()  {
		pages.SwitchToPage("choose_problem")
	})
	list.AddItem("Random Question Within Difficulty", "Generates a random question within the specified difficulty.", 'c', func () {
		pages.SwitchToPage("rand_difficulty")
	})
	list.AddItem("Random Question Within Company", "Generates a random question asked by specific companies.", 'd', func() {
		pages.SwitchToPage("choose_company")
	})
	list.AddItem("Random Question Within Topic", "Generates a random question within a given topic.", 'e', func() {
		pages.SwitchToPage("choose_topic")
	})

	framer := AddFrameWrapper(list, "Welcome to LeetCode Framer CLI and TUI!", "A program designed to practice Leetcode questions within your favorite editor", "`Ctrl + c` to exit program...")
	
	return framer
}