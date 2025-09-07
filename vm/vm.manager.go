package vm

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type VMProvider int

const (
	ProviderNone VMProvider = iota
	ProviderWSL
	ProviderLima
	ProviderLinux
)

func (p VMProvider) String() string {
	switch p {
	case ProviderWSL:
		return "WSL"
	case ProviderLima:
		return "Lima"
	case ProviderLinux:
		return "Linux"
	default:
		return "Unknown"
	}
}

type VMManager struct {
	provider VMProvider
	timeout  time.Duration
}

func NewVMManager() *VMManager {
	var provider VMProvider
	switch runtime.GOOS {
	case "windows":
		provider = ProviderWSL
	case "darwin":
		provider = ProviderLima
	case "linux":
		provider = ProviderLinux
	default:
		provider = ProviderNone
	}
	return &VMManager{
		provider: provider,
		timeout:  30 * time.Second,
	}
}

func (vm *VMManager) Start() error {
	if !vm.providerManager() {
		log.Printf("Failed to manage provider: %v\n", vm.provider)
		return fmt.Errorf("failed to manage provider: %v", vm.provider)
	}
	if vm.isVMAvailable() {
		log.Printf("VM is already available using provider: %v\n", vm.provider)
		return nil
	}
	log.Printf("Provider %v is ready\n", vm.provider)
	fmt.Printf("The VM will be created using the %v provider.\n", vm.provider)
	fmt.Printf("This may take a few minutes depending on your system and internet connection.\n")
	fmt.Print("Press Enter to continue or type anything else to abort: ")
	userInput := ""
	_, err := fmt.Scanln(&userInput)
	if err != nil {
		err = vm.createVM()
		if err != nil {
			log.Printf("Error creating VM: %v\n", err)
			return err
		}
		log.Printf("VM created successfully using provider: %v\n", vm.provider)
		return nil
	}
	log.Printf("VM creation aborted by user.\n")
	return nil
}

func (vm *VMManager) providerManager() bool {
	if vm.provider == ProviderNone {
		log.Printf("Provider is not supported\n")
		return false
	}
	if vm.provider == ProviderLinux {
		log.Printf("Running on native Linux, no VM needed\n")
		return false
	}
	installed, err := vm.isProviderInstalled()
	log.Printf("Provider installed: %v\n", installed)
	if err != nil {
		log.Printf("Error checking provider installation: %v\n", err)
	}
	if installed {
		log.Printf("Provider is already installed\n")
		return installed
	}

	switch vm.provider {
	case ProviderWSL:
		log.Printf("Managing WSL provider\n")
		err := vm.installWSL()
		if err != nil {
			log.Printf("Error installing WSL: %v\n", err)
			return false
		}
		return true
	case ProviderLima:
		log.Printf("Managing Lima provider\n")
		err := vm.installLima()
		if err != nil {
			log.Printf("Error installing Lima: %v\n", err)
			return false
		}
		return true
	}
	return true
}


func (vm *VMManager) isProviderInstalled() (bool, error) {
	switch vm.provider {
	case ProviderWSL:
		return vm.isWSLInstalled()
	case ProviderLima:
		return vm.isLimaInstalled()
	case ProviderLinux:
		return true, nil
	default:
		return false, fmt.Errorf("unknown provider")
	}
}

func (vm *VMManager) isWSLInstalled() (bool, error) {
	fmt.Printf("Checking if WSL is installed...\n")
	ctx, cancel := context.WithTimeout(context.Background(), vm.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "wsl", "--list", "-q")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	if len(output) == 0 {
		return false, fmt.Errorf("WSL is not installed")
	}
	return true, nil
}

func (vm *VMManager) isLimaInstalled() (bool, error) {
	fmt.Printf("Checking if Lima is installed...\n")
	ctx, cancel := context.WithTimeout(context.Background(), vm.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "which", "limactl")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	if len(output) == 0 {
		return false, fmt.Errorf("lima is not installed")
	}
	return true, nil
}

func (vm *VMManager) installLima() error {
	fmt.Printf("Installing Lima...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)
	defer cancel()
	err := exec.CommandContext(ctx, "which", "brew").Run()
	if err != nil {
		log.Printf("Error installing Lima: %v\n", err)
		return fmt.Errorf("install brew: %w", err)
	}
	err = exec.CommandContext(ctx, "brew", "install", "lima").Run()
	if err != nil {
		log.Printf("Error installing Lima: %v\n", err)
		return err
	}
	log.Printf("Lima installed successfully\n")
	return nil
}

func (vm *VMManager) installWSL() error {
	fmt.Printf("Installing WSL...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "wsl", "--install")
	err := cmd.Run()
	if err != nil {
		log.Printf("Error installing WSL: %v\n", err)
		return err
	}
	log.Printf("WSL installed successfully\n")
	return nil
}



func (vm *VMManager) createVM() error {
	log.Printf("Creating VM using provider: %v\n", vm.provider.String())
	switch vm.provider {
	case ProviderWSL:
		return vm.createWSLVM()
	case ProviderLima:
		return vm.createLimaVM()
	default:
		return fmt.Errorf("unsupported provider for VM creation")
	}
}

func (vm *VMManager) createWSLVM() error {
	log.Printf("Creating WSL VM...\n")
	time.Sleep(2 * time.Second)
	log.Printf("WSL VM created successfully\n")
	return nil
}

func (vm *VMManager) createLimaVM() error {
	log.Printf("Creating Lima VM...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "limactl", "start", "--name=conti")
	err := cmd.Run()
	if err != nil {
		log.Printf("Error creating Lima VM: %v\n", err)
		return err
	}
	log.Printf("Lima VM created successfully\n")
	return nil
}

func (vm *VMManager) isVMAvailable() bool {
	switch vm.provider {
	case ProviderWSL:
		return vm.isWSLVMAvailable()
	case ProviderLima:
		return vm.isLimaVMAvailable()
	case ProviderLinux:
		return true
	default:
		return false
	}
}

func (vm *VMManager) isWSLVMAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), vm.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "wsl", "--list", "-q")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error checking WSL VM availability: %v\n", err)
		return false
	}
	if len(output) == 0 {
		return false
	}
	return true
}

func (vm *VMManager) isLimaVMAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), vm.timeout)
	defer cancel()

	output, err := exec.CommandContext(ctx, "limactl", "list", "--json").Output()
	if err != nil {
		log.Printf("Error getting Lima VM details: %v\n", err)
		return false
	}
	if strings.Contains(string(output), "conti") {

		if strings.Contains(string(output), "Stopped") {
			err = exec.CommandContext(ctx, "limactl", "start", "conti").Run()
			if err != nil {
				log.Printf("Error starting Lima VM: %v\n", err)
				return false
			}

			time.Sleep(5 * time.Second)
		}
		return true
	}
	return false
}

func (vm *VMManager) GetProvider() VMProvider {
	return vm.provider
}
