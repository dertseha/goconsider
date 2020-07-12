package text

import (
	"strings"
)

// Wordify processes a text and returns a simplified version of it.
// It reduces all whitespace to single spaces and adds a single whitespace at the begin and end.
// It splits up MixedCaseWords or hyphenated-words.
// Any punctuation is replaced with whitespace.
// Finally, it returns everything lowercase.
func Wordify(s string) string {
	w := s
	w = removePunctuation(w)
	w = mergeBlanks(w)

	oldwords := strings.Split(w, " ")
	newwords := make([]string, 0, len(oldwords))
	for _, word := range oldwords {
		currentPart := ""
		lastCase := runeCase(0)
		addWord := func() {
			if len(currentPart) == 0 {
				return
			}
			newwords = append(newwords, currentPart)
			currentPart = ""
		}
		for _, r := range word {
			currentCase := runeCaseFrom(r)
			newUpper := currentCase == 1 && lastCase != 1
			newLower := currentCase == -1 && lastCase == 0
			if newUpper || newLower {
				addWord()
			}
			lastCase = currentCase
			currentPart += string(r)
		}
		addWord()
	}
	if len(newwords) == 0 {
		return ""
	}
	return " " + strings.ToLower(strings.Join(newwords, " ")) + " "
}

func removePunctuation(w string) string {
	for _, old := range []string{".", "?", "!", ";", ":", "-", "/", "(", ")", "\n", "\r"} {
		w = strings.ReplaceAll(w, old, " ")
	}
	return w
}

func mergeBlanks(w string) string {
	ref := w
	w = strings.ReplaceAll(ref, "  ", " ")
	for len(w) < len(ref) {
		ref = w
		w = strings.ReplaceAll(ref, "  ", " ")
	}
	return w
}

type runeCase int

func runeCaseFrom(r rune) runeCase {
	noLowerChange := strings.ToLower(string(r)) == string(r)
	noUpperChange := strings.ToUpper(string(r)) == string(r)

	if noLowerChange == noUpperChange {
		return 0
	}
	if noUpperChange {
		return 1
	}
	return -1
}
