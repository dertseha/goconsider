// Package main provides the main entry function for the standalone linter executable.
package main

import (
	"github.com/dertseha/goconsider/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzerFromFlags())
}
