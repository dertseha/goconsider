package goconsider_test

import (
	"testing"

	"github.com/dertseha/goconsider"
)

func TestDefaultSettingsPhrasesComeWithReferences(t *testing.T) {
	settings := goconsider.DefaultSettings()
	for _, phrase := range settings.Phrases {
		n := len(phrase.References)
		if n == 0 {
			t.Errorf("Phrase <%v> has no reference."+
				" This makes it difficult to understand why the phrase is flagged.", phrase.Synonyms[0])
		}
	}
}
