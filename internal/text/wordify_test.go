package text_test

import (
	"fmt"
	"testing"

	"github.com/dertseha/goconsider/internal/text"
)

func ExampleWordify() {
	w := text.Wordify("This is a SpecialTest     of\nsomething-true.")
	fmt.Printf("'%s'\n", w)
	// Output:
	// ' this is a special test of something true '
}

func TestWordifyReturnsEmptyStringIfEmpty(t *testing.T) {
	w := text.Wordify("  ")
	if len(w) != 0 {
		t.Errorf("Expected empty string, got '" + w + "'")
	}
}
