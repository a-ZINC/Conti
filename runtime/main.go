//go:build linux

package main

import (
	"os"

	"github.com/a-ZINC/conti/runtime/pkg/runner"
)

func main() {
	run := runner.NewRunner("/rootfs")
	switch os.Args[1] {
	case "run":
		run.CreateContainerProcess()
	case "init":
		run.ExecuteContainerProcess()
	default:
		run.ExecuteContainerProcess()
	}
}
