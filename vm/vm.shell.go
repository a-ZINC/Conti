package vm

import (
	"fmt"
	"os/exec"
)

type Shell struct {
	provider VMProvider
}

func NewShell(provider VMProvider) *Shell {
	return &Shell{provider: provider}
}

func (s *Shell) ExecuteCommand(cmd string) (string, error) {
	switch s.provider {
	case ProviderWSL:
		return s.openWSLShell(cmd)
	case ProviderLima:
		return s.openLimaShell(cmd)
	case ProviderLinux:
		return s.openLinuxShell(cmd)
	default:
		return "", fmt.Errorf("unsupported provider for shell execution: %v", s.provider)
	}
}

func (s *Shell) openWSLShell(cmd string) (string, error) {
	wslCmd := exec.Command("wsl.exe", "-e", "bash", "-c", cmd)
	output, err := wslCmd.Output()
	return string(output), err
}

func (s *Shell) openLimaShell(cmd string) (string, error) {
	limaCmd := exec.Command("limactl", "shell", "conti", "--", "bash", "-c", cmd)
	output, err := limaCmd.Output()
	return string(output), err
}

func (s *Shell) openLinuxShell(cmd string) (string, error) {
	linuxCmd := exec.Command("bash", "-c", cmd)
	output, err := linuxCmd.Output()
	return string(output), err
}
