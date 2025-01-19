package stateservice

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// StateToJSON writes the given state object to a JSON file at the specified file path.
// It creates the file if it does not exist, or truncates it if it does.
// The JSON output is formatted with indentation for readability.
// Returns an error if the file cannot be created or if encoding fails.
func StateToJSON(state any, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func JSONToState(state any, filePath string) error  {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return errors.New("file is empty")
	}

	if err := json.Unmarshal(data, state); err != nil {
		return err
	}

	return nil
}

