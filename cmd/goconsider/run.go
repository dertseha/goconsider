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
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/dertseha/goconsider"
)

type arguments struct {
	help      bool
	filenames []string
	settings  string
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

const implicitSettingsFilename = ".goconsider.yaml"

func defaultSettings() (goconsider.Settings, error) {
	filename := path.Join(".", implicitSettingsFilename)
	if _, err := os.Stat(filename); err != nil {
		return goconsider.DefaultSettings(), nil
	}
	return parseSettings(filename)
}

func parseSettings(filename string) (goconsider.Settings, error) {
	settingsData, err := ioutil.ReadFile(filename)
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
	--settings <filename>  Name of a settings file. Defaults to '.goconsider' in current working directory.
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
