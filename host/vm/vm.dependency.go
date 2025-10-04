package vm

import (
	"fmt"
	"strings"
)

type Dependency struct {
	manager *VMManager
}

const (
	Git = "git"
	Go  = "go"
)

func NewDependency(manager *VMManager) *Dependency {
	return &Dependency{
		manager: manager,
	}
}

func (d *Dependency) isExist(tool string) (bool, error) {
	test := fmt.Sprintf("which %s", tool)
	output, err := d.manager.Shell.ExecuteCommandInVM(test)
	if err != nil {
		return false, err
	}
	return strings.Contains(output, tool), nil
}

func (d *Dependency) installGo() error {
	arch := d.manager.config.arch
	if d.manager.config.arch == "aarch64" {
		arch = "arm64"
	}
	goBinary := fmt.Sprintf("go1.21.2.%s-%s.tar.gz", d.manager.config.osName, arch)
	goBinaryURL := fmt.Sprintf("https://golang.org/dl/%s", goBinary)
	wgetStr := fmt.Sprintf("cd /tmp && wget %s", goBinaryURL)
	_, err := d.manager.Shell.ExecuteCommandInVM(wgetStr)
	if err != nil {
		return err
	}

	tarStr := fmt.Sprintf("cd /tmp && tar -xf %s", goBinary)
	_, err = d.manager.Shell.ExecuteCommandInVM(tarStr)
	if err != nil {
		return err
	}

	_, err = d.manager.Shell.ExecuteCommandInVM("cd /tmp && sudo mv go /usr/local/go")
	if err != nil {
		return err
	}

	_, err = d.manager.Shell.ExecuteCommandInVM(`echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bashrc`)
	if err != nil {
		return err
	}

	_, err = d.manager.Shell.ExecuteCommandInVM("source ~/.bashrc")
	if err != nil {
		return err
	}

	return nil
}

func (d *Dependency) InstallEssentialToolInVM() error {

	exist, err := d.isExist(Go)
	fmt.Printf("exist: %t err: %v \n", exist, err)
	if err != nil {
		return err
	}
	if !exist {
		fmt.Printf("Go doesnt exist so downloading it \n")
		err := d.installGo()
		if err != nil {
			return err
		}
	}

	



	return nil
}
