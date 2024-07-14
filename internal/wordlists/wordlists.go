// Package wordlists is a simple wrapper to maintain the word list files as embeds, and then retrieve them as structures.
package wordlists

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed files/allowed.txt
var allowedTextString string

//go:embed files/answers.txt
var answersTextString string

// AllowedList reads from the local embedded file and creates a string slice of words that are supposed to be allowed.
func AllowedList() []string {
	allowedWords := make([]string, 0, 3000)
	scanner := bufio.NewScanner(strings.NewReader(allowedTextString))
	for scanner.Scan() {
		allowedWords = append(allowedWords, strings.ToUpper(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("splitting allowed text lines via Scanner: %w", err))
	}
	return allowedWords
}

// AnswerSet reads from the local embedded file and create a string set of words that were previous answers.
func AnswerSet() map[string]struct{} {
	answerWords := make(map[string]struct{})
	scanner := bufio.NewScanner(strings.NewReader(answersTextString))
	for scanner.Scan() {
		word := strings.Split(strings.ToUpper(scanner.Text()), " ")[0]
		answerWords[word] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("splitting answer text lines via Scanner: %w", err))
	}
	return answerWords
}