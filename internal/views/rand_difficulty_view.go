package views

import (
	"github.com/rivo/tview"
	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/controllers"
	stateservice "github.com/vincent-buchner/leetcode-framer/internal/services/state_service"
)

func ShowRandomDifficulty(pages *tview.Pages) *tview.Frame {

	list := tview.NewList()
	list.AddItem("Easy", "", 'a', func() {
		question, err := controllers.GetRandomDifficulty("Easy")

		// TODO: Modularize this, make it DRY
		if err != nil {
			pages.AddAndSwitchToPage("error", ShowError(pages, err), true)
		} else {
			config.STATE.Question = question
			stateservice.StateToJSON(config.STATE, "./configs/config.json")
			pages.SwitchToPage("pick_lang")
		}
	})
	list.AddItem("Medium", "", 'b', func() {

		question, err := controllers.GetRandomDifficulty("Medium")
		if err != nil {
			pages.AddAndSwitchToPage("error", ShowError(pages, err), true)
		} else {
			config.STATE.Question = question
			stateservice.StateToJSON(config.STATE, "./configs/config.json")
			pages.SwitchToPage("pick_lang")
		}
	})
	list.AddItem("Hard", "", 'c', func() {

		question, err := controllers.GetRandomDifficulty("Hard")
		if err != nil {
			pages.AddAndSwitchToPage("error", ShowError(pages, err), true)
		} else {
			config.STATE.Question = question
			stateservice.StateToJSON(config.STATE, "./configs/config.json")
			pages.SwitchToPage("pick_lang")
		}
	})
	list.AddItem("Back to Main Menu", "", 'd', func ()  {
		pages.SwitchToPage("menu")
	})

	frame := AddFrameWrapper(list, "Choose a difficulty", "This will pick a random problem within the given difficulty.", "")
	return frame
}