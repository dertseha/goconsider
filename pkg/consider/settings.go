package consider

// Settings contain all the parameters for the analysis.
type Settings struct {
	// References is a key-value map of short keys to a reference, typically a stable link.
	// They indicate resources that can help understand why phrases are flagged, or
	// at least give examples of other (larger) peer groups that considered rewording.
	References map[string]string `yaml:"references"`
	// Phrases describe all the texts the linter should look for.
	Phrases []Phrase `yaml:"phrases"`
	// Formatting describes how the messages shall be formatted.
	Formatting Formatting `yaml:"formatting"`
}

// Phrase describes an expression, with optional alternatives, that the linter flags.
type Phrase struct {
	// Synonyms are one or more expressions that have the same meaning and proposed alternatives.
	Synonyms []string `yaml:"synonyms"`
	// Alternatives are zero, one, or more expressions that are provided as replacement.
	Alternatives []string `yaml:"alternatives"`
	// References is a list of either direct, or keyed references into the global map of references.
	References []string `yaml:"references"`
}

// Formatting descries how messages shall be formatted.
type Formatting struct {
	// WithReferences indicates whether the long-form of references shall be added.
	// This is not done by default as this is done in separate lines.
	WithReferences *bool `yaml:"withReferences"`
}
