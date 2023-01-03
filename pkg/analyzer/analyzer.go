package analyzer

import (
	"go/token"
	"os"
	"path"

	"github.com/dertseha/goconsider/pkg/consider"
	"github.com/dertseha/goconsider/pkg/settings"
	"golang.org/x/tools/go/analysis"
)

const (
	implicitSettingsFilename = ".goconsider.yaml"

	analyzerName  = "goconsider"
	documentation = "proposes alternatives for words or phrases found in source"
)

// NewAnalyzer returns a new instance with the given settings.
func NewAnalyzer(s consider.Settings) *analysis.Analyzer {
	an := newBaseAnalyzer()
	an.Run = runnerWithSettingsFrom(func() (consider.Settings, error) { return s, nil })
	return an
}

// NewAnalyzerFromSettingsFile returns a new instance that will load the settings from a file at given path.
// If the given string is empty, defaults will apply.
func NewAnalyzerFromSettingsFile(settingsFile string) *analysis.Analyzer {
	an := newBaseAnalyzer()
	an.Run = runnerWithSettingsFrom(func() (consider.Settings, error) { return resolveSettings(settingsFile) })
	return an
}

// NewAnalyzerFromFlags returns an instance that defers to configuration via flags.
func NewAnalyzerFromFlags() *analysis.Analyzer {
	an := newBaseAnalyzer()
	settingsFile := an.Flags.String("settings", "",
		"name of a settings file (defaults to '"+implicitSettingsFilename+"' in current working directory)")
	an.Run = runnerWithSettingsFrom(func() (consider.Settings, error) { return resolveSettings(*settingsFile) })
	return an
}

func newBaseAnalyzer() *analysis.Analyzer {
	an := &analysis.Analyzer{
		Name: analyzerName,
		Doc:  documentation,
	}
	an.Flags.Var(versionFlag{}, "V", "print version and exit")
	return an
}

func runnerWithSettingsFrom(factory func() (consider.Settings, error)) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		s, err := factory()
		if err != nil {
			return nil, err
		}
		return run(s, pass)
	}
}

func resolveSettings(settingsFile string) (consider.Settings, error) {
	if len(settingsFile) != 0 {
		return readSettings(settingsFile)
	}
	return defaultSettings()
}

func defaultSettings() (consider.Settings, error) {
	filename := path.Join(".", implicitSettingsFilename)
	if _, err := os.Stat(filename); err != nil {
		return settings.Default(), nil
	}
	return readSettings(filename)
}

func readSettings(settingsFile string) (consider.Settings, error) {
	settingsData, err := os.ReadFile(settingsFile)
	if err != nil {
		return consider.Settings{}, err
	}
	s, err := settings.FromYaml(settingsData)
	if err != nil {
		return consider.Settings{}, err
	}
	return s, nil
}

func run(settings consider.Settings, pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		linter := consider.NewLinter(settings, reporterFuncFor(pass))
		linter.CheckFile(f, pass.Fset.File(f.Package))
	}
	return nil, nil
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
