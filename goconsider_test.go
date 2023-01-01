package goconsider_test

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dertseha/goconsider/pkg/analyzer"
	"github.com/dertseha/goconsider/pkg/consider"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestReporting(t *testing.T) {
	settings := consider.Settings{
		Phrases: []consider.Phrase{
			{Synonyms: []string{"abcd"}, Alternatives: nil},
			{Synonyms: []string{"a long phrase that may go across multiple lines"}, Alternatives: nil},
		},
	}

	analysistest.Run(t, testdataDir(t, "reporting"), analyzer.NewAnalyzer(settings), "./...")
}

func TestSettingsDefault(t *testing.T) {
	cdWorkingDir(t, "settings", "default")
	a := analyzer.NewAnalyzerFromFlags()
	_ = a.Flags.Parse([]string{})
	analysistest.Run(t, testdataDir(t, "settings", "default"), a, "./...")
}

func TestSettingsImplicit(t *testing.T) {
	cdWorkingDir(t, "settings", "implicit")
	a := analyzer.NewAnalyzerFromFlags()
	_ = a.Flags.Parse([]string{})
	analysistest.Run(t, testdataDir(t, "settings", "implicit"), a, "./...")
}

func TestSettingsExplicit(t *testing.T) {
	cdWorkingDir(t, "settings", "explicit")
	a := analyzer.NewAnalyzerFromFlags()
	_ = a.Flags.Parse([]string{"-settings", "explicit.yaml"})
	analysistest.Run(t, testdataDir(t, "settings", "explicit"), a, "./...")
}

func TestSettingsWithReferences(t *testing.T) {
	cdWorkingDir(t, "settings", "references")
	a := analyzer.NewAnalyzerFromFlags()
	_ = a.Flags.Parse([]string{})
	analysistest.Run(t, testdataDir(t, "settings", "references"), a, "./...")
}

func cdWorkingDir(tb testing.TB, nested ...string) {
	tb.Helper()
	base := testBaseDir(tb)
	tb.Cleanup(func() { _ = os.Chdir(base) })
	allPaths := []string{testBaseDir(tb), "testdata"}
	allPaths = append(allPaths, nested...)
	err := os.Chdir(path.Join(allPaths...))
	if err != nil {
		tb.Fatalf("Failed to change test directory: %v", err)
	}
}

func testdataDir(tb testing.TB, nested ...string) string {
	tb.Helper()
	allPaths := []string{testBaseDir(tb), "testdata"}
	allPaths = append(allPaths, nested...)
	return filepath.Join(allPaths...)
}

func testBaseDir(tb testing.TB) string {
	tb.Helper()
	_, testFilename, _, ok := runtime.Caller(1)
	if !ok {
		tb.Fatalf("unable to get current test filename")
	}
	return filepath.Dir(testFilename)
}
