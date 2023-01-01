package goconsider_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dertseha/goconsider"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	settings := goconsider.Settings{
		Phrases: []goconsider.Phrase{
			{Synonyms: []string{"abcd"}, Alternatives: nil},
		},
	}

	analysistest.Run(t, testdataDir(), goconsider.NewAnalyzer(settings), "./...")
}

func testdataDir() string {
	_, testFilename, _, ok := runtime.Caller(1)
	if !ok {
		panic("unable to get current test filename")
	}
	return filepath.Join(filepath.Dir(testFilename), "testdata")
}
