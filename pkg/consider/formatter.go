package consider

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

type referenceModel struct {
	// Short refers to the key that is directly associated with the finding.
	Short string
	// Long is a resolved string identified by the short string. Empty if not found.
	Long string
}

type formatModel struct {
	// Context is describing where the triggering phrase was found.
	Context string
	// Found is the triggering phrase.
	Found string
	// Alternatives is the list of possibilities that can replace the phrase.
	Alternatives []string
	// References is the list of sources for the reasoning.
	References []referenceModel

	// PrintReferences is true if the long form shall be added to the message.
	PrintReferences bool
}

type formatter struct {
	templ *template.Template
}

//go:embed default_format.gotext
var defaultFormat string

func newFormatter() *formatter {
	funcs := template.FuncMap{"join": strings.Join}
	templ, err := template.New("message").Funcs(funcs).Parse(defaultFormat)
	if err != nil {
		panic(fmt.Sprintf("failed to parse built-in format: %v", err))
	}
	return &formatter{templ: templ}
}

func (f *formatter) Format(model formatModel) string {
	buf := bytes.NewBuffer(nil)
	err := f.templ.Execute(buf, model)
	if err != nil {
		return fmt.Sprintf("failed to format message: %v", err)
	}
	return buf.String()
}
