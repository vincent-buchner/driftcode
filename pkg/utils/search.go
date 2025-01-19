package utils

import (
	"fmt"
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// FuzzySearch returns a string of the top maxResults matches to the query.
// The matches are separated by commas and are ranked by their fuzzy search score.
// If the number of matches is less than maxResults, the function will return all of the matches.
func FuzzySearch(query string, searchList []string, maxResults int) string {
	query_ranks := fuzzy.RankFind(query, searchList)
	sort.Sort(query_ranks)

	if query_ranks.Len() >= maxResults {
		query_ranks = query_ranks[:maxResults]
	} else {
		maxResults = query_ranks.Len()
	}

	var resultString string
	for _, rank := range query_ranks[:maxResults] {
		resultString += fmt.Sprintf("%s, ", rank.Target)
	}

	return resultString
}