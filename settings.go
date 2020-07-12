package goconsider

// Settings contain all the parameters for the analysis.
type Settings struct {
	Phrases []Phrase
}

// Phrase describes an expression, with optional alternatives, that the linter flags.
type Phrase struct {
	// Synonyms are one or more expressions the linter looks for.
	Synonyms []string
	// Alternatives are zero, one, or more expressions that are provided as replacement.
	Alternatives []string
	// References are one or more resources that can help understand why the phrase is flagged.
	// Ideally, a reference is a stable link.
	References []string
}

// DefaultSettings return a settings instance for common use.
func DefaultSettings() Settings {
	settings := minimalSettings()
	settings = forEnglish(settings)
	return settings
}

func minimalSettings() Settings {
	return Settings{}
}

func forEnglish(settings Settings) Settings {
	return settings
}
