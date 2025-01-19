package config

import (
	"time"

	"github.com/rivo/tview"
	"github.com/vincent-buchner/leetcode-framer/internal/models"
	"github.com/vincent-buchner/leetcode-framer/internal/services/state_service"
)

var STATE ApplicationState
var APP *tview.Application

type ApplicationState struct {
	Language      string               `json:"lang"`
	Username      string               `json:"username"`
	Question      models.QuestionModel `json:"question"`
	DailyQuestion models.QuestionModel `json:"dailyQuestion"`
	TodayDate     time.Time            `json:"todayDate"`
}

func Init() (*tview.Application, *ApplicationState, error) {

	// Create App
	APP = tview.NewApplication().EnableMouse(true)

	// Try to load the state from file
	stateservice.JSONToState(&STATE, "./config.json")
	return APP, &STATE, nil
}

