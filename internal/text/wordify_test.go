package text_test

import (
	"fmt"

	"github.com/dertseha/goconsider/internal/text"
)

func ExampleWordify() {
	w := text.Wordify("This is a SpecialTest     of\nsomething-true.")
	fmt.Printf("'%s'\n", w)
	// Output:
	// ' this is a special test of something true '
}
