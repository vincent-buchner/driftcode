package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/db"
	"github.com/vincent-buchner/leetcode-framer/internal/models"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// isToday checks if the given date (configDate) matches the current date.
// It returns true if the year, month, and day of configDate are the same as today's date, otherwise false.
func isToday(configDate time.Time) bool {

	today := time.Now()
	return configDate.Year() == today.Year() &&
		configDate.Month() == today.Month() &&
		configDate.Day() == today.Day()
}

// GetDailyQuestion retrieves the daily question from a local server API endpoint.
// If a question has already been requested today, it returns the cached question.
// Otherwise, it sends a GET request to the "/daily" endpoint and parses the
// response into a QuestionModel. The function also processes the topic tags from 
// the response and joins them into a comma-separated string. The result is cached 
// in the application config for future requests. Returns a QuestionModel and an 
// error if any step fails during the process.
func GetDailyQuestion() (models.QuestionModel, error) {

	// Have we requested the daily question already and is it expired?
	if config.STATE.DailyQuestion != (models.QuestionModel{}) && isToday(config.STATE.TodayDate) {
		return config.STATE.DailyQuestion, nil
	}

	// Creates a new request to be sent
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/daily", os.Getenv("API_ENDPOINT")), nil)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't create the request: %s", err)
	}

	// Takes the created request and sends it to the given endpoint
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't send the request: %s", err)
	}

	// Get body (close at end of function)
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't read the body: %s", err)
	}

	// This is the message you'll get if you exceed the api rate limit
	if string(body) == "Too many request from this IP, try again in 1 hour" {
		return models.QuestionModel{}, errors.New("rate limit exceeded on API: 1 hour wait")
	}

	// Take the response and insert it's data into the struct instance
	var questionInfo models.QuestionModel
	err = json.Unmarshal(body, &questionInfo)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't marshall the response body to question: %s", err)
	}

	// Because the topic comes back as an array of objects, we have to define these two structs to get that data
	// We will later convert it to the string
	type TopicTag struct {
		Name string `json:"name"`
	}

	type TopicTagsResponse struct {
		TopicTags []TopicTag `json:"topicTags"`
	}

	var topicResponse TopicTagsResponse
	err = json.Unmarshal([]byte(body), &topicResponse)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't marshall the response body to topics: %s", err)
	}

	// Extract the names and join them with commas
	names := []string{}
	for _, tag := range topicResponse.TopicTags {
		names = append(names, tag.Name)
	}
	result := strings.Join(names, ", ")
	questionInfo.Topics = result

	// If empty, then we couldn't get result from api
	if questionInfo == (models.QuestionModel{}) {
		return models.QuestionModel{}, fmt.Errorf("no response from API, check internet connection and try again")
	}

	// Save to config
	config.STATE.DailyQuestion = questionInfo
	config.STATE.TodayDate = time.Now()

	return questionInfo, nil
}

// GetSpecificQuestion queries the database for a question with the given id and then
// requests that question from the API. It returns the question model from the API.
//
// If the record is not found in the database, it returns an error with a message stating
// that the record was not found. If the API request fails, it returns an error with
// the error message from the request. If the response body can't be marshalled to a
// question model, it returns an error with that message. If the response body is empty,
// it returns an error with a message stating that there was no response from the API.
func GetSpecificQuestion(id string) (models.QuestionModel, error) {

	// Define this struct because the naming convention in the API is different
	type ExtractLink struct {
		Link string `json:"link"`
	}

	// Initialize db and auto migrate
	db := db.DB
	// if err != nil {
	// 	return models.QuestionModel{}, err
	// }

	// Query the question from the db
	var db_question models.QuestionDB
	db.First(&db_question, "id = ?", id)
	if db_question == (models.QuestionDB{}) {
		return models.QuestionModel{}, fmt.Errorf("couldn't get record from db: %s", id)
	}

	// Use the queried data to send a request to the API
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/select?titleSlug=%s", os.Getenv("API_ENDPOINT"), db_question.NameSlug), nil)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't create the request: %s", err)
	}

	// Send request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't send the request: %s", err)
	}

	// Get body (close at end of function)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't read the body: %s", err)
	}

	// JSON to struct
	var questionInfo models.QuestionModel
	err = json.Unmarshal(body, &questionInfo)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't marshall the response body: %s", err)
	}

	// We need this struct to grab the link to the question
	var questionLinkStruct ExtractLink
	err = json.Unmarshal(body, &questionLinkStruct)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't marshall the response body: %s", err)
	}

	questionInfo.QuestionLink = questionLinkStruct.Link

	// Topics to string
	// Define a struct for the topic tag
	type TopicTag struct {
		Name string `json:"name"`
	}

	// Define a struct for the top-level JSON object
	type TopicTagsResponse struct {
		TopicTags []TopicTag `json:"topicTags"`
	}

	var topicResponse TopicTagsResponse
	err = json.Unmarshal([]byte(body), &topicResponse)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't marshall the response body: %s", err)
	}

	// Extract the names and join them with commas
	names := []string{}
	for _, tag := range topicResponse.TopicTags {
		names = append(names, tag.Name)
	}
	result := strings.Join(names, ", ")
	questionInfo.Topics = result

	// If empty, then we couldn't get result from api
	if questionInfo == (models.QuestionModel{}) {
		return models.QuestionModel{}, fmt.Errorf("no response from API, check internet connection and try again")
	}

	return questionInfo, nil

}

// GetRandomDifficulty returns a random question of the given difficulty from the local db.
// Difficulty must be one of the following: Easy, Medium, Hard.
// The function returns a models.QuestionModel which contains the question's title, link, description, and tags.
// If the function fails to get the question or can't read the response body, it returns an error.
func GetRandomDifficulty(difficulty string) (models.QuestionModel, error) {

	type ExtractLink struct {
		Link string `json:"link"`
	}

	// Initialize db and auto migrate
	db := db.DB
	db.AutoMigrate(&models.QuestionDB{})

	// Get all the questions with that difficulty
	var diff_questions []models.QuestionDB
	db.Where("difficulty = ?", difficulty).Find(&diff_questions)

	if len(diff_questions) == 0 {
		return models.QuestionModel{}, errors.New("couldn't retrieve items")
	}

	// Get a random index from those questions and index the question
	rand_index := rand.Intn(len(diff_questions))
	rand_diff_question := diff_questions[rand_index]

	// Use the queried data to send a request to the API
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/select?titleSlug=%s", os.Getenv("API_ENDPOINT"), rand_diff_question.NameSlug), nil)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't create the request: %s", err)
	}

	// Send request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't send the request: %s", err)
	}

	// Get body (close at end of function)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't read the body: %s", err)
	}

	// JSON to struct
	var questionInfo models.QuestionModel
	err = json.Unmarshal(body, &questionInfo)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't marshall the response body: %s", err)
	}
	var questionLink ExtractLink
	err = json.Unmarshal(body, &questionLink)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't marshall the response body: %s", err)
	}
	questionInfo.QuestionLink = questionLink.Link

	// Topics to string
	// Define a struct for the topic tag
	type TopicTag struct {
		Name string `json:"name"`
	}

	// Define a struct for the top-level JSON object
	type TopicTagsResponse struct {
		TopicTags []TopicTag `json:"topicTags"`
	}

	var topicResponse TopicTagsResponse
	err = json.Unmarshal([]byte(body), &topicResponse)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't marshall the response body: %s", err)
	}

	// Extract the names and join them with commas
	names := []string{}
	for _, tag := range topicResponse.TopicTags {
		names = append(names, tag.Name)
	}
	result := strings.Join(names, ", ")
	questionInfo.Topics = result

	// If empty, then we couldn't get result from api
	if questionInfo == (models.QuestionModel{}) {
		return models.QuestionModel{}, fmt.Errorf("no response from API, check internet connection and try again")
	}

	return questionInfo, nil
}

// GetQuestionsFromCompany takes a company name and returns the list of questions associated with that company from the db.
// If the company is not found, it returns an error.
// If the db connection fails, it returns an error.
func GetQuestionsFromCompany(company string) ([]models.QuestionDB, error) {
	// db, err := config.NewDBConnection()
	db := db.DB
	// if err != nil {
	// 	return []models.QuestionDB{}, fmt.Errorf("couldn't connect to db")
	// }

	// Get the company from DB
	var db_company models.CompanyDB
	result := db.First(&db_company, "name = ?", strings.TrimSpace(company))
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []models.QuestionDB{}, fmt.Errorf("NOT FOUND: %s", db_company)
	}

	// Split it's questions into a list
	questionLinks := strings.Split(db_company.Questions, ", ")

	// Query the questions from the Companies to the DB
	var questions []models.QuestionDB
	for _, link := range questionLinks {

		var question models.QuestionDB
		db.First(&question, "link = ?", link)
		questions = append(questions, question)

	}

	return questions, nil
}

// GetRandomTopicQuestion returns a random question of the given topic from the API.
// The topic must be a string that matches one of the topics in the API.
// The maxLength is the maximum number of questions to return.
// The function returns a models.QuestionModel which contains the question's title, link, description, and tags.
// If the function fails to get the question or can't read the response body, it returns an error.
func GetRandomTopicQuestion(topic string, maxLength int) (models.QuestionModel, error)  {

	type ProblemEntry struct {
		Accuracy float64 `json:"acRate"`
		TitleSlug string `json:"titleSlug"`
		ID string `json:"questionFrontendId"`
	}

	type RequestResponse struct {
		ProblemList []ProblemEntry `json:"problemsetQuestionList"`
	}

	if maxLength <= 0 {
		return models.QuestionModel{}, fmt.Errorf("invalid max length, needs to be greater than zero: %d", maxLength)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/problems?tags=%s&limit=%d", os.Getenv("API_ENDPOINT"), topic, maxLength), nil)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't create the request: %s", err)
	}

	// Send request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't send the request: %s", err)
	}

	// Get body (close at end of function)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't read the body: %s", err)
	}

	// JSON to struct
	var topicResponse RequestResponse
	err = json.Unmarshal([]byte(body), &topicResponse)
	if err != nil {
		return models.QuestionModel{}, fmt.Errorf("couldn't marshall the response body: %s", err)
	}

	if len(topicResponse.ProblemList) == 0 {
		return models.QuestionModel{}, fmt.Errorf("empty problem list, couldn't query: %s", topic)
	}

	// Sometimes there are records that can't be found in the db, so let's try 5 times before throwing an error
	var finalQuestion models.QuestionModel
	for i := 0; i < 5; i++ {
		randProblem := topicResponse.ProblemList[rand.Intn(len(topicResponse.ProblemList))]

		question, err := GetSpecificQuestion(randProblem.ID)

		if err == nil {
			finalQuestion = question
			break
		}
	}

	// If after 5 attempts and still nothing, throw error
	if finalQuestion == (models.QuestionModel{}) {
		return models.QuestionModel{}, fmt.Errorf("after 5 attempts, couldn't find question in DB: %s", topic)
	}

	return finalQuestion, nil

}
