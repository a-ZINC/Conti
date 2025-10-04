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

const (
	GoVersion = "1.21.2"
)

const (
	conti = "https://github.com/a-ZINC/Conti.git"
)

func NewDependency(manager *VMManager) *Dependency {
	return &Dependency{
		manager: manager,
	}
}

func (d *Dependency) installGo() error {

	checkStr := `[ -x /usr/local/go/bin/go ] && echo 'exists' || echo 'nope'`
	output, err := d.manager.Shell.ExecuteCommandInVM(checkStr)
	if err != nil {
		return err
	}

	if strings.Contains(output, "exists") {
		fmt.Printf("go already installed \n")
		return nil
	}

	arch := d.manager.config.arch
	if d.manager.config.arch == "aarch64" {
		arch = "arm64"
	}
	goBinary := fmt.Sprintf("go%s.%s-%s.tar.gz", GoVersion,  d.manager.config.osName, arch)
	goBinaryURL := fmt.Sprintf("https://golang.org/dl/%s", goBinary)
	wgetStr := fmt.Sprintf("cd /tmp && wget %s", goBinaryURL)
	_, err = d.manager.Shell.ExecuteCommandInVM(wgetStr)
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

	//inside lima not working need to work on it

	// _, err = d.manager.Shell.ExecuteCommandInVM(`echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.bashrc`)
	// if err != nil {
	// 	return err
	// }

	// _, err = d.manager.Shell.ExecuteCommandInVM("source ~/.bashrc")
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (d *Dependency) installGit() error {
	checkStr := "which git"
	out, err := d.manager.Shell.ExecuteCommandInVM(checkStr)
	if err != nil {
		return err
	}
	if strings.Contains(out, "git") {
		fmt.Printf("git already installed \n")
		return nil
	}
	_, err = d.manager.Shell.ExecuteCommandInVM("sudo apt update")
	if err != nil {
		return err
	}
	_, err = d.manager.Shell.ExecuteCommandInVM("sudo apt update")
	if err != nil {
		return err
	}

	_, err = d.manager.Shell.ExecuteCommandInVM("sudo apt install -y git")
	if err != nil {
		return err
	}

	return nil
}

func (d *Dependency) installRuntime() error {
	_, err := d.manager.Shell.ExecuteCommandInVM(fmt.Sprintf("cd /tmp && git clone %s", conti))
	if err != nil {
		return err
	}

	_, err = d.manager.Shell.ExecuteCommandInVM("cd /tmp/Conti/runtime && /usr/local/go/bin/go build -o runtime")
	if err != nil {
		return err
	}

	_, err = d.manager.Shell.ExecuteCommandInVM(" sudo mv /tmp/Conti/runtime/runtime /usr/local/runtime")
	if err != nil {
		return err
	}

	// _, err = d.manager.Shell.ExecuteCommandInVM(`sed -i '1i export PATH=/usr/local/runtime:$PATH' ~/.bashrc`)
	// if err != nil {
	// 	return err
	// }

	// _, err = d.manager.Shell.ExecuteCommandInVM("source ~/.bashrc")
	// if err != nil {
	// 	return err
	// }
	return nil

}

func (d *Dependency) isRuntimeExist() (bool, error) {
	existStr := `[ -x /usr/local/runtime ] && echo 'exists' || echo 'none'`
	output, err := d.manager.Shell.ExecuteCommandInVM(existStr)
	if err != nil {
		return false, err
	}

	return strings.Contains(output, "exists"), nil
}

func (d *Dependency) InstallEssentialToolInVM() error {

	exist, err := d.isRuntimeExist()
	if err != nil {
		return err
	}
	if exist {
		fmt.Printf("runtime is already present \n")
		return nil
	}

	err = d.installGo()
	if err != nil {
		return err
	}

	err = d.installGit()
	if err != nil {
		return err
	}

	err = d.installRuntime()
	if err != nil {
		return err
	}

	return nil
}
