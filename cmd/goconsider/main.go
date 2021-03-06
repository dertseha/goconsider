package main

import (
	"fmt"
	"os"
)

func main() {
	err := run(os.Stdout, os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
