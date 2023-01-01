package consider

// Settings contain all the parameters for the analysis.
type Settings struct {
	// Phrases describe all the texts the linter should look for
	Phrases []Phrase
}

// Phrase describes an expression, with optional alternatives, that the linter flags.
type Phrase struct {
	// Synonyms are one or more expressions that have the same meaning and proposed alternatives.
	Synonyms []string
	// Alternatives are zero, one, or more expressions that are provided as replacement.
	Alternatives []string
	// References are one or more resources that can help understand why the phrase is flagged, or
	// at least give examples of other (larger) peer groups that considered rewording.
	// Ideally, a reference is a stable link.
	References []string
}
