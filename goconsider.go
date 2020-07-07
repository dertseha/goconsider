package goconsider

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type Phrase struct {
	Synonyms     []string
	Alternatives []string
}

// Settings contain all the parameters for the analysis.
type Settings struct {
	Phrases []Phrase
	Escapes []string
}

// Issue describes an occurrence of an unwanted phrase.
type Issue struct {
	Pos     token.Position
	Message string
}

// Lint runs analysis on the provided code.
func Lint(file *ast.File, fset *token.FileSet, settings Settings) []Issue {
	var issues []Issue

	for _, group := range file.Comments {
		found := checkComments(fset, group, settings)
		issues = append(issues, found...)
	}

	return issues
}

func checkComments(fset *token.FileSet, group *ast.CommentGroup, settings Settings) []Issue {
	var issues []Issue
	text := group.Text()
	for _, phrase := range settings.Phrases {
		for _, synonym := range phrase.Synonyms {
			if strings.Contains(text, synonym) {
				issue := Issue{
					Pos:     fset.Position(group.Pos()),
					Message: considerMessage("Comment contains", synonym, phrase.Alternatives),
				}
				issues = append(issues, issue)
			}
		}
	}
	return issues
}

func considerMessage(prefix, synonym string, alternatives []string) string {
	return fmt.Sprintf("%s '%s', consider rephrasing to %s", prefix, synonym, alternative(alternatives))
}

func alternative(list []string) string {
	if len(list) == 0 {
		return "something else"
	}
	return "one of [" + strings.Join(list, ", ") + "]"
}
