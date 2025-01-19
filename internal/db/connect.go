package db

import (
	"fmt"

	"github.com/vincent-buchner/leetcode-framer/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes a SQLite database connection using the provided file path.
// It performs automatic migrations for the QuestionDB and CompanyDB models.
// Returns a pointer to the gorm.DB instance and an error if the connection fails.
func InitDatabase(filePath string, verbose bool) (*gorm.DB, error) {

	var db_config *gorm.Config

	if !verbose {
		db_config = &gorm.Config{
			Logger: logger.Discard,
		}
	} else {
		db_config = &gorm.Config{}
	}

	db, err := gorm.Open(sqlite.Open(filePath), db_config)
	if err != nil {
		return &gorm.DB{}, fmt.Errorf("couldn't connect to db: %s", err)
	}
	db.AutoMigrate(&models.QuestionDB{})
	db.AutoMigrate(&models.CompanyDB{})

	DB = db

	return db, nil
}