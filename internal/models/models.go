package models

type QuestionModel struct {
	Id string `json:"questionId"`
	QuestionTitle string `json:"questionTitle"`
	QuestionTitleSlug string `json:"titleSlug"`
	QuestionLink string `json:"questionLink"`
	Question string `json:"question"`
	TestCases string `json:"exampleTestcases"`
	Difficulty string `json:"difficulty"`
	Topics string `json:"topicTagsString"`
}

type QuestionDB struct {
	Name string `json:"name"`
	NameSlug string `json:"nameSlug"`
	Link string `json:"link"`
	Difficulty string `json:"difficulty"`
	ID int `json:"id"`
}

type CompanyDB struct {
	Name string `json:"name"`
	Questions string `json:"questions"`
}
