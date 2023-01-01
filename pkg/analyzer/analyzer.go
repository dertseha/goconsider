package analyzer

import (
	"go/token"

	"github.com/dertseha/goconsider/pkg/consider"
	"golang.org/x/tools/go/analysis"
)

// NewAnalyzer returns a new instance with the given settings.
func NewAnalyzer(settings consider.Settings) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "goconsider",
		Doc:  "proposes alternatives for words or phrases found in source",
		Run:  func(pass *analysis.Pass) (interface{}, error) { return run(settings, pass) },
	}
}

type reporterFunc func(pos token.Pos, message string)

func (f reporterFunc) Report(pos token.Pos, message string) {
	f(pos, message)
}

func reporterFuncFor(pass *analysis.Pass) reporterFunc {
	return func(pos token.Pos, message string) {
		pass.Report(analysis.Diagnostic{
			Pos:     pos,
			Message: message,
		})
	}
}

func run(settings consider.Settings, pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		linter := consider.NewLinter(settings, reporterFuncFor(pass))
		linter.CheckFile(f, pass.Fset.File(f.Package))
	}
	return nil, nil
}
