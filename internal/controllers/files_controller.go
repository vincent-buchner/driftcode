package controllers

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/vincent-buchner/leetcode-framer/internal/config"
	"github.com/vincent-buchner/leetcode-framer/internal/models"
	"github.com/vincent-buchner/leetcode-framer/internal/constants"
	"golang.org/x/net/html"
)

// CleanHTML takes in a string of HTML content and returns a plain text string.
// This is useful for stripping away HTML tags from the LeetCode problem text.
func cleanHTML(input string) string {

	// Parse the HTML input using the html package from the standard library.
	// This creates a DOM tree from the input string.
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		// If there's an error parsing the HTML, print the error and return an empty string.
		fmt.Println("Error parsing HTML:", err)
		return ""
	}

	// We'll define a recursive function extractText to traverse the DOM tree and
	// extract all text content from the HTML document.
	var buf bytes.Buffer
	var extractText func(*html.Node)

	// Define the extractText function. This function takes an HTML node as input
	// and adds all text content to the buf buffer.
	extractText = func(n *html.Node) {
		// If the node is a text node, add its content to the buffer
		if n.Type == html.TextNode {
			// n.Data is the text content of the node
			buf.WriteString(n.Data)
		}

		// If the node is a <sup> tag, prepend "^" to its text content
		if n.Type == html.ElementNode && n.Data == "sup" {
			buf.WriteString("^")
		}

		// Recursively traverse the DOM tree by calling extractText on all child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			// c is the current child node
			extractText(c)
		}

		// Close the ^ if the current node is <sup>
		if n.Type == html.ElementNode && n.Data == "sup" {
			buf.WriteString("")
		}
	}

	// Call extractText on the root node of the DOM tree (doc)
	extractText(doc)

	// Finally, return the plain text content of the HTML document, but strip any
	// extra spaces and return the plain text
	return strings.TrimSpace(buf.String())
}

func toCamelCase(input string) string {
	words := strings.Split(input, "-")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

func createCodeFile(path string, question models.QuestionModel, format string, extension string) {

	// Create file
	new_file_path := fmt.Sprintf("%s/%s.%s", path, question.QuestionTitleSlug, extension)
	f, err := os.Create(new_file_path)
	if err != nil {
		fmt.Println(err)
		return 
	}
	defer f.Close()

	fmt.Fprint(f, format)
}

func CreateQuestionProject(path string, question *models.QuestionModel) error  {

	// Make directory
	folder_name := question.QuestionTitleSlug
	new_folder_path := fmt.Sprintf("%s/%s_%s", path, question.Id, folder_name)
	err := os.MkdirAll(new_folder_path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't make the directory: %v", err)
	}

	// Handling empty questions (paid only request)
	if question.Question == "" {
		question.Question = "There was an issue grabbing the question (likely paid-only)"
	}

	// Switch between the different requests of languages
	var format string
	switch config.STATE.Language {
	case "python":
		format = fmt.Sprintf(constants.PYTHON_FILE, cleanHTML(question.Question), question.QuestionLink, strings.ReplaceAll(question.QuestionTitleSlug, "-", "_"))
		createCodeFile(new_folder_path, *question, format, "py")
	case "java":
		format = fmt.Sprintf(constants.JAVA_FILE, cleanHTML(question.Question), question.QuestionLink, toCamelCase(question.QuestionTitleSlug))
		createCodeFile(new_folder_path, *question, format, "java")
	case "golang":
		format = fmt.Sprintf(constants.GOLANG_FILE, cleanHTML(question.Question), question.QuestionLink, toCamelCase(question.QuestionTitleSlug))
		createCodeFile(new_folder_path, *question, format, "go")
	default:
		return fmt.Errorf("given a language that's not supported: %s", config.STATE.Language)
	}

	return nil
}