//go:build linux

package main

import (
	"os"

	"github.com/a-ZINC/conti/runtime/pkg/runner"
)

func main() {
	switch os.Args[1] {
	case "run":
		runner.CreateContainerProcess()
	case "init":
		runner.ExecuteContainerProcess()
	default:
		runner.ExecuteContainerProcess()
	}
}
