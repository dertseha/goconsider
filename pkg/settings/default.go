package settings

import (
	_ "embed"

	"github.com/dertseha/goconsider/pkg/consider"
)

//go:embed default.yaml
var defaultYaml []byte

// Default returns the default settings of the linter package.
// They are reasonable for common use, yet they may also change over time with different versions.
//
// If you only operate with default settings, it may be that newer versions of the linter will suddenly report
// new issues in old code.
func Default() consider.Settings {
	s, err := FromYaml(defaultYaml)
	if err != nil {
		panic("embedded default settings cannot be parsed! This is a packaging issue.")
	}
	return s
}
