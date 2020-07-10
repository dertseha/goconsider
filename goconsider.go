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

func considerMessage(prefix, synonym string, alternatives []string) string {
	return fmt.Sprintf("%s '%s', consider rephrasing to %s", prefix, synonym, alternative(alternatives))
}

func alternative(list []string) string {
	if len(list) == 0 {
		return "something else"
	}
	return "one of [" + strings.Join(list, ", ") + "]"
}

// Lint runs analysis on the provided code.
func Lint(file *ast.File, fset *token.FileSet, settings Settings) []Issue {
	col := issueCollector{
		settings: settings,
		fset:     fset,
	}
	// TODO: package name
	col.checkCommentGroups(file.Comments)
	col.checkDecls(file.Decls)
	return col.issues
}

func (col *issueCollector) addIssue(typeString string, pos token.Pos, synonym string, alternatives []string) {
	issue := Issue{
		Pos:     col.fset.Position(pos),
		Message: considerMessage(typeString+" contains", synonym, alternatives),
	}
	col.issues = append(col.issues, issue)
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

func (col *issueCollector) checkIdents(idents []*ast.Ident, prefix string) {
	for _, ident := range idents {
		col.checkIdent(ident, prefix)
	}
}

func (col *issueCollector) checkIdent(ident *ast.Ident, typeString string) {
	col.checkGeneric(ident.Name, typeString, ident.NamePos)
}

func (col *issueCollector) checkCommentGroups(groups []*ast.CommentGroup) {
	for _, group := range groups {
		col.checkCommentGroup(group)
	}
}

func (col *issueCollector) checkCommentGroup(group *ast.CommentGroup) {
	col.checkGeneric(group.Text(), "Comment", group.Pos())
}

func (col *issueCollector) checkDecls(decls []ast.Decl) {
	for _, decl := range decls {
		col.checkDecl(decl)
	}
}

func (col *issueCollector) checkDecl(decl ast.Decl) {
	if decl == nil {
		return
	}
	switch typedDecl := decl.(type) {
	case *ast.GenDecl:
		col.checkGenDecl(typedDecl)
	case *ast.FuncDecl:
		col.checkFuncDecl(typedDecl)
	}
}

func (col *issueCollector) checkGenDecl(typedDecl *ast.GenDecl) {
	col.checkSpecs(typedDecl.Specs)
}

func (col *issueCollector) checkSpecs(specs []ast.Spec) {
	for _, spec := range specs {
		col.checkSpec(spec)
	}
}

func (col *issueCollector) checkSpec(spec ast.Spec) {
	if spec == nil {
		return
	}
	switch typedSpec := spec.(type) {
	case *ast.ImportSpec: // TODO: local package name
	case *ast.ValueSpec: // TODO: names
	case *ast.TypeSpec:
		col.checkType(typedSpec)
	}
}

func (col *issueCollector) checkType(typeSpec *ast.TypeSpec) {
	col.checkIdent(typeSpec.Name, "Type name")
	col.checkTypeExpr(typeSpec.Type)
}

func (col *issueCollector) checkTypeExpr(typeExpr ast.Expr) {
	if typeExpr == nil {
		return
	}
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
		col.checkField(field, prefix)
	}
}

func (col *issueCollector) checkField(field *ast.Field, prefix string) {
	col.checkIdents(field.Names, prefix)
	col.checkTypeExpr(field.Type)
}

func (col *issueCollector) checkFuncDecl(funcDecl *ast.FuncDecl) {
	col.checkIdent(funcDecl.Name, "Function name")
	col.checkFieldList(funcDecl.Recv, "Function receiver")
	col.checkFuncType(funcDecl.Type)
	col.checkBlockStmt(funcDecl.Body)
}

func (col *issueCollector) checkBlockStmt(block *ast.BlockStmt) {
	if block == nil {
		return
	}
	col.checkStmts(block.List)
}

func (col *issueCollector) checkStmts(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		col.checkStmt(stmt)
	}
}

func (col *issueCollector) checkStmt(stmt ast.Stmt) {
	if stmt == nil {
		return
	}
	switch typedStmt := stmt.(type) {
	case *ast.DeclStmt:
		col.checkDecl(typedStmt.Decl)
	case *ast.LabeledStmt:
		col.checkLabelStmt(typedStmt)
	case *ast.ExprStmt:
		col.checkExprStmt(typedStmt)
	case *ast.SendStmt:
	case *ast.IncDecStmt:
	case *ast.AssignStmt:
	case *ast.GoStmt:
	case *ast.DeferStmt:
	case *ast.ReturnStmt:
	case *ast.BranchStmt:
	case *ast.BlockStmt:
		col.checkBlockStmt(typedStmt)
	case *ast.IfStmt:
		col.checkIfStmt(typedStmt)
	case *ast.CaseClause:
		col.checkCaseClause(typedStmt)
	case *ast.SwitchStmt:
		col.checkSwitchStmt(typedStmt)
	case *ast.TypeSwitchStmt:
		col.checkTypeSwitchStmt(typedStmt)
	case *ast.CommClause:
		col.checkCommClause(typedStmt)
	case *ast.SelectStmt:
		col.checkSelectStmt(typedStmt)
	case *ast.ForStmt:
		col.checkForStmt(typedStmt)
	case *ast.RangeStmt:
		col.checkRangeStmt(typedStmt)
	}
}

func (col *issueCollector) checkLabelStmt(stmt *ast.LabeledStmt) {
	col.checkIdent(stmt.Label, "Label")
	col.checkStmt(stmt.Stmt)
}

func (col *issueCollector) checkExprStmt(stmt *ast.ExprStmt) {
	col.checkExpr(stmt.X)
}

func (col *issueCollector) checkExprs(exprs []ast.Expr) {
	for _, expr := range exprs {
		col.checkExpr(expr)
	}
}

func (col *issueCollector) checkExpr(expr ast.Expr) {
	if expr == nil {
		return
	}
	switch typedStmt := expr.(type) {
	case *ast.Ident:
	case *ast.Ellipsis:
	case *ast.BasicLit:
	case *ast.FuncLit:
		col.checkFuncLit(typedStmt)
	case *ast.CompositeLit:
	case *ast.ParenExpr:
	case *ast.SelectorExpr:
	case *ast.IndexExpr:
	case *ast.SliceExpr:
	case *ast.TypeAssertExpr:
	case *ast.CallExpr:
	case *ast.StarExpr:
	case *ast.UnaryExpr:
	case *ast.BinaryExpr:
	case *ast.KeyValueExpr:
	}
}

func (col *issueCollector) checkIfStmt(stmt *ast.IfStmt) {
	col.checkStmt(stmt.Init)
	col.checkBlockStmt(stmt.Body)
	col.checkStmt(stmt.Else)
}

func (col *issueCollector) checkCaseClause(stmt *ast.CaseClause) {
	col.checkExprs(stmt.List)
	col.checkStmts(stmt.Body)
}

func (col *issueCollector) checkSwitchStmt(stmt *ast.SwitchStmt) {
	col.checkStmt(stmt.Init)
	col.checkExpr(stmt.Tag)
	col.checkBlockStmt(stmt.Body)
}

func (col *issueCollector) checkTypeSwitchStmt(stmt *ast.TypeSwitchStmt) {
	col.checkStmt(stmt.Init)
	col.checkStmt(stmt.Assign)
	col.checkBlockStmt(stmt.Body)
}

func (col *issueCollector) checkCommClause(stmt *ast.CommClause) {
	col.checkStmt(stmt.Comm)
	col.checkStmts(stmt.Body)
}

func (col *issueCollector) checkSelectStmt(stmt *ast.SelectStmt) {
	col.checkBlockStmt(stmt.Body)
}

func (col *issueCollector) checkForStmt(stmt *ast.ForStmt) {
	col.checkStmt(stmt.Init)
	col.checkExpr(stmt.Cond)
	col.checkStmt(stmt.Post)
	col.checkBlockStmt(stmt.Body)
}

func (col *issueCollector) checkRangeStmt(stmt *ast.RangeStmt) {
	col.checkExpr(stmt.Key)
	col.checkExpr(stmt.Value)
	col.checkExpr(stmt.X)
	col.checkBlockStmt(stmt.Body)
}

func (col *issueCollector) checkFuncLit(funcLit *ast.FuncLit) {
	col.checkFuncType(funcLit.Type)
	col.checkBlockStmt(funcLit.Body)
}
