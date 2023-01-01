package settings_test

import (
	"testing"

	"github.com/dertseha/goconsider/pkg/settings"
)

func TestDefaultSettingsHaveSynonyms(t *testing.T) {
	s := settings.Default()
	if len(s.Phrases) == 0 {
		t.Errorf("Default settings have no phrases.")
	}
	for index, phrase := range s.Phrases {
		if len(phrase.Synonyms) == 0 {
			t.Errorf("Phrase at index %d has no synonyms.", index)
		}
	}
}
