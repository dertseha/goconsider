package goconsider_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dertseha/goconsider/pkg/analyzer"
	"github.com/dertseha/goconsider/pkg/consider"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	settings := consider.Settings{
		Phrases: []consider.Phrase{
			{Synonyms: []string{"abcd"}, Alternatives: nil},
		},
	}

	analysistest.Run(t, testdataDir(), analyzer.NewAnalyzer(settings), "./...")
}

func testdataDir() string {
	_, testFilename, _, ok := runtime.Caller(1)
	if !ok {
		panic("unable to get current test filename")
	}
	return filepath.Join(filepath.Dir(testFilename), "testdata")
}
