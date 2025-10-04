package main

import (
	"fmt"

	"github.com/a-ZINC/conti/client"
	// "github.com/a-ZINC/conti/container/container"
	"github.com/a-ZINC/conti/container/manager"
	"github.com/a-ZINC/conti/vm"
)

type AppConfig struct {
	VMEnabled bool
	Manager manager.ContainerManager
}

var RootFs = "/var/lib/conti"

func main() {
	fmt.Println("Hello, World!")
	// Initialize VM Manager
	vmManager := vm.NewVMManager()
	if err := vmManager.Start(); err != nil {
		fmt.Printf("Error starting VM Manager: %v\n", err)
		return
	}


	shell := vm.NewShell(vmManager.GetProvider())
	vmManager.Shell = shell
	// output, err := shell.ExecuteCommand("echo Hello from VM shell $USER")
	// if err != nil {
	// 	fmt.Printf("Error executing command in VM shell: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Command output: %s\n", output)

	// Ensure runtime is available in the VM
	dependency := vm.NewDependency(vmManager)
	err := dependency.InstallEssentialToolInVM()
	if err != nil {
		fmt.Printf("bruh bitch %v", err)
		return
	}


	containerManager := manager.NewContainerManager(shell)
	// container := container.NewContainer("bro", "bruh", "for i in {1..10}; do echo \"$USER - $i\"; sleep 2; done")
	// containerManager.AddContainer(container)
	// idStr := container.Id.String()
	// go containerManager.Run(idStr)
	// time.Sleep(3 * time.Second)
	// err = containerManager.Stop(idStr)
	// if err != nil {
	// 	return
	// }
	

	client := client.NewClient(containerManager)
	client.Start()
}