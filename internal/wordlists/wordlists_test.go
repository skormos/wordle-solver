package wordlists_test

import (
	"testing"

	"github.com/skormos/wordle-solver/internal/wordlists"
)

func TestAllowedList_AllWordsAreFiveLong(t *testing.T) {
	t.Parallel()

	allowedList := wordlists.AllowedList()

	for _, word := range allowedList {
		if 5 != len(word) {
			t.Errorf("got %d, want %d", len(word), 5)
		}
	}
}

func TestAllowedList_AllUniqueEntries(t *testing.T) {
	t.Parallel()

	allowedList := wordlists.AllowedList()
	allowedSet := createAllowedSet(t)
	if len(allowedList) != len(allowedSet) {
		t.Errorf("got %d, want %d", len(allowedSet), len(allowedList))
	}
}

func TestAnswerSet_AllInAllowed(t *testing.T) {
	t.Parallel()

	answerSet := wordlists.AnswerSet()
	allowedSet := createAllowedSet(t)

	for answer := range answerSet {
		if _, contains := allowedSet[answer]; !contains {
			t.Errorf("answer %q not in allowedSet", answer)
		}
	}
}

func createAllowedSet(t *testing.T) map[string]struct{} {
	t.Helper()

	allowedList := wordlists.AllowedList()
	allowedSet := make(map[string]struct{}, len(allowedList))
	for _, word := range allowedList {
		allowedSet[word] = struct{}{}
	}
	return allowedSet
}
