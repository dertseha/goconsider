package consider

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

type formatModel struct {
	// Context is describing where the triggering phrase was found.
	Context string
	// Found is the triggering phrase.
	Found string
	// Alternatives is the list of possibilities that can replace the phrase.
	Alternatives []string
	// ShortReferences is the list of guidance for the reasoning, expressed as short keys.
	ShortReferences []string
	// References is the list of long references for the reasoning.
	References []string

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
	templ, _ := template.New("message").Funcs(funcs).Parse(defaultFormat)
	return &formatter{templ: templ}
}

func (f *formatter) Format(model formatModel) string {
	buf := bytes.NewBuffer(nil)
	_ = f.templ.Execute(buf, model)
	return buf.String()
}
