package main

import (
	"fmt"

	"github.com/a-ZINC/conti/vm"
)

func main() {
	fmt.Println("Hello, World!")
	// Initialize VM Manager
	vmManager := vm.NewVMManager()
	vmManager.Start()
}