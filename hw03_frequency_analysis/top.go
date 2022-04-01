package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var regExp = regexp.MustCompile("[^a-zA-Zа-яА-Я-]+")

func Top10(text string) []string {
	parsed := strings.ToLower(regExp.ReplaceAllString(text, " "))
	words := strings.Fields(parsed)

	occurrences := make(map[string]int)
	for _, word := range words {
		if word == "-" {
			continue
		}
		occurrences[word]++
	}

	unsortedWords := make([]string, len(occurrences))
	var idx int
	for word := range occurrences {
		unsortedWords[idx] = word
		idx++
	}

	sort.Slice(unsortedWords, func(i, j int) bool {
		if occurrences[unsortedWords[i]] == occurrences[unsortedWords[j]] {
			switch strings.Compare(unsortedWords[i], unsortedWords[j]) {
			case 1:
				return false
			default:
				return true
			}
		}

		return occurrences[unsortedWords[i]] > occurrences[unsortedWords[j]]
	})

	if len(unsortedWords) < 10 {
		return unsortedWords
	}

	return unsortedWords[:10]
}
