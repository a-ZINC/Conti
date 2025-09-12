package container

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/a-ZINC/conti/vm"
	"github.com/google/uuid"
)

type Container struct {
	Name     string
	Pid      int
	Image    string
	Command  string
	Id       uuid.UUID
	Cmd      *exec.Cmd
	Provider vm.VMProvider
}

func NewContainer(name, image, command string) *Container {
	id := uuid.New()
	return &Container{
		Name:    name,
		Image:   image,
		Command: command,
		Id:      id,
	}
}

func (c *Container) ExecuteCommand(cmd string) (string, error) {
	switch c.Provider {
	case vm.ProviderWSL:
		return c.openWSLContainer(cmd)
	case vm.ProviderLima:
		return c.openLimaContainer(cmd)
	case vm.ProviderLinux:
		return c.openLinuxContainer(cmd)
	default:
		return "", fmt.Errorf("unsupported provider for Container execution: %v", c.Provider)
	}
}

func (c *Container) openWSLContainer(cmd string) (string, error) {
	wslCmd := exec.Command("wsl.exe", "-e", "bash", "-c", cmd)
	c.Cmd = wslCmd
	return c.streamingOutput(wslCmd)
}

func (c *Container) openLimaContainer(cmd string) (string, error) {
	limaCmd := exec.Command("limactl", "Container", "conti", "--", "bash", "-c", cmd)
	c.Cmd = limaCmd
	return c.streamingOutput(limaCmd)
}

func (c *Container) openLinuxContainer(cmd string) (string, error) {
	linuxCmd := exec.Command("bash", "-c", cmd)
	c.Cmd = linuxCmd
	return c.streamingOutput(linuxCmd)
}

func (c *Container) streamingOutput(cmd *exec.Cmd) (string, error) {
	var buff bytes.Buffer

	stdout := io.MultiWriter(os.Stdout, &buff)
	stderr := io.MultiWriter(os.Stderr, &buff)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Start()
	if err != nil {
		return "", err
	}
	return buff.String(), err
}

func (c *Container) Stop() error {
	return c.Cmd.Process.Kill()
}

