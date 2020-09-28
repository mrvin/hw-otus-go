// A function that takes text as input and returns a slice with the 10 most
// common words in the text.
package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

const numTOP = 10

func Top10(str string) []string {
	wordCounts := make(map[string]int)

	words := splitStr(str)

	for _, word := range words {
		word = strings.ToLower(word)
		wordCounts[word]++
	}

	uniqueWords := make([]string, len(wordCounts))
	i := 0
	for word := range wordCounts {
		uniqueWords[i] = word
		i++
	}

	sort.Slice(uniqueWords, func(i, j int) bool {
		return wordCounts[uniqueWords[i]] > wordCounts[uniqueWords[j]]
	})

	numTopWords := numTOP
	if len(uniqueWords) < numTOP {
		numTopWords = len(uniqueWords)
	}

	return uniqueWords[:numTopWords]
}

// splitStr splits the string into words.
func splitStr(str string) []string {
	regExp := regexp.MustCompile(`[\p{L}\p{M}\p{Nd}\p{Pc}]+([\p{Pd}'][\p{L}\p{M}\p{Nd}\p{Pc}]+)*`)
	result := regExp.FindAllString(str, -1)

	return result
}
