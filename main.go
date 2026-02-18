package main

import (
	"github.com/ersanisk/sieve/cmd"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	rootCmd := cmd.NewRootCmd(version, buildTime)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
