package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var regExp, _ = regexp.Compile("[^a-zA-Zа-яА-Я-]+")

var delSyms = []string{"-"}

func Top10(text string) []string {
	parsed := regExp.ReplaceAllString(text, " ")
	parsed = strings.ToLower(parsed)
	words := strings.Fields(parsed)

	var occurences = make(map[string]int)
	for _, word := range words {
		occurences[word] += 1
	}
	for _, delSym := range delSyms {
		delete(occurences, delSym)
	}

	var unsortedWords []string
	for word := range occurences {
		unsortedWords = append(unsortedWords, word)
	}

	sort.Slice(unsortedWords, func(i, j int) bool {
		if occurences[unsortedWords[i]] == occurences[unsortedWords[j]] {
			switch strings.Compare(unsortedWords[i], unsortedWords[j]) {
			case 1:
				return false
			default:
				return true
			}
		}

		return occurences[unsortedWords[i]] > occurences[unsortedWords[j]]
	})

	if len(unsortedWords) < 10 {
		return unsortedWords[:]
	}

	return unsortedWords[:10]
}
