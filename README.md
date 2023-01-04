# goconsider

[![Go Report Card](https://goreportcard.com/badge/github.com/dertseha/goconsider)](https://goreportcard.com/report/github.com/dertseha/goconsider)
[![Lint Status](https://github.com/dertseha/goconsider/workflows/golangci-lint/badge.svg)](https://github.com/dertseha/goconsider/actions)

`goconsider` is a linter for [Go](https://golang.org) that proposes different words or phrases found in identifiers or comments.

The tool considers comments, filenames, and any identifier that is free to be chosen.
For example, it will raise an issue for the name of a declared type, but not if the code uses such a type.

It comes with a default set of phrases to support an inclusive language.  

## Example

```
code:
type MasterIndex int

output:
file.go:1:6: Type name contains 'master', consider rephrasing to one of ['primary', 'leader', 'main'].
```

## Install

### Manual download

Download the pre-compiled binaries from the [releases page](https://github.com/dertseha/goconsider/releases) and
copy them to the desired location.

### Via Go

```sh
go install github.com/dertseha/goconsider/cmd/goconsider@latest
```

> This puts the binary into `GOPATH/bin`.

## Run

```sh
cd go-project-dir
goconsider ./...
```

### Usage
```
> goconsider --help
goconsider: proposes alternatives for words or phrases found in source

Usage: goconsider [-flag] [package]

Flags:
  ... (several flags supported by go/analysis) 
  -settings string
        name of a settings file (defaults to '.goconsider.yaml' in current working directory)
  ...
```

## Configuration

### Default
#### Implicit
The tool will look for a `.goconsider.yaml` file in the current working directory.
See "explicit" configuration, below, for an example of its format.

If no such file exists, then the internal defaults will be used. 

#### Internal
The tool comes with a list of English phrases that are considered inappropriate.
It also proposes alternatives. See file [`default.yaml`](pkg/settings/default.yaml).

### Explicit
Command argument `-settings <filename>` will load the given `YAML` file.
Example: 
```
references:
  guide: https://example.com/guide
  dsl: https://example.com/dsl
  req: https://example.com/requirements

formatting:
  # By default false, a setting of true causes the long references to be printed for each issue.
  withReferences: true

phrases:
  - synonyms: [unwanted, variant]
    alternatives: [better, also good]
    references: [guide]
  - synonyms: [not good, worse]
    alternatives: [only this]
    references: [dsl, req]
```

## Algorithm

The algorithm is simple, yet effective enough to handle most likely cases.

The tool considers comments and identifier (names) that the developer has control over and can change.

First, the tool removes all punctuation from texts (in case of comments), as well as any casing.
This also separates CamelCase words, and the tool tries to keep abbreviations as one word.
A block of comment is considered as one long text. 

For example, the following texts all result in "this is an example" for further processing:
```
ThisIsAnExample
This-is-an-example
this is. An Example

as well as

// this is
// an
// example
```

The settings then specify which phrases to look for. Phrases allow looking for "word combinations".
So, the phrase `bad thing maker` will be found in identifiers such as`theBadThingMakerErr`,
or a comment like `The bad thing maker does stuff`.

## Recommendations

### References for phrases

For provided configurations, it is not required that phrases have any reference.
However, having references helps to trace to the origin why a phrase is found.
This could also be achieved by commit messages, yet these are not printed in the result list.

## Limits

* There is no ignore system for "false positives". This could be handled by using a linter framework, such as `golangci-lint`.
* The word-finding algorithm is simple and can probably be tricked. If someone uses this tool *and* circumvents it this way, it's not an issue of the tool.
* There is no concept of automatic singular/plural detection. For such phrases, provide additional variants as synonyms.

## License

MIT license, see [license file](LICENSE).
