package network

import (
	"fmt"
	"os/exec"
	"strings"
)

type NetworkManager struct {
	BridgeName        string
	BridgeIP          string
	VethContainerName string
	VethBridgeName    string
	ContainerIP       string
}

var containerCount = 2

func generateContainerIP() string {
	containerCount++
	return fmt.Sprintf("28.28.28.%d/24", containerCount)
}

func NewNetworkManager(bridgename string, containerName string) *NetworkManager {
	return &NetworkManager{
		BridgeName:        bridgename,
		BridgeIP:          "28.28.28.1/24",
		VethContainerName: "veth0",
		VethBridgeName:    fmt.Sprintf("veth-%s", containerName),
		ContainerIP:       generateContainerIP(),
	}
}

func (nm *NetworkManager) SetupNetworking() error {
	if err := nm.SetupBridge(); err != nil {
		fmt.Printf("network bridge err: %v \n", err)
		return err
	}
	if err := nm.SetupVethPair(); err != nil {
		fmt.Printf("network veth err: %v \n", err)
		return err
	}
	return nil
}

func (nm *NetworkManager) SetupBridge() error {
	cmd := exec.Command("which", "ip")
	output, err := cmd.Output()
	fmt.Printf("output %s \n", string(output))
	if err != nil {
		return fmt.Errorf("error interface list %v", err)
	}
	fmt.Printf("output %s \n", string(output))

	cmd = exec.Command("ip", "link", "list")
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("error interface list %v", err)
	}
	fmt.Printf("output: %v\n", string(output))
	outputStr := strings.Trim(string(output), " ")
	if strings.Contains(outputStr, nm.BridgeName) {
		fmt.Printf("Bridge already exist %s \n", nm.BridgeName)
		cmd = exec.Command("ip", "link", "set", nm.BridgeName, "up")
		err = cmd.Run()
		return err
	}

	cmd = exec.Command("ip", "link", "add", nm.BridgeName, "type", "bridge")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error bridge creation %v", err)
	}

	cmd = exec.Command("ip", "addr", "add", nm.BridgeIP, "dev", nm.BridgeName)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error bridge ip %v", err)
	}

	cmd = exec.Command("ip", "link", "set", nm.BridgeName, "up")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error bridge up %v", err)
	}
	return nil
}

func (nm *NetworkManager) SetInterfaceUp(netInterface string) error {
	cmd := exec.Command("ip", "link", "set", netInterface, "up")
	err := cmd.Run()
	return err
}

func (nm *NetworkManager) SetupVethPair() error {
	fmt.Printf("veth bridge: %s, veth container: %s\n", nm.VethBridgeName, nm.VethContainerName)
	cmd := exec.Command("ip", "link", "list")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error fetching available interfaces")
	}

	fmt.Printf("available interface: %s \n", output)

	cmd = exec.Command("ip", "link", "add", nm.VethBridgeName, "type", "veth", "peer", "name", nm.VethContainerName)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("veth creation error: %v", err)
	}

	cmd = exec.Command("ip", "link", "set", nm.VethBridgeName, "master", nm.BridgeName)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("veth master error")
	}
	cmd = exec.Command("ip", "link", "set", nm.VethBridgeName, "up")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("veth master error")
	}
	return nil
}

func (nm *NetworkManager) SetupContainerNetwork(pid int) error {
	pidStr := fmt.Sprintf("%d", pid)
	cmd := exec.Command("ip", "link", "set", nm.VethContainerName, "netns", pidStr)
	err := cmd.Run()
	if err != nil {
		return nil
	}

	cmd = exec.Command("nsenter", "-t", pidStr, "-n", "ip", "link", "set", "lo", "up")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("err setting up conatiner lo interface: %v", err)
	}

	cmd = exec.Command("nsenter", "-t", pidStr, "-n", "ip", "addr", "add", nm.ContainerIP, "dev", nm.VethContainerName)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("err setting up conatiner veth interface: %v", err)
	}

	cmd = exec.Command("nsenter", "-t", pidStr, "-n", "ip", "link", "set", nm.VethContainerName, "up")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("err setting up conatiner veth interface: %v", err)
	}

	cmd = exec.Command("nsenter", "-t", pidStr, "-n", "ip", "route", "add", "default", "via", strings.Split(nm.BridgeIP, "/")[0])
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("err setting up conatiner default route: %v", err)
	}

	return nil
}
