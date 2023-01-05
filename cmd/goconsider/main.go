// Package main provides the main entry function for the standalone linter executable.
package main

import (
	"github.com/dertseha/goconsider/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	an := analyzer.NewAnalyzerFromFlags()
	an.Flags.Var(versionFlag{}, "V", "print version and exit")
	singlechecker.Main(an)
}
