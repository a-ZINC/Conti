package vm

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	return s.streamingOutput(wslCmd)
}

func (s *Shell) openLimaShell(cmd string) (string, error) {
	limaCmd := exec.Command("limactl", "shell", "conti", "--", "bash", "-c", cmd)
	return s.streamingOutput(limaCmd)
}

func (s *Shell) openLinuxShell(cmd string) (string, error) {
	linuxCmd := exec.Command("bash", "-c", cmd)
	return s.streamingOutput(linuxCmd)
}

func (s *Shell) streamingOutput(cmd *exec.Cmd) (string, error) {
	var buff bytes.Buffer

	stdout := io.MultiWriter(os.Stdout, &buff)
	stderr := io.MultiWriter(os.Stderr, &buff)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	} 
	return buff.String(), err
}