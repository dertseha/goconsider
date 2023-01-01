package consider

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/dertseha/goconsider/internal/text"
)

// Reporter is the outgoing interface for detected issues.
type Reporter interface {
	// Report is called for each detected issue.
	Report(pos token.Pos, message string)
}

// Linter is the main type of the linting functionality.
type Linter struct {
	settings  Settings
	formatter *formatter
	reporter  Reporter

	issuesSuppressed bool
}

// NewLinter returns a new instance for given parameters.
func NewLinter(settings Settings, reporter Reporter) *Linter {
	return &Linter{
		settings:         settings,
		formatter:        newFormatter(),
		reporter:         reporter,
		issuesSuppressed: false,
	}
}

// CheckFile runs the analysis on given file.
func (l *Linter) CheckFile(file *ast.File, rawFile *token.File) {
	l.issuesSuppressed = false

	l.checkFilename(file, rawFile)
	l.checkIdent(file.Name, "Package name")
	l.checkCommentGroups(file.Comments)
	l.checkDecls(file.Decls)
}

func (l *Linter) suppressIssues(on bool) func() {
	currentSuppression := l.issuesSuppressed
	l.issuesSuppressed = on
	return func() { l.issuesSuppressed = currentSuppression }
}

func (l *Linter) addIssue(typeString string, pos token.Pos, synonym string, phrase Phrase) {
	if l.issuesSuppressed {
		return
	}
	l.reporter.Report(pos, l.formatMessage(typeString, synonym, phrase))
}

func (l *Linter) checkGeneric(s string, typeString string, pos token.Pos) {
	worded := text.Wordify(s)
	for _, phrase := range l.settings.Phrases {
		for _, synonym := range phrase.Synonyms {
			if strings.Contains(worded, " "+synonym+" ") {
				l.addIssue(typeString, pos, synonym, phrase)
			}
		}
	}
}

func (l *Linter) checkIdents(idents []*ast.Ident, prefix string) {
	for _, ident := range idents {
		l.checkIdent(ident, prefix)
	}
}

func (l *Linter) checkIdent(ident *ast.Ident, typeString string) {
	if ident == nil {
		return
	}
	l.checkGeneric(ident.Name, typeString, ident.NamePos)
}

func (l *Linter) checkFilename(file *ast.File, rawFile *token.File) {
	if rawFile == nil {
		return
	}
	_, filename := filepath.Split(rawFile.Name())
	l.checkGeneric(filename, "File name", file.Package)
}

func (l *Linter) checkCommentGroups(groups []*ast.CommentGroup) {
	for _, group := range groups {
		l.checkCommentGroup(group)
	}
}

func (l *Linter) checkCommentGroup(group *ast.CommentGroup) {
	l.checkGeneric(group.Text(), "Comment", group.Pos())
}

func (l *Linter) checkDecls(decls []ast.Decl) {
	for _, decl := range decls {
		l.checkDecl(decl)
	}
}

func (l *Linter) checkDecl(decl ast.Decl) {
	if decl == nil {
		return
	}
	reset := l.suppressIssues(false)
	switch typedDecl := decl.(type) {
	case *ast.GenDecl:
		l.checkGenDecl(typedDecl)
	case *ast.FuncDecl:
		l.checkFuncDecl(typedDecl)
	}
	reset()
}

func (l *Linter) checkGenDecl(typedDecl *ast.GenDecl) {
	l.checkSpecs(typedDecl.Specs)
}

func (l *Linter) checkSpecs(specs []ast.Spec) {
	for _, spec := range specs {
		l.checkSpec(spec)
	}
}

func (l *Linter) checkSpec(spec ast.Spec) {
	if spec == nil {
		return
	}
	switch typedSpec := spec.(type) {
	case *ast.ImportSpec:
		l.checkImportSpec(typedSpec)
	case *ast.ValueSpec:
		l.checkValueSpec(typedSpec)
	case *ast.TypeSpec:
		l.checkType(typedSpec)
	}
}

func (l *Linter) checkImportSpec(spec *ast.ImportSpec) {
	l.checkIdent(spec.Name, "Package alias")
}

func (l *Linter) checkValueSpec(spec *ast.ValueSpec) {
	l.checkIdents(spec.Names, "Value name")
}

func (l *Linter) checkType(spec *ast.TypeSpec) {
	l.checkIdent(spec.Name, "Type name")
	l.checkTypeExpr(spec.Type)
}

func (l *Linter) checkTypeExpr(typeExpr ast.Expr) {
	if typeExpr == nil {
		return
	}
	switch spec := typeExpr.(type) {
	case *ast.StructType:
		l.checkFieldList(spec.Fields, "Member name")
	case *ast.FuncType:
		l.checkFuncType(spec)
	case *ast.InterfaceType:
		l.checkFieldList(spec.Methods, "Method name")
	}
}

func (l *Linter) checkFuncType(funcType *ast.FuncType) {
	l.checkFieldList(funcType.Params, "Parameter name")
	l.checkFieldList(funcType.Results, "Result name")
}

func (l *Linter) checkFieldList(fields *ast.FieldList, prefix string) {
	if fields == nil {
		return
	}
	for _, field := range fields.List {
		l.checkField(field, prefix)
	}
}

func (l *Linter) checkField(field *ast.Field, prefix string) {
	l.checkIdents(field.Names, prefix)
	l.checkTypeExpr(field.Type)
}

func (l *Linter) checkFuncDecl(funcDecl *ast.FuncDecl) {
	l.checkIdent(funcDecl.Name, "Function name")
	l.checkFieldList(funcDecl.Recv, "Function receiver")
	l.checkFuncType(funcDecl.Type)
	l.checkBlockStmt(funcDecl.Body)
}

func (l *Linter) checkBlockStmt(block *ast.BlockStmt) {
	if block == nil {
		return
	}
	l.checkStmts(block.List)
}

func (l *Linter) checkStmts(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		l.checkStmt(stmt)
	}
}

func (l *Linter) checkStmt(stmt ast.Stmt) {
	if stmt == nil {
		return
	}
	switch typedStmt := stmt.(type) {
	case *ast.DeclStmt:
		l.checkDecl(typedStmt.Decl)
	case *ast.LabeledStmt:
		l.checkLabelStmt(typedStmt)
	case *ast.ExprStmt:
		l.checkExprStmt(typedStmt)
	case *ast.SendStmt:
	case *ast.IncDecStmt:
	case *ast.AssignStmt:
		l.checkAssignStmt(typedStmt)
	case *ast.GoStmt:
	case *ast.DeferStmt:
	case *ast.ReturnStmt:
	case *ast.BranchStmt:
	case *ast.BlockStmt:
		l.checkBlockStmt(typedStmt)
	case *ast.IfStmt:
		l.checkIfStmt(typedStmt)
	case *ast.CaseClause:
		l.checkCaseClause(typedStmt)
	case *ast.SwitchStmt:
		l.checkSwitchStmt(typedStmt)
	case *ast.TypeSwitchStmt:
		l.checkTypeSwitchStmt(typedStmt)
	case *ast.CommClause:
		l.checkCommClause(typedStmt)
	case *ast.SelectStmt:
		l.checkSelectStmt(typedStmt)
	case *ast.ForStmt:
		l.checkForStmt(typedStmt)
	case *ast.RangeStmt:
		l.checkRangeStmt(typedStmt)
	}
}

func (l *Linter) checkLabelStmt(stmt *ast.LabeledStmt) {
	l.checkIdent(stmt.Label, "Label")
	l.checkStmt(stmt.Stmt)
}

func (l *Linter) checkExprStmt(stmt *ast.ExprStmt) {
	l.checkExpr(stmt.X)
}

func (l *Linter) checkAssignStmt(stmt *ast.AssignStmt) {
	reset := l.suppressIssues(stmt.Tok != token.DEFINE)
	l.checkExprs(stmt.Lhs)
	reset()

	reset = l.suppressIssues(true)
	l.checkExprs(stmt.Rhs)
	reset()
}

func (l *Linter) checkExprs(exprs []ast.Expr) {
	for _, expr := range exprs {
		l.checkExpr(expr)
	}
}

func (l *Linter) checkExpr(expr ast.Expr) {
	if expr == nil {
		return
	}
	switch typedStmt := expr.(type) {
	case *ast.Ident:
		l.checkIdent(typedStmt, "Identifier")
	case *ast.Ellipsis:
	case *ast.BasicLit:
	case *ast.FuncLit:
		l.checkFuncLit(typedStmt)
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

func (l *Linter) checkIfStmt(stmt *ast.IfStmt) {
	l.checkStmt(stmt.Init)
	l.checkBlockStmt(stmt.Body)
	l.checkStmt(stmt.Else)
}

func (l *Linter) checkCaseClause(stmt *ast.CaseClause) {
	l.checkExprs(stmt.List)
	l.checkStmts(stmt.Body)
}

func (l *Linter) checkSwitchStmt(stmt *ast.SwitchStmt) {
	l.checkStmt(stmt.Init)
	l.checkExpr(stmt.Tag)
	l.checkBlockStmt(stmt.Body)
}

func (l *Linter) checkTypeSwitchStmt(stmt *ast.TypeSwitchStmt) {
	l.checkStmt(stmt.Init)
	l.checkStmt(stmt.Assign)
	l.checkBlockStmt(stmt.Body)
}

func (l *Linter) checkCommClause(stmt *ast.CommClause) {
	l.checkStmt(stmt.Comm)
	l.checkStmts(stmt.Body)
}

func (l *Linter) checkSelectStmt(stmt *ast.SelectStmt) {
	l.checkBlockStmt(stmt.Body)
}

func (l *Linter) checkForStmt(stmt *ast.ForStmt) {
	l.checkStmt(stmt.Init)
	l.checkExpr(stmt.Cond)
	l.checkStmt(stmt.Post)
	l.checkBlockStmt(stmt.Body)
}

func (l *Linter) checkRangeStmt(stmt *ast.RangeStmt) {
	l.checkExpr(stmt.Key)
	l.checkExpr(stmt.Value)
	l.checkExpr(stmt.X)
	l.checkBlockStmt(stmt.Body)
}

func (l *Linter) checkFuncLit(funcLit *ast.FuncLit) {
	reset := l.suppressIssues(false)
	l.checkFuncType(funcLit.Type)
	l.checkBlockStmt(funcLit.Body)
	reset()
}

func (l *Linter) formatMessage(context, found string, phrase Phrase) string {
	model := formatModel{
		Context:         context,
		Found:           found,
		Alternatives:    phrase.Alternatives,
		ShortReferences: phrase.References,
		References:      nil,

		PrintReferences: (l.settings.Formatting.WithReferences != nil) && *l.settings.Formatting.WithReferences,
	}
	for _, short := range model.ShortReferences {
		long := l.settings.References[short]
		if len(long) > 0 {
			model.References = append(model.References, long)
		} else {
			model.References = append(model.References, short)
		}
	}
	return l.formatter.Format(model)
}
