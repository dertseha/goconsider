package goconsider_test

import (
	"fmt"
	"go/parser"
	"go/token"
	"testing"

	"github.com/dertseha/goconsider"
)

func TestLint(t *testing.T) {
	tt := []struct {
		name     string
		expected []goconsider.Issue
	}{
		{
			name:     "testdata/issueFree.go",
			expected: nil,
		},
		{
			name: "testdata/issueInPackageComment.go",
			expected: []goconsider.Issue{
				{
					Pos: token.Position{
						Filename: "testdata/issueInPackageComment.go",
						Line:     1,
						Column:   1,
					},
					Message: "Comment contains 'abcd', consider rephrasing to something else",
				},
			},
		},
		{
			name: "testdata/issueInFreefloatingComment.go",
			expected: []goconsider.Issue{
				{
					Pos: token.Position{
						Filename: "testdata/issueInFreefloatingComment.go",
						Line:     3,
						Column:   1,
					},
					Message: "Comment contains 'abcd', consider rephrasing to something else",
				},
			},
		},
		{
			name: "testdata/issueInInlineComment.go",
			expected: []goconsider.Issue{
				{
					Pos: token.Position{
						Filename: "testdata/issueInInlineComment.go",
						Line:     4,
						Column:   2,
					},
					Message: "Comment contains 'abcd', consider rephrasing to something else",
				},
				{
					Pos: token.Position{
						Filename: "testdata/issueInInlineComment.go",
						Line:     5,
						Column:   21,
					},
					Message: "Comment contains 'abcd', consider rephrasing to something else",
				},
				{
					Pos: token.Position{
						Filename: "testdata/issueInInlineComment.go",
						Offset:   0,
						Line:     9,
						Column:   22,
					},
					Message: "Comment contains 'abcd', consider rephrasing to something else",
				},
			},
		},
		{
			name: "testdata/issueInType.go",
			expected: []goconsider.Issue{
				{
					Pos: token.Position{
						Filename: "testdata/issueInType.go",
						Line:     3,
						Column:   6,
					},
					Message: "Type name contains 'abcd', consider rephrasing to something else",
				},
				{
					Pos: token.Position{
						Filename: "testdata/issueInType.go",
						Line:     4,
						Column:   2,
					},
					Message: "Member name contains 'abcd', consider rephrasing to something else",
				},
				{
					Pos: token.Position{
						Filename: "testdata/issueInType.go",
						Line:     7,
						Column:   6,
					},
					Message: "Type name contains 'abcd', consider rephrasing to something else",
				},
				{
					Pos: token.Position{
						Filename: "testdata/issueInType.go",
						Line:     7,
						Column:   27,
					},
					Message: "Parameter name contains 'abcd', consider rephrasing to something else",
				},
				{
					Pos: token.Position{
						Filename: "testdata/issueInType.go",
						Line:     7,
						Column:   46,
					},
					Message: "Result name contains 'abcd', consider rephrasing to something else",
				},
			},
		},
	}

	settings := goconsider.Settings{
		Phrases: []goconsider.Phrase{
			{Synonyms: []string{"abcd"}, Alternatives: nil},
		},
		Escapes: nil,
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
			expectedString := fmt.Sprintf("%v", tc.expected)
			if issuesString != expectedString {
				t.Errorf("Reported issues are not expected.\nGot: %+v\nWanted: %+v", issuesString, expectedString)
			}
		})
	}
}
