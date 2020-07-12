package main

import (
	"fmt"
	"os"
)

func main() {
	err := run(os.Args[1:], os.Stdout)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
