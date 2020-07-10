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

type issueCollector struct {
	settings Settings
	fset     *token.FileSet
	issues   []Issue
}

// Lint runs analysis on the provided code.
func Lint(file *ast.File, fset *token.FileSet, settings Settings) []Issue {
	col := issueCollector{
		settings: settings,
		fset:     fset,
	}

	for _, group := range file.Comments {
		col.checkComments(group)
	}
	for _, decl := range file.Decls {
		switch typedDecl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range typedDecl.Specs {
				switch typedSpec := spec.(type) {
				case *ast.TypeSpec:
					col.checkType(typedSpec)
				}
			}
		case *ast.FuncDecl:
			col.checkFunction(typedDecl)
		}

	}
	return col.issues
}

func (col *issueCollector) checkComments(group *ast.CommentGroup) {
	col.checkGeneric(group.Text(), "Comment", group.Pos())
}

func (col *issueCollector) checkType(typeSpec *ast.TypeSpec) {
	col.checkGeneric(typeSpec.Name.Name, "Type name", typeSpec.Name.Pos())
	col.checkTypeExpr(typeSpec.Type)
}

func (col *issueCollector) checkTypeExpr(typeExpr ast.Expr) {
	switch spec := typeExpr.(type) {
	case *ast.StructType:
		col.checkFieldList(spec.Fields, "Member name")
	case *ast.FuncType:
		col.checkFuncType(spec)
	case *ast.InterfaceType:
		col.checkFieldList(spec.Methods, "Method name")
	}
}

func (col *issueCollector) checkFuncType(funcType *ast.FuncType) {
	col.checkFieldList(funcType.Params, "Parameter name")
	col.checkFieldList(funcType.Results, "Result name")
}

func (col *issueCollector) checkFieldList(fields *ast.FieldList, prefix string) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		for _, name := range field.Names {
			col.checkGeneric(name.Name, prefix, name.Pos())
		}
		col.checkTypeExpr(field.Type)
	}
}

func (col *issueCollector) checkFunction(funcDecl *ast.FuncDecl) {
	col.checkGeneric(funcDecl.Name.Name, "Function name", funcDecl.Name.Pos())
	col.checkFieldList(funcDecl.Recv, "Function receiver")
	col.checkFuncType(funcDecl.Type)
	// TODO: body
}

func (col *issueCollector) checkGeneric(s string, typeString string, pos token.Pos) {
	worded := text.Wordify(s)
	for _, phrase := range col.settings.Phrases {
		for _, synonym := range phrase.Synonyms {
			if strings.Contains(worded, " "+synonym+" ") {
				col.addIssue(typeString, pos, synonym, phrase.Alternatives)
			}
		}
	}
}

func (col *issueCollector) addIssue(typeString string, pos token.Pos, synonym string, alternatives []string) {
	issue := Issue{
		Pos:     col.fset.Position(pos),
		Message: considerMessage(typeString+" contains", synonym, alternatives),
	}
	col.issues = append(col.issues, issue)
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
