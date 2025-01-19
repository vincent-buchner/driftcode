package controllers

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/db"
	"github.com/vincent-buchner/leetcode-framer/internal/models"
	stateservice "github.com/vincent-buchner/leetcode-framer/internal/services/state_service"
)

func TestGetDailyQuestion(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	t.Run("Get daily question - no cache", func(t *testing.T) {
		test_config := config.ApplicationState{
			Language: "python",
		}
		err := stateservice.StateToJSON(test_config, filepath.Join(tmpDir, "configs/config.json"))
		assert.NoError(t, err)

		question, err := GetDailyQuestion()
		assert.NoError(t, err)

		assert.NotNil(t, question)
	})

	t.Run("Get daily question - with cache", func(t *testing.T) {
		test_config := config.ApplicationState{
			Language: "python",
			DailyQuestion: models.QuestionModel{
				Id: "1",
				QuestionTitle: "Two Sum",
			},
			TodayDate: time.Now(),
		}

		err := stateservice.StateToJSON(test_config, filepath.Join(tmpDir, "configs/config.json"))
		assert.NoError(t, err)

		config.STATE = test_config

		question, err := GetDailyQuestion()
		assert.NoError(t, err)

		assert.NotNil(t, question)
		assert.Equal(t, "Two Sum", question.QuestionTitle)
	})

}

func TestGetSpecificQuestion(t *testing.T) {

	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	testcases := []struct{
		name string
		id string
		title string
		expectedError error
	}{
		{name: "get specific question", id: "1", title: "Two Sum", expectedError: nil},
		{name: "invalid id", id: "hello world", title: "", expectedError: errors.New("Invalid ID")},
		{name: "paid question", id: "527", title: "Word Abbreviation", expectedError: nil},
	}

    err := os.MkdirAll(filepath.Join(tmpDir, "data"), os.ModePerm)
    assert.NoError(t, err)

	conn, cleanDB := setupDatabase(t)
	defer cleanDB()

	db.DB = conn

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			question, err := GetSpecificQuestion(tc.id)
			
			if tc.expectedError != nil {
				assert.NotNil(t, err)
				return
			}

			assert.Equal(t, tc.id, question.Id)
		})
	}
}

func TestGetRandomDifficultyQuestions(t *testing.T) {

	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	testcases := []struct{
		name string
		difficulty string
		expectedError error
	}{
		{name: "get easy diff.", difficulty: "Easy", expectedError: nil},
		{name: "get medium diff.", difficulty: "Medium", expectedError: nil},
		{name: "get hard diff.", difficulty: "Hard", expectedError: nil},
		{name: "get non existing diff.", difficulty: "Arthur Morgan", expectedError: errors.New("Not found diff")},
	}

    err := os.MkdirAll(filepath.Join(tmpDir, "data"), os.ModePerm)
    assert.NoError(t, err)

	conn, cleanDB := setupDatabase(t)
	defer cleanDB()

	db.DB = conn

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GetRandomDifficulty(tc.difficulty)
			
			if tc.expectedError != nil {
				assert.NotNil(t, err)
				return
			}

		})
	}
}

func TestGetQuestionsFromCompany(t *testing.T) {

	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	testcases := []struct{
		name string
		company string
		expectedError error
	}{
		{name: "get company #1", company: "Microsoft", expectedError: nil},
		{name: "get company #2", company: "Adobe", expectedError: nil},
		{name: "get company #3", company: "Google", expectedError: nil},
		{name: "get non-exist company", company: "I am a company", expectedError: errors.New("No company")},
	}

    err := os.MkdirAll(filepath.Join(tmpDir, "data"), os.ModePerm)
    assert.NoError(t, err)

	conn, cleanDB := setupDatabase(t)
	defer cleanDB()

	db.DB = conn

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			question, err := GetQuestionsFromCompany(tc.company)
			
			if tc.expectedError != nil {
				assert.NotNil(t, err)
				return
			}

			assert.NotEmpty(t, question)
		})
	}
}

func TestGetRandomTopicQuestion(t *testing.T) {

	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	testcases := []struct{
		name string
		topic string
		maxLength int
		expectedError error
	}{
		{name: "topic #1", topic: "Strings", maxLength: 15, expectedError: nil},
		{name: "topic #2", topic: "Arrays", maxLength: 20, expectedError: nil},
		{name: "topic #3", topic: "Trees", maxLength: 10, expectedError: nil},
		{name: "Topic doesn't exist - random top", topic: "Yoga", maxLength: 10, expectedError: nil},
		{name: "Invalid max length", topic: "Tries", maxLength: -1, expectedError: errors.New("Invalid max length")},
	}

    err := os.MkdirAll(filepath.Join(tmpDir, "data"), os.ModePerm)
    assert.NoError(t, err)

	conn, cleanDB := setupDatabase(t)
	defer cleanDB()

	db.DB = conn

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			topic := strings.ToLower(strings.Replace(tc.topic, " ", "-", -1))
			question, err := GetRandomTopicQuestion(topic, tc.maxLength)
			
			if tc.expectedError != nil {
				assert.NotNil(t, err)
				return
			}

			assert.NotEmpty(t, question)
		})
	}
}
