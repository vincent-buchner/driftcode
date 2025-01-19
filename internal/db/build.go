package db

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vincent-buchner/leetcode-framer/internal/constants"
	"github.com/vincent-buchner/leetcode-framer/internal/models"
	urlservice "github.com/vincent-buchner/leetcode-framer/internal/services/url_service"
	"gorm.io/gorm"
)

func LoadCompanies(db *gorm.DB, filePath string) {

	// Drop all tables before rerunning
	err := db.Migrator().DropTable(&models.CompanyDB{})
	if err != nil {
		log.Println("Could't drop table, continuing: ", err)
	}
	db.AutoMigrate(&models.CompanyDB{})

	// Load in the csv file
	// ./data/leetcode_problems_and_companies.csv
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Couldn't load file: ", err)
	}
	defer file.Close()

	// Create a new reader and get all records
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Couldn't read records", err)
	}

	for _, company := range constants.COMPANIES {
		
		var questionLinks []string
		for _, record := range records {

			if strings.EqualFold(company, record[2]) {
				questionLinks = append(questionLinks, record[0])
			}
		}

		linksText := strings.Join(questionLinks, ", ")

		db.Create(&models.CompanyDB{
			Name: company,
			Questions: linksText,
		})
	}
}

func LoadQuestions(db *gorm.DB, filePath string) {

	// Drop all tables before rerunning
	err := db.Migrator().DropTable(&models.QuestionDB{})
	if err != nil {
		log.Println("Could't drop table, continuing: ", err)
	}

	// Auto migrate it with the schema
	db.AutoMigrate(&models.QuestionDB{})

	// Load in the csv file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Couldn't load file: ", err)
	}
	defer file.Close()

	// Create a new reader and get all records
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Couldn't read records", err)
	}

	// Iterate over records, adding them to the sqlite file
	// 0 - Name
	// 1 - Link
	// 2 - Difficulty
	// 3 - Solution/Number
	for _, record := range records {

		problemNumber, err := urlservice.ExtractProblemNumber(record[3])
		if err != nil {
			continue
		}
		nameSlug, err := urlservice.ExtractNameSlug(record[1])
		if err != nil {
			log.Println("Couldn't find slug", err)
			continue
		}

		intProblemNumber, _ := strconv.Atoi(problemNumber)
		db.Create(&models.QuestionDB{
			Name: record[0],
			NameSlug: nameSlug,
			Link: record[1],
			Difficulty: strings.TrimSpace(record[2]),
			ID: intProblemNumber,
		})
	}
}