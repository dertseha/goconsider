package version

// The following "constants" need to be initialized at build-time, through ldflags.
var scmTag string
var scmCommit string

// Info describes the version of the library.
type Info struct {
	// CoreWithPreRelease is the combination of the three numeric identifiers and an optional pre-release suffix.
	CoreWithPreRelease string
	// Build is the auxiliary information for the version.
	Build string
}

// Version returns the version specified at build time.
func Version() Info {
	coreWithPreRelease := scmTag
	if len(coreWithPreRelease) == 0 {
		coreWithPreRelease = "0.0.0-alpha"
	}
	build := scmCommit
	if len(build) == 0 {
		build = "manual"
	}
	return Info{
		CoreWithPreRelease: coreWithPreRelease,
		Build:              build,
	}
}
