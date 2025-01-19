package views

import (
	"github.com/rivo/tview"
	"github.com/vincent-buchner/leetcode-framer/internal/config"
	stateservice "github.com/vincent-buchner/leetcode-framer/internal/services/state_service"
)

func ShowSupportedLanguages(pages *tview.Pages) *tview.Frame {
	list := tview.NewList()

	list.AddItem("Python", "", 'a', func() {
		config.STATE.Language = "python"
		stateservice.StateToJSON(config.STATE, "./configs/config.json")
		pages.AddAndSwitchToPage("pick_dir", ShowChoseDirectory(pages), true)
	}).
	AddItem("Java", "", 'b', func() {
		config.STATE.Language = "java"
		stateservice.StateToJSON(config.STATE, "./configs/config.json")
		pages.AddAndSwitchToPage("pick_dir", ShowChoseDirectory(pages), true)
	}).
	AddItem("Golang", "", 'c', func() {
		config.STATE.Language = "golang"
		stateservice.StateToJSON(config.STATE, "./configs/config.json")
		pages.AddAndSwitchToPage("pick_dir", ShowChoseDirectory(pages), true)
	}).
	AddItem("Back to main menu...", "", 'd', func() {
		pages.SwitchToPage("menu")
	})
	frame := AddFrameWrapper(list, "Pick your preferred language...", "(More coming soon!)", "")
	return frame
}