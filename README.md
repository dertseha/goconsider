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

## Configuration

### Default
By default, the tool comes with a list of English phrases that are considered inappropriate.
It also proposes alternatives. See file [`settings.go`](settings.go).

## TODO

* Better documentation on how the tool finds phrases.
* Figure out whether and how references should be reported.
* Consider integration in `golangci-lint`.
* Implement check for filenames themselves.

## Limits

* There is no ignore system for "false positives". This could be handled by using a linter framework, such as `golangci-lint`.
* The word-finding algorithm is simple and can probably be tricked. If someone uses this tool and circumvents it this way, it's not an issue of the tool. 

## License

MIT license, see [license file](LICENSE).
