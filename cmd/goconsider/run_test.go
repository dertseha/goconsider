package main // nolint: testpackage

import (
	"bytes"
	"os"
	"path"
	"testing"
)

func TestReportIssuesWithDefaultSettings(t *testing.T) {
	out := bytes.NewBuffer(nil)
	err := run(out, []string{path.Join(".", "testdata", "issues", "default.go")})
	if err == nil {
		t.Errorf("error expected")
	} else if _, isOK := err.(issuesFoundError); !isOK {
		t.Errorf("unexpected error returned: %v", err)
	}
}

func TestReportIssuesWithExplicitSettings(t *testing.T) {
	out := bytes.NewBuffer(nil)
	err := run(out, []string{
		"--settings", path.Join(".", "testdata", "settings", "explicit.yaml"),
		path.Join(".", "testdata", "issues", "explicit.go"),
	})
	if err == nil {
		t.Errorf("error expected")
	} else if _, isOK := err.(issuesFoundError); !isOK {
		t.Errorf("unexpected error returned: %v", err)
	}
}

func TestReportIssuesWithImplicitSettings(t *testing.T) {
	out := bytes.NewBuffer(nil)
	err := os.Chdir(path.Join(".", "testdata", "settings", "implicit"))
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
	}
	err = run(out, []string{
		path.Join("..", "..", "issues", "implicit.go"),
	})
	if err == nil {
		t.Errorf("error expected")
	} else if _, isOK := err.(issuesFoundError); !isOK {
		t.Errorf("unexpected error returned: %v", err)
	}
}
