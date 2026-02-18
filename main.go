package main

import (
	"fmt"
	"os"

	"github.com/ersanisk/sieve/cmd"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	if err := cmd.Execute(version, buildTime); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
