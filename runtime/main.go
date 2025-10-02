//go:build linux

package main

import (
	"fmt"
	"os"

	"github.com/a-ZINC/conti/runtime/pkg/runner"
)

func main() {
	run := runner.NewRunner("/rootfs")
	fmt.Printf("Runtime started with PID %d\n", os.Getpid())
	switch os.Args[1] {
	case "run":
		run.CreateContainerProcess()
	case "init":
		run.ExecuteContainerProcess()
	default:
		run.ExecuteContainerProcess()
	}
}
