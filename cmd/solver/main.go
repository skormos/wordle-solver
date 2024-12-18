package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/skormos/wordle-solver/internal/wordlists"
)

const (
	MATCH_GRAY = iota
	MATCH_YELLOW
	MATCH_GREEN
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

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(writer io.Writer, _ []string) error {
	logger := slog.New(slog.NewTextHandler(writer, nil))

	/**
	* 1. Collect file sources
	* 2. Create initial word list
	* 3. Create initial regex
	  4. Read word and match result
	* 5. Update regex
	* 6. Apply regex to word list
	* 7. Go to 4.
	*/

	allowedWords := wordlists.AllowedList()
	puzzleDate := time.Now()
	//puzzleDate, err := time.Parse(time.DateOnly, "2024-02-21")
	//if err != nil {
	//	return fmt.Errorf("could not parse puzzle date: %w", err)
	//}
	answers := wordlists.AnswerSet(puzzleDate)

	testList := make([]string, 0, len(allowedWords))
	for _, word := range allowedWords {
		if _, ok := answers[word]; !ok {
			testList = append(testList, word)
		}
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
			return fmt.Errorf("while reading user input: %w", err)
		}

		input := strings.Split(text[:len(text)-1], " ")
		if input[0] == "" {
			i--
			continue
		}
		if len(input) == 1 || input[0] == input[1] {
			_, _ = fmt.Fprintf(writer, "Congratulations!")
			return nil
		}
		for c := 0; c < 5; c++ {
			guess := string(input[0][c])
			result := string(input[1][c])
			match := -1
			switch result {
			case "_":
				match = MATCH_GRAY
			case "*":
				match = MATCH_YELLOW
			default:
				match = MATCH_GREEN
			}
			switch match {
			case MATCH_GRAY: // delete from all
				for _, set := range alphaSets {
					set.delete(guess)
				}
			case MATCH_YELLOW: // delete from this set only
				containsSet[guess] = struct{}{}
				alphaSets[c].delete(guess)
			case MATCH_GREEN: // delete from others
				alphaSets[c].keepOnly(guess)
			}
		}

		remaining := make([]string, 0, len(testList))
		logger.Info(calcRegex(alphaSets))
		pattern := regexp.MustCompile("^" + calcRegex(alphaSets) + "$")
	words:
		for _, word := range testList {
			if pattern.MatchString(word) {
				for a := range containsSet {
					if !strings.Contains(word, a) {
						continue words
					}
				}
				remaining = append(remaining, word)
			}
		}
		testList = remaining
		sort.Strings(testList)
		_, _ = fmt.Fprintf(writer, "--- Remaining words ---\n")
		_, _ = fmt.Fprintf(writer, "%s\n", formatRemaining(testList, 5))
		_, _ = fmt.Fprintf(writer, "=== Count :: %d\n", len(testList))
	}

	return nil
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

func formatRemaining(remaining []string, width int) string {
	output := ""

	for count, word := range remaining {
		output += word
		if (count+1)%width == 0 {
			output += "\n"
		} else {
			output += " "
		}
	}

	return output
}
