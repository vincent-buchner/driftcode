package controllers

import (
	// "encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	// "path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/models"
	stateservice "github.com/vincent-buchner/leetcode-framer/internal/services/state_service"
)


func TestCreateQuestionProject(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	var testcases = []struct {
		name      string
		lang      string
		extension string
        fileName string
        expectedError error
        id string
	}{
		{name: "create java file", lang: "java", extension: ".java", id: "0", fileName: "hello_world", expectedError: nil},
		{name: "create python file", lang: "python", extension: ".py", id: "0", fileName: "hello_world", expectedError: nil},
		{name: "create go file", lang: "golang", extension: ".go", id: "0", fileName: "hello_world", expectedError: nil},
		{name: "create html file - error", lang: "html", extension: ".html", id: "0", fileName: "hello_world", expectedError: assert.AnError},
	}

	t.Run("Creates new Python Project", func(t *testing.T) {

		for _, tc := range testcases {

			t.Run(tc.name, func(t *testing.T) {

				app_config := &config.ApplicationState{
					Language:      tc.lang,
					Username:      "test-user",
					Question:      models.QuestionModel{},
					DailyQuestion: models.QuestionModel{},
					TodayDate:     time.Now(),
				}
				config.STATE = *app_config
				err := stateservice.StateToJSON(config.STATE, "./configs/config.json")
				assert.NoError(t, err)

				err = CreateQuestionProject(tmpDir, &models.QuestionModel{
					Id:                tc.id,
					QuestionTitle:     "Hello World",
					QuestionTitleSlug: tc.fileName,
					QuestionLink:      "google.com",
					Question:          "What's your favorite color",
					TestCases:         "n == 2",
					Topics:            "Personal",
					Difficulty:        "Pretty Good",
				})

                // Did we expect and error and are we getting an error?
                if tc.expectedError != nil {
                    assert.NotNil(t, err)
                    return
                }

				// Did we create the dir
				_, err = os.Stat(filepath.Join(tmpDir, fmt.Sprintf("%s_%s", tc.id, tc.fileName)))
				assert.NoError(t, err)

				// Is there a file in the dir
				_, err = os.Stat(filepath.Join(tmpDir, fmt.Sprintf("%s_%s/%s%s",tc.id, tc.fileName, tc.fileName, tc.extension)))
				assert.NoError(t, err)
			})

		}
	})
}
