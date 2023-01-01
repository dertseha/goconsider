package goconsider_test

import (
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dertseha/goconsider"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLint(t *testing.T) {
	makeIssue := func(line, col int, prefix string) goconsider.Issue {
		return goconsider.Issue{
			Pos:     token.Position{Line: line, Column: col},
			Message: prefix + " contains 'abcd', consider rephrasing to something else",
		}
	}
	tt := []struct {
		name     string
		expected []goconsider.Issue
	}{
		{
			name:     "testdata/issueFree.go",
			expected: nil,
		},
		{
			name: "testdata/issueInAbcdFilename.go",
			expected: []goconsider.Issue{
				makeIssue(1, 1, "Filename"),
			},
		},
		{
			name: "testdata/issueInPackageComment.go",
			expected: []goconsider.Issue{
				makeIssue(1, 1, "Comment"),
			},
		},
		{
			name: "testdata/issueInFreefloatingComment.go",
			expected: []goconsider.Issue{
				makeIssue(3, 1, "Comment"),
			},
		},
		{
			name: "testdata/issueInInlineComment.go",
			expected: []goconsider.Issue{
				makeIssue(4, 2, "Comment"),
				makeIssue(5, 21, "Comment"),
				makeIssue(9, 22, "Comment"),
			},
		},
		{
			name: "testdata/abcd/issueInPackageName.go",
			expected: []goconsider.Issue{
				makeIssue(1, 9, "Package name"),
			},
		},
		{
			name: "testdata/issueInImportName.go",
			expected: []goconsider.Issue{
				makeIssue(4, 2, "Package alias"),
			},
		},
		{
			name: "testdata/issueInValueName.go",
			expected: []goconsider.Issue{
				makeIssue(5, 7, "Value name"),
				makeIssue(7, 5, "Value name"),
				makeIssue(10, 8, "Value name"),
			},
		},
		{
			name: "testdata/issueInType.go",
			expected: []goconsider.Issue{
				makeIssue(3, 6, "Type name"),
				makeIssue(4, 2, "Member name"),
				makeIssue(7, 6, "Type name"),
				makeIssue(7, 27, "Parameter name"),
				makeIssue(7, 46, "Result name"),
				makeIssue(9, 6, "Type name"),
				makeIssue(10, 2, "Method name"),
				makeIssue(10, 11, "Parameter name"),
				makeIssue(10, 30, "Result name"),
			},
		},
		{
			name: "testdata/issueInFunction.go",
			expected: []goconsider.Issue{
				makeIssue(5, 36, "Function name"),
				makeIssue(5, 7, "Function receiver"),
				makeIssue(5, 53, "Parameter name"),
				makeIssue(5, 72, "Result name"),
				makeIssue(6, 2, "Identifier"),
				makeIssue(6, 22, "Parameter name"),
				makeIssue(9, 2, "Identifier"),
			},
		},
	}

	settings := goconsider.Settings{
		Phrases: []goconsider.Phrase{
			{Synonyms: []string{"abcd"}, Alternatives: nil},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, tc.name, nil, parser.ParseComments)
			if err != nil {
				t.Errorf("File could not be parsed")
				return
			}

			issues := goconsider.Lint(file, fset, settings)
			issuesString := fmt.Sprintf("%v", issues)
			for i := 0; i < len(tc.expected); i++ {
				tc.expected[i].Pos.Filename = tc.name
			}
			expectedString := fmt.Sprintf("%v", tc.expected)
			if issuesString != expectedString {
				t.Errorf("Reported issues are not expected.\nGot: %+v\nWanted: %+v", issuesString, expectedString)
			}
		})
	}
}

func TestLintIssueCount(t *testing.T) {
	tt := []struct {
		name     string
		expected int
	}{
		{name: "testdata/issueFree.go", expected: 0},
		{name: "testdata/issueCountIgnoreTypeUse.go", expected: 1},
		{name: "testdata/issueInFunction.go", expected: 7},
	}

	settings := goconsider.Settings{
		Phrases: []goconsider.Phrase{
			{Synonyms: []string{"abcd"}, Alternatives: nil},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, tc.name, nil, parser.ParseComments)
			if err != nil {
				t.Errorf("File could not be parsed")
				return
			}

			issues := goconsider.Lint(file, fset, settings)

			if tc.expected != len(issues) {
				t.Errorf("Reported issue count is not expected.\nGot: %v\nWanted: %v", len(issues), tc.expected)
			}
		})
	}
}

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
