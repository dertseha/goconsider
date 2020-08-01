package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dertseha/goconsider"
)

type arguments struct {
	help      bool
	filenames []string
}

type issuesFoundError int

func (err issuesFoundError) Error() string {
	return fmt.Sprintf("%v issues were found", int(err))
}

func run(out io.Writer, rawArgs []string) error {
	args, err := parseArguments(rawArgs)
	if err != nil {
		return err
	}
	if args.help {
		printUsage(out)
		return nil
	}

	settings := goconsider.DefaultSettings()

	fset, files, err := parseFiles(args.filenames)
	if err != nil {
		return err
	}
	issueCount := lintAndReport(out, fset, files, settings)
	if issueCount > 0 {
		return issuesFoundError(issueCount)
	}
	return nil
}

func printUsage(out io.Writer) {
	const usage = `Usage:
godot [OPTIONS] [FILES]
Options:
	-h, --help      show this message
`
	_, _ = fmt.Fprint(out, usage)
}

type unknownArgumentErr string

func (err unknownArgumentErr) Error() string {
	return fmt.Sprintf("unknown argument '%s'", string(err))
}

func parseArguments(rawArgs []string) (arguments, error) {
	var args arguments

	for _, arg := range rawArgs {
		if !strings.HasPrefix(arg, "-") {
			args.filenames = append(args.filenames, arg)
			continue
		}

		switch arg {
		case "-h", "--help":
			args.help = true
		default:
			return arguments{}, unknownArgumentErr(arg)
		}
	}
	return args, nil
}

type pathOrErr struct {
	filepath string
	err      error
}

func parseFiles(filenames []string) (*token.FileSet, []*ast.File, error) {
	var files []*ast.File
	fset := token.NewFileSet()
	for _, filename := range filenames {
		paths := allGoFilesIn(filename)
		for p := range paths {
			if p.err != nil {
				return nil, nil, p.err
			}

			file, err := parser.ParseFile(fset, p.filepath, nil, parser.ParseComments)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse file '%s': %w", p.filepath, err)
			}
			files = append(files, file)
		}
	}
	return fset, files, nil
}

func allGoFilesIn(root string) chan pathOrErr {
	paths := make(chan pathOrErr)

	go func() {
		defer close(paths)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if hasVendorPath(path) {
				return nil
			}
			if err != nil {
				paths <- pathOrErr{err: err}
				return nil
			}
			if isAGoFile(info) {
				paths <- pathOrErr{filepath: path}
			}
			return nil
		})
		if err != nil {
			paths <- pathOrErr{err: err}
		}
	}()

	return paths
}

func hasVendorPath(s string) bool {
	const sep = string(filepath.Separator)
	return strings.HasPrefix(s, "vendor"+sep) || strings.Contains(s, sep+"vendor"+sep)
}

func isAGoFile(info os.FileInfo) bool {
	return !info.IsDir() && strings.HasSuffix(info.Name(), ".go")
}

func lintAndReport(out io.Writer, fset *token.FileSet, files []*ast.File, settings goconsider.Settings) int {
	issueCount := 0
	for _, file := range files {
		issues := goconsider.Lint(file, fset, settings)
		issueCount += len(issues)
		for _, issue := range issues {
			_, _ = fmt.Fprintf(out, "%s: %s\n", issue.Pos, issue.Message)
		}
	}
	return issueCount
}
