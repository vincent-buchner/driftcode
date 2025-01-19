package urlservice

import (
	"fmt"
	"regexp"
	"strings"
)

func ExtractProblemNumber(url string) (string, error) {
	// Define a regular expression to match the problem number.
	re := regexp.MustCompile(`_([0-9]+)\.java`)

	// Find the first match
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1], nil // Return the captured group
	}

	return "", fmt.Errorf("problem number not found in URL")
}

func ExtractNameSlug(url string) (string, error)  {
	// Check if the URL contains the "/problems/" segment
	const segment = "/problems/"
	if !strings.Contains(url, segment) {
		return "", fmt.Errorf("URL does not contain %q segment", segment)
	}
	
	// Split the URL at "/problems/"
	parts := strings.Split(url, segment)
	if len(parts) < 2 {
		return "", fmt.Errorf("URL is malformed")
	}
	
	// Get the part after "/problems/"
	slugPart := parts[1]
	
	// Remove any trailing slashes or query parameters
	slug := strings.Split(slugPart, "/")[0]
	
	return slug, nil
}