package controllers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vincent-buchner/leetcode-framer/internal/db"
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
	// os.Chdir(tmpDir)

	// Return cleanup function
	return tmpDir, func() {
		os.Chdir(originalDir)
		os.RemoveAll(tmpDir)
	}
}

func setupDatabase(t *testing.T) (*gorm.DB, func())  {

	tmpDir, err := os.MkdirTemp("", "leetcode-test-*")
    if err != nil {
        t.Fatal(err)
    }

	// original_dir, err := os.Getwd()
	// if err != nil {
	// 	t.Fatal(err)
	// }

    // Create database path
    dbPath := filepath.Join(tmpDir, "leetcode.db")
	
	conn, err := db.InitDatabase(dbPath, true)
	if err != nil {
		t.Fatal(err)
	}

	// Create cleanup function
    cleanup := func() {
        os.RemoveAll(tmpDir)
    }
	db.LoadQuestions(conn, "../../data/static/leetcode_problems.csv")
	db.LoadCompanies(conn, "../../data/static/leetcode_problems_and_companies.csv")

	return conn, cleanup
}
