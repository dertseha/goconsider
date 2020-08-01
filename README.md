# goconsider

`goconsider` is a linter for [Go](https://golang.org) that proposes different words or phrases
found in identifier or comments.

The tool considers comments, and any identifier that is free to be chosen.
For example, it will raise an issue for the name of a declared type, but not if the code uses such a type.

It comes with a default set of phrases to support an inclusive language.  

## Example

```
code:
type MasterIndex int

output:
file.go:1:6: Type name contains 'master', consider rephrasing to one of [primary, leader, main]
```

## Install

Build from source

```sh
go get -u github.com/dertseha/goconsider/cmd/goconsider
```

## Run

```sh
goconsider ./project-dir
```

### Usage
```
> goconsider --help
Usage:
goconsider [OPTIONS] [FILES]
Options:
        -h, --help             Show this message
        --settings <filename>  Name of a settings file. Defaults to '.goconsider' in current working directory.
        --noReferences         Skip printing references as per settings for any found issues.
```

## Configuration

### Default
#### Implicit
The tool will look for a `.goconsider.yaml` file in the current working directory.
See "explicit" configuration, below, for an example.

If no such file exists, then the internal defaults will be used. 

#### Internal
The tool comes with a list of English phrases that are considered inappropriate.
It also proposes alternatives. See file [`settings.go`](settings.go).

### Explicit
Command argument `--settings <filename>` will load the given `YAML` file.
Example: 
```
phrases:
  - synonyms: [unwanted, variant]
    alternatives: [better, also good]
    references: [https://example.com/guide]
  - synonyms: [not good, worse]
    alternatives: [only this]
    references: [https://example.com/dsl, https://example.com/requirements]
```

## Algorithm

The algorithm is simple, yet effective enough to handle most likely cases.

The tool considers comments and identifier that the user has control over.

First, the tool removes all punctuation from texts (in case of comments), as well as any casing.
This also separates CamelCase words, and the tool tries to keep abbreviations as one word.
A Block of comment is considered as one long text. 

For example, the following texts all result in "this is an example" for further processing:
```
ThisIsAnExample
This-is-an-example
this is. An Example
```

The settings then specify which phrases to look for. Phrases allow looking for "word combinations".
So, the phrase `bad thing maker` will be found in identifiers such as`theBadThingMakerErr`,
or a comment like `The bad thing maker does stuff`.


## TODO

* Consider integration in `golangci-lint`.
* Implement check for filenames themselves.

## Limits

* There is no ignore system for "false positives". This could be handled by using a linter framework, such as `golangci-lint`.
* The word-finding algorithm is simple and can probably be tricked. If someone uses this tool and circumvents it this way, it's not an issue of the tool.
* There is no concept of automatic singular/plural detection. For such phrases, provide additional variants as synonyms.

## License

MIT license, see [license file](LICENSE).
