package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vincent-buchner/leetcode-framer/internal/db"
	"github.com/vincent-buchner/leetcode-framer/internal/models"
	stateservice "github.com/vincent-buchner/leetcode-framer/internal/services/state_service"
	"gorm.io/gorm"
)

func setupTestDir(t *testing.T) (string, func()) {
    // Create a temporary directory for test files
    tmpDir, err := os.MkdirTemp("", "config-test-*")
    if err != nil {
        t.Fatal(err)
    }

    // Change to the temp directory
    originalDir, err := os.Getwd()
    if err != nil {
        t.Fatal(err)
    }
    os.Chdir(tmpDir)

    // Return cleanup function
    return tmpDir, func() {
        os.Chdir(originalDir)
        os.RemoveAll(tmpDir)
    }
}

func TestInit(t *testing.T) {
    app, state, err := Init()
    assert.NoError(t, err)
    assert.NotNil(t, app)
    assert.NotNil(t, state)
    assert.NotNil(t, APP) // Global variable should be set
}

func TestNewDBConnection(t *testing.T) {
    tmpDir, cleanup := setupTestDir(t)
    defer cleanup()

    // Create sqlite_db directory
    err := os.MkdirAll(filepath.Join(tmpDir, "data"), 0755)
    assert.NoError(t, err)

    // Create empty db file
    _, err = os.Create(filepath.Join(tmpDir, "data", "leetcode.db"))
    assert.NoError(t, err)

    conn, err := db.InitDatabase("data/leetcode.db", true)
    assert.NoError(t, err)
    assert.NotNil(t, conn)
    assert.NotNil(t, db.DB) // Global variable should be set
}

func TestApplicationConfig_LoadConfig(t *testing.T) {
    _, cleanup := setupTestDir(t)
    defer cleanup()

    t.Run("Creates new config if doesn't exist", func(t *testing.T) {
        config := &ApplicationState{
            Language: "go",
            Username: "test-user",
            Question: models.QuestionModel{},
            DailyQuestion: models.QuestionModel{},
            TodayDate: time.Now(),
        }
        err := stateservice.StateToJSON(config, "./configs/config.json")
        assert.NoError(t, err)


        // Verify file was created
        _, err = os.Stat("./configs/config.json")
        assert.NoError(t, err)

        // Read and verify content
        var loadedState ApplicationState
        stateservice.JSONToState(&loadedState, "./configs/config.json")
        assert.NoError(t, err)

        assert.NoError(t, err)
        assert.Equal(t, config.Language, loadedState.Language)
        assert.Equal(t, config.Username, loadedState.Username)
    })

    t.Run("Loads existing config", func(t *testing.T) {
        // Create a test config file
        testConfig := ApplicationState{
            Language: "python",
            Username: "existing-user",
            TodayDate: time.Now(),
        }

        configBytes, err := json.Marshal(testConfig)
        assert.NoError(t, err)

        dir := filepath.Dir("./configs/config.json")
        err = os.MkdirAll(dir, os.ModePerm)
        assert.NoError(t, err)

        err = os.WriteFile("./configs/config.json", configBytes, 0644)
        assert.NoError(t, err)

        // Load the config
        var loadedConfig ApplicationState
        stateservice.JSONToState(&loadedConfig, "./configs/config.json")

        assert.Equal(t, testConfig.Language, loadedConfig.Language)
        assert.Equal(t, testConfig.Username, loadedConfig.Username)
    })
}

func TestApplicationConfig_SaveConfig(t *testing.T) {
    _, cleanup := setupTestDir(t)
    defer cleanup()

    config := &ApplicationState{
        Language: "java",
        Username: "save-test-user",
        TodayDate: time.Now(),
    }

    stateservice.StateToJSON(config, "./configs/config.json")

    // Verify file exists
    _, err := os.Stat("./configs/config.json")
    assert.NoError(t, err)

    // Read and verify content
    content, err := os.ReadFile("./configs/config.json")
    assert.NoError(t, err)

    var loadedConfig ApplicationState
    err = json.Unmarshal(content, &loadedConfig)
    assert.NoError(t, err)
    assert.Equal(t, config.Language, loadedConfig.Language)
    assert.Equal(t, config.Username, loadedConfig.Username)
}

func TestLoadConfig_Errors(t *testing.T) {
    tmpDir, cleanup := setupTestDir(t)
    defer cleanup()

    t.Run("Error on config file creation", func(t *testing.T) {
        // Make directory read-only
        err := os.Chmod(tmpDir, 0444)
        if err != nil {
            t.Fatal(err)
        }
        defer os.Chmod(tmpDir, 0755)

        config := &ApplicationState{}
        err = stateservice.StateToJSON(config, "./configs/config.json")

        assert.NotEqual(t, err, nil)
        
    })

    t.Run("Error on config file open", func(t *testing.T) {
        // Create config file with no read permissions
        dir := filepath.Dir("./configs/config.json")
        err := os.MkdirAll(dir, os.ModePerm)
        assert.NoError(t, err)

        _, err = os.OpenFile("./configs/config.json", os.O_CREATE, 0222)
        assert.NoError(t, err)
        defer os.Remove("./configs/config.json")

        config := &ApplicationState{}
        err = stateservice.JSONToState(config, "./configs/config.json")
        assert.NotNil(t, err)
        
    })
}

func TestSaveConfig_Errors(t *testing.T) {
    tmpDir, cleanup := setupTestDir(t)
    defer cleanup()

    t.Run("Error on config file creation for save", func(t *testing.T) {
        // Make directory read-only
        err := os.Chmod(tmpDir, 0444)
        if err != nil {
            t.Fatal(err)
        }
        defer os.Chmod(tmpDir, 0755)

        config := &ApplicationState{}
        err = stateservice.StateToJSON(config, "./configs/config.json")
        assert.NotNil(t, err)

    })
}

func TestNewDBConnection_Errors(t *testing.T) {
    tmpDir, cleanup := setupTestDir(t)
    defer cleanup()

    t.Run("Error on invalid DB path", func(t *testing.T) {
        // Don't create sqlite_db directory - should cause error
        conn, err := db.InitDatabase("data/leetcode.db", true)
        assert.Error(t, err)
        assert.Equal(t, &gorm.DB{}, conn)
    })

    t.Run("Error on inaccessible DB file", func(t *testing.T) {
        // Create sqlite_db directory but make it inaccessible
        err := os.MkdirAll(filepath.Join(tmpDir, "data"), 0000)
        assert.NoError(t, err)
        defer os.Chmod(filepath.Join(tmpDir, "data"), 0755)

        conn, err := db.InitDatabase("data/leetcode.db", true)
        assert.Error(t, err)
        assert.Equal(t, &gorm.DB{}, conn)
    })
}

