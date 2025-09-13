//go:build linux

package main

import (
	"os"

	"github.com/a-ZINC/conti/runtime/pkg"
)

func main() {
	switch os.Args[1] {
	case "run":
		pkg.CreateContainerProcess()
	case "init":
		pkg.ExecuteContainerProcess()
	default:
		pkg.ExecuteContainerProcess()
	}
}
