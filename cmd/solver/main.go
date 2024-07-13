package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"
)

type alphabetSet struct {
	store map[string]struct{}
}

func (s *alphabetSet) delete(target string) {
	if len(s.store) == 1 {
		return
	}
	delete(s.store, strings.ToUpper(target))
}

func (s *alphabetSet) keepOnly(target string) {
	if len(s.store) > 1 {
		s.store = map[string]struct{}{
			strings.ToUpper(target): {},
		}
	}
}

func (s *alphabetSet) String() string {
	letters := make([]string, 0, 26)
	for l := range s.store {
		letters = append(letters, l)
	}

	sort.Strings(letters)
	return strings.Join(letters, "")
}

func allowedList(filename string) ([]string, error) {
	b, err := os.ReadFile("allowed.txt")
	if err != nil {
		return nil, fmt.Errorf("reading file %w", err)
	}

	allowedWords := make([]string, 0, 3000)
	scanner := bufio.NewScanner(strings.NewReader(string(b)))
	for scanner.Scan() {
		allowedWords = append(allowedWords, strings.ToUpper(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("splitting lines: %w", err)
	}
	return allowedWords, nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	allowedWords, err := allowedList("allowed.txt")
	if err != nil {
		logger.Error("Could not read allowed list.", "error", err)
		os.Exit(1)
	}

	alphaSets := make([]alphabetSet, 5)
	for i := 0; i < 5; i++ {
		alphaSets[i] = newAlphabetSet()
	}

	containsSet := make(map[string]struct{})
	for i := 0; i < 6; i++ {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(fmt.Sprintf("%d > ", i+1))
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input := strings.Split(text[:len(text)-1], " ")
		if len(input) == 1 || input[0] == input[1] {
			logger.Info("Congratulations!")
			os.Exit(0)
		}
		for c := 0; c < 5; c++ {
			if string(input[1][c]) == string(input[0][c]) { // letter is in the right spot
				alphaSets[c].keepOnly(string(input[0][c]))
			} else if string(input[1][c]) == "*" { // letter is in the word, not this spot
				containsSet[string(input[0][c])] = struct{}{}
				alphaSets[c].delete(string(input[0][c]))
			} else if string(input[1][c]) == "_" { // letter is not in this word
				for _, set := range alphaSets {
					set.delete(string(input[0][c]))
				}
			}
		}

		remaining := make([]string, 0, len(allowedWords))
		fmt.Println(calcRegex(alphaSets))
		pattern := regexp.MustCompile("^" + calcRegex(alphaSets) + "$")
	words:
		for _, word := range allowedWords {
			if pattern.MatchString(word) {
				for a := range containsSet {
					if !strings.Contains(word, a) {
						continue words
					}
				}
				remaining = append(remaining, word)
				fmt.Println(word)
			}
		}
		allowedWords = remaining
		fmt.Println(len(allowedWords))
	}

}

func calcRegex(sets []alphabetSet) string {
	out := ""
	for _, set := range sets {
		out += "[" + set.String() + "]"
	}

	return out
}

func newAlphabetSet() alphabetSet {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	out := alphabetSet{
		store: make(map[string]struct{}),
	}
	for _, c := range alphabet {
		out.store[string(c)] = struct{}{}
	}

	return out
}