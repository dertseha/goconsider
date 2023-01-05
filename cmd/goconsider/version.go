package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dertseha/goconsider/internal/version"
)

// versionFlag attempts to adhere to the (apparently unspecified) contract of the -V functionality.
//
// Resources:
// golang.org/x/tools/go/analysis/internal/analysisflags/flags.go adds -V, and provides comments about contract.
// go1.19.4/src/cmd/internal/objabi/flag.go includes further notes about how it should behave.
type versionFlag struct{}

func (versionFlag) IsBoolFlag() bool { return true }
func (versionFlag) Get() interface{} { return nil }
func (versionFlag) String() string   { return "" }
func (versionFlag) Set(s string) error {
	name := os.Args[0]
	name = name[strings.LastIndex(name, `/`)+1:]
	name = name[strings.LastIndex(name, `\`)+1:]
	name = strings.TrimSuffix(name, ".exe")

	ver := version.Version()
	fullSuffix := ""

	if s == "full" {
		fullSuffix += " buildID=" + ver.Build
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s version %s%s\n", name, ver.CoreWithPreRelease, fullSuffix)
	os.Exit(0)
	return nil
}
