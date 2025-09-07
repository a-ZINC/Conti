package main

import (
	"fmt"

	"github.com/a-ZINC/conti/vm"
)

type AppConfig struct {
	VMEnabled bool
}

func main() {
	fmt.Println("Hello, World!")
	// Initialize VM Manager
	vmManager := vm.NewVMManager()
	if err := vmManager.Start(); err != nil {
		fmt.Printf("Error starting VM Manager: %v\n", err)
		return
	}
	shell := vm.NewShell(vmManager.GetProvider())
	output, err := shell.ExecuteCommand("echo Hello from VM shell $USER")
	if err != nil {
		fmt.Printf("Error executing command in VM shell: %v\n", err)
		return
	}
	fmt.Printf("Command output: %s\n", output)
}