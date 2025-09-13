//go:build linux

package pkg

import (
	"os"
	"os/exec"
	"syscall"
)

func CreateProcess() {
	cmd := exec.Command("/bin/bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
