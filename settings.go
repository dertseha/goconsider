package goconsider

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

// DefaultSettings return a settings instance for common use.
func DefaultSettings() Settings {
	settings := minimalSettings()
	settings = forEnglish(settings)
	return settings
}

func minimalSettings() Settings {
	return Settings{}
}

func synonyms(s ...string) func(Phrase) Phrase {
	return func(p Phrase) Phrase {
		p.Synonyms = append(p.Synonyms, s...)
		return p
	}
}

func alternatives(s ...string) func(Phrase) Phrase {
	return func(p Phrase) Phrase {
		p.Alternatives = append(p.Alternatives, s...)
		return p
	}
}

func references(s ...string) func(Phrase) Phrase {
	return func(p Phrase) Phrase {
		p.References = append(p.References, s...)
		return p
	}
}

func phraseWith(mod ...func(Phrase) Phrase) Phrase {
	var p Phrase
	for _, m := range mod {
		p = m(p)
	}
	return p
}

func forEnglish(settings Settings) Settings {
	add := func(p Phrase) {
		settings.Phrases = append(settings.Phrases, p)
	}

	const linuxKernel = "https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=49decddd39e5f6132ccd7d9fdc3d7c470b0061bb"
	const twitter = "https://www.cnet.com/news/twitter-engineers-replace-racially-loaded-tech-terms-like-master-slave/"
	const google = "https://developers.google.com/style/inclusive-documentation"
	const googlePronouns = "https://developers.google.com/style/pronouns#gender-neutral-pronouns"

	add(phraseWith(synonyms("master"),
		alternatives("primary", "leader", "main"),
		references(linuxKernel, twitter),
	))
	add(phraseWith(synonyms("slave"),
		alternatives("secondary", "follower", "replica", "standby"),
		references(linuxKernel, twitter),
	))

	add(phraseWith(synonyms("whitelist"),
		alternatives("allowlist", "passlist"),
		references(linuxKernel, twitter),
	))
	add(phraseWith(synonyms("blacklist"),
		alternatives("denylist", "blocklist"),
		references(linuxKernel, twitter),
	))

	add(phraseWith(synonyms("grandfathered"),
		alternatives("legacy status"),
		references(twitter),
	))

	add(phraseWith(synonyms("guys"),
		alternatives("people", "folks", "you all"),
		references(twitter),
	))
	add(phraseWith(synonyms("he", "his", "him"),
		alternatives("their", "them"),
		references(googlePronouns, twitter),
	))
	add(phraseWith(synonyms("man hours"),
		alternatives("person hours", "engineer hours"),
		references(google, twitter),
	))

	add(phraseWith(synonyms("dummy"),
		alternatives("placeholder", "sample"),
		references(google, twitter),
	))
	add(phraseWith(synonyms("sanity check"),
		alternatives("quick check"),
		references(google, twitter),
	))
	return settings
}
