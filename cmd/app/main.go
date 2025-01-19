package main

import (
	"github.com/rivo/tview"
	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/db"
	"github.com/vincent-buchner/leetcode-framer/internal/views"
	"github.com/joho/godotenv"
)


func main() {

	// Start app
	app, _, err := config.Init()
	if err != nil {
		panic(err)
	}

	// Connect db
	db.InitDatabase("data/leetcode.db", false)

	// Load ENV
	err = godotenv.Load()
    if err != nil {
		panic(err)
    }
	pages := tview.NewPages()

	// Static Pages, no dynamic content
	pages.AddPage("menu", views.ShowMainMenu(pages), true, true)
	pages.AddPage("pick_lang", views.ShowSupportedLanguages(pages), true, false)
	pages.AddPage("choose_problem", views.ShowSpecificQuestion(pages), true, false)
	pages.AddPage("rand_difficulty", views.ShowRandomDifficulty(pages), true, false)
	pages.AddPage("choose_company", views.ShowChooseCompany(pages), true, false)
	pages.AddPage("choose_topic", views.ShowChooseTopic(pages), true, false)

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}

}
