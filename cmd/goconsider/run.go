package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/dertseha/goconsider"
)

type arguments struct {
	help         bool
	noReferences bool
	filenames    []string
	settings     string
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

	var settings goconsider.Settings
	if len(args.settings) > 0 {
		settings, err = parseSettings(args.settings)
	} else {
		settings, err = defaultSettings()
	}
	if err != nil {
		return err
	}

	fset, files, err := parseFiles(reportTo(out), args.filenames)
	if err != nil {
		return err
	}
	issueCount := lintAndReport(out, reportTo(out), fset, files, settings, !args.noReferences)
	if issueCount > 0 {
		return issuesFoundError(issueCount)
	}
	return nil
}

const implicitSettingsFilename = ".goconsider.yaml"

func defaultSettings() (goconsider.Settings, error) {
	filename := path.Join(".", implicitSettingsFilename)
	if _, err := os.Stat(filename); err != nil {
		return goconsider.DefaultSettings(), nil
	}
	return parseSettings(filename)
}

func parseSettings(filename string) (goconsider.Settings, error) {
	settingsData, err := ioutil.ReadFile(filename) // nolint: gosec
	if err != nil {
		return goconsider.Settings{}, err
	}
	var settings goconsider.Settings
	err = yaml.Unmarshal(settingsData, &settings)
	if err != nil {
		return goconsider.Settings{}, err
	}
	return settings, nil
}

func printUsage(out io.Writer) {
	const usage = `Usage:
goconsider [OPTIONS] [FILES]
Options:
	-h, --help             Show this message
	--settings <filename>  Name of a settings file. Defaults to '` + implicitSettingsFilename + `' in current working directory.
	--noReferences         Skip printing references as per settings for any found issues.
`
	_, _ = fmt.Fprint(out, usage)
}

type unknownArgumentErr string

func (err unknownArgumentErr) Error() string {
	return fmt.Sprintf("unknown argument '%s'", string(err))
}

type missingParameterErr string

func (err missingParameterErr) Error() string {
	return fmt.Sprintf("argument '%s' is missing a parameter", string(err))
}

func parseArguments(rawArgs []string) (arguments, error) {
	var args arguments

	for i := 0; i < len(rawArgs); i++ {
		arg := rawArgs[i]
		if !strings.HasPrefix(arg, "-") {
			args.filenames = append(args.filenames, arg)
			continue
		}

		switch arg {
		case "-h", "--help":
			args.help = true
		case "--settings":
			i++
			if i >= len(rawArgs) {
				return arguments{}, missingParameterErr(arg)
			}
			args.settings = rawArgs[i]
		case "--noReferences":
			args.noReferences = true
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

func parseFiles(report issueReporter, filenames []string) (*token.FileSet, []*ast.File, error) {
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
				report(token.Position{
					Filename: p.filepath,
					Offset:   0,
					Line:     1,
					Column:   1,
				}, "failed to parse file")
				continue
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

func lintAndReport(out io.Writer, reporter issueReporter,
	fset *token.FileSet, files []*ast.File,
	settings goconsider.Settings, reportRef bool) int {
	issueCount := 0
	references := make(map[string]struct{})
	for _, file := range files {
		issues := goconsider.Lint(file, fset, settings)
		issueCount += len(issues)
		for _, issue := range issues {
			reporter(issue.Pos, issue.Message)
			for _, ref := range issue.References {
				references[ref] = struct{}{}
			}
		}
	}
	if reportRef && len(references) > 0 {
		refList := make([]string, 0, len(references))
		for ref := range references {
			refList = append(refList, ref)
		}
		sort.StringSlice(refList).Sort()
		_, _ = fmt.Fprintf(out, "References:\n")
		for _, ref := range refList {
			_, _ = fmt.Fprintf(out, "%s\n", ref)
		}
	}
	return issueCount
}

type issueReporter func(pos token.Position, message string)

func reportTo(out io.Writer) issueReporter {
	return func(pos token.Position, message string) {
		_, _ = fmt.Fprintf(out, "%s: %s\n", pos, message)
	}
}
