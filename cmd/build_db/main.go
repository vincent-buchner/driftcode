package main

import (
	"github.com/vincent-buchner/leetcode-framer/internal/db"
)

func main() {
	conn, err := db.InitDatabase("data/leetcode.db", true)
	if err != nil {
		panic(err)
	}
	db.LoadQuestions(conn, "data/static/leetcode_problems.csv")
	db.LoadCompanies(conn, "data/static/leetcode_problems_and_companies.csv")
}