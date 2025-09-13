//go:build linux

package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func ExecuteContainerProcess() {
	fmt.Printf("Container process: %d\n", os.Getpid())
	if err := syscall.Sethostname([]byte("container")); err != nil {
		fmt.Printf("Error setting hostname: %v\n", err)
		return
	}
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %v\n", err)
		return
	}
	fmt.Printf("Command %s finished\n", os.Args[2])
}

func CreateContainerProcess() {
	cmd := exec.Command("/proc/self/exe", append([]string{"init"}, os.Args[2:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	fmt.Printf("Started process with PID %d (parent PID: %d)\n", cmd.Process.Pid, os.Getpid())

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for command: %v\n", err)
		return
	}
	fmt.Printf("Process %d exited\n", cmd.Process.Pid)
}
