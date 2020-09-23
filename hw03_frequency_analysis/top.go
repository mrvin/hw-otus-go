// A function that takes text as input and returns a slice with the 10 most
// common words in the text.
package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

const numTOP = 10

type wordCount struct {
	word  string
	count int
}

func Top10(str string) []string {
	wordCounts := make(map[string]int)

	words := splitStr(str)

	for _, word := range words {
		word = strings.ToLower(word)
		wordCounts[word]++
	}

	sortWordCounts := make([]wordCount, len(wordCounts))
	i := 0
	for word := range wordCounts {
		sortWordCounts[i] = wordCount{word, wordCounts[word]}
		i++
	}

	sort.Slice(sortWordCounts, func(i, j int) bool {
		return sortWordCounts[i].count > sortWordCounts[j].count
	})

	numTopWords := numTOP
	if len(sortWordCounts) < numTOP {
		numTopWords = len(sortWordCounts)
	}
	top10Words := make([]string, numTopWords)
	for i, strutWordCount := range sortWordCounts[:numTopWords] {
		top10Words[i] = strutWordCount.word
	}

	return top10Words
}

// splitStr splits the string into words.
func splitStr(str string) []string {
	regExp := regexp.MustCompile(`[\p{L}\p{M}\p{Nd}\p{Pc}]+([\p{Pd}'][\p{L}\p{M}\p{Nd}\p{Pc}]+)*`)
	result := regExp.FindAllString(str, -1)

	return result
}
