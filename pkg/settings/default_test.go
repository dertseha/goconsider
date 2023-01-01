package settings_test

import (
	"testing"

	"github.com/dertseha/goconsider/pkg/settings"
)

func TestDefaultSettingsExist(t *testing.T) {
	s := settings.Default()
	if len(s.Phrases) == 0 {
		t.Errorf("Default settings have no phrases.")
	}
}

func TestDefaultSettingsHaveSynonyms(t *testing.T) {
	s := settings.Default()
	for index, phrase := range s.Phrases {
		if len(phrase.Synonyms) == 0 {
			t.Errorf("Phrase at index %d has no synonyms.", index)
		}
	}
}

func TestDefaultSettingsHaveReferences(t *testing.T) {
	s := settings.Default()
	for index, phrase := range s.Phrases {
		if len(phrase.References) == 0 {
			t.Errorf("Phrase at index %d has no references. References are required for the default settings.", index)
		}
		for _, shortRef := range phrase.References {
			if long := s.References[shortRef]; len(long) == 0 {
				t.Errorf("Phrase at index %d has unknown short reference '%s'. "+
					"References for the default settings must be traceable.", index, shortRef)
			}
		}
	}
}
