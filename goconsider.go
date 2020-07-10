package goconsider

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/dertseha/goconsider/internal/text"
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
	for _, decl := range file.Decls {
		switch typedDecl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range typedDecl.Specs {
				switch typedSpec := spec.(type) {
				case *ast.TypeSpec:
					issues = append(issues, checkType(fset, typedSpec, settings)...)
				}
			}
		case *ast.FuncDecl:
			issues = append(issues, checkFunction(fset, typedDecl, settings)...)
		}

	}

	return issues
}

func checkComments(fset *token.FileSet, group *ast.CommentGroup, settings Settings) []Issue {
	return checkGeneric(group.Text(), settings, "Comment", fset.Position(group.Pos()))
}

func checkType(fset *token.FileSet, typeSpec *ast.TypeSpec, settings Settings) []Issue {
	issues := checkGeneric(typeSpec.Name.Name, settings, "Type name", fset.Position(typeSpec.Name.Pos()))
	issues = append(issues, checkTypeExpr(fset, typeSpec.Type, settings)...)
	return issues
}

func checkTypeExpr(fset *token.FileSet, typeExpr ast.Expr, settings Settings) []Issue {
	var issues []Issue
	switch spec := typeExpr.(type) {
	case *ast.StructType:
		issues = append(issues, checkFieldList(fset, spec.Fields, "Member name", settings)...)
	case *ast.FuncType:
		issues = append(issues, checkFuncType(fset, spec, settings)...)
	case *ast.InterfaceType:
		issues = append(issues, checkFieldList(fset, spec.Methods, "Method name", settings)...)
	}
	return issues
}

func checkFuncType(fset *token.FileSet, funcType *ast.FuncType, settings Settings) []Issue {
	var issues []Issue
	issues = append(issues, checkFieldList(fset, funcType.Params, "Parameter name", settings)...)
	issues = append(issues, checkFieldList(fset, funcType.Results, "Result name", settings)...)
	return issues
}

func checkFieldList(fset *token.FileSet, fields *ast.FieldList, prefix string, settings Settings) []Issue {
	var issues []Issue
	if fields == nil {
		return nil
	}
	for _, field := range fields.List {
		for _, name := range field.Names {
			issues = append(issues, checkGeneric(name.Name, settings, prefix, fset.Position(name.Pos()))...)
		}
		issues = append(issues, checkTypeExpr(fset, field.Type, settings)...)
	}
	return issues
}

func checkFunction(fset *token.FileSet, funcDecl *ast.FuncDecl, settings Settings) []Issue {
	var issues []Issue
	issues = append(issues, checkGeneric(funcDecl.Name.Name, settings, "Function name", fset.Position(funcDecl.Name.Pos()))...)
	issues = append(issues, checkFieldList(fset, funcDecl.Recv, "Function receiver", settings)...)
	issues = append(issues, checkFuncType(fset, funcDecl.Type, settings)...)
	// TODO: body
	return issues
}

func checkGeneric(s string, settings Settings, typeString string, pos token.Position) []Issue {
	var issues []Issue
	addIssue := func(synonym string, alternatives []string) {
		issue := Issue{
			Pos:     pos,
			Message: considerMessage(typeString+" contains", synonym, alternatives),
		}
		issues = append(issues, issue)
	}
	worded := text.Wordify(s)
	for _, phrase := range settings.Phrases {
		for _, synonym := range phrase.Synonyms {
			if strings.Contains(worded, " "+synonym+" ") {
				addIssue(synonym, phrase.Alternatives)
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
