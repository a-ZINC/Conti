package filesytem

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

type RootFileSystem struct {
	RootFS      string
	TimeoutTime time.Duration
}

func NewRootFileSystem(rootfs string) *RootFileSystem {
	return &RootFileSystem{
		RootFS:      rootfs,
		TimeoutTime: 2 * time.Minute,
	}
}

func (rfs *RootFileSystem) copyBinary(binaryPath string) error {
	containerBinaryPath := filepath.Join(rfs.RootFS, binaryPath)
	if err := rfs.copyFile(binaryPath, containerBinaryPath); err != nil {
		fmt.Printf("error copying file bro check it out: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), rfs.TimeoutTime)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ldd", binaryPath)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		var libPath string
		if len(fields) >= 3 && fields[1] == "=>" {
			libPath = fields[2]
		} else if len(fields) >= 1 && strings.HasPrefix(fields[0], "/") {
			libPath = fields[0]
		}
		if libPath != "" && libPath != "not" {
			destPath := filepath.Join(rfs.RootFS, libPath)
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("failed creating lib dir: %w", err)
			}
			if err := rfs.copyFile(libPath, destPath); err != nil {
				return fmt.Errorf("failed copying lib %s: %w", libPath, err)
			}
		}
	}

	fmt.Printf("copied binary and its dependencies\n")
	return nil
}

func (rfs *RootFileSystem) copyFile(binaryPath, containerBinaryPath string) error {
	binaryFile, err := os.Create(containerBinaryPath)
	if err != nil {
		return err
	}
	defer binaryFile.Close()

	hostBinaryFile, err := os.Open(binaryPath)
	if err != nil {
		return err
	}
	defer hostBinaryFile.Close()

	_, err = io.Copy(binaryFile, hostBinaryFile)
	if err != nil {
		return err
	}

	if err := os.Chmod(containerBinaryPath, 0755); err != nil {
		return err
	}
	return nil
}

func (rfs *RootFileSystem) essentialRootFilesystem() error {
	dirList := []string{
		"bin", "etc", "lib", "lib64", "usr", "proc", "sys", "dev", "tmp", "var", "home", "root", "usr/bin",
	}

	for _, dir := range dirList {
		dirPath := filepath.Join(rfs.RootFS, dir)
		if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
			continue
		}
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Printf("error occurred creating %s: %v\n", dir, err)
		}
	}

	essentailBinary := []string{
		"/bin/bash",
		"/bin/ls",
		"/bin/cat",
		"/bin/mount",
		"/bin/umount",
		"/bin/ps",
		"/bin/sleep",
		"/bin/ip",
		"/bin/nsenter",
		"/bin/python3.13",
		"/usr/bin/head",
		"/usr/bin/tr",
		"/usr/local/",
	}
	for _, binary := range essentailBinary {
		containerBinaryPath := filepath.Join(rfs.RootFS, binary)
		if _, err := os.Stat(containerBinaryPath); !os.IsNotExist(err) {
			fmt.Printf("Already available essential binary %s to container\n", binary)
			continue
		}
		err := rfs.copyBinary(binary)
		if err != nil {
			fmt.Printf("bro somewhere in creating FS we fucked up \n")
		}
	}
	return nil

}

func (rfs *RootFileSystem) PermissionEssentialFile() error {
	psswd := `
	root:x:0:0:root:/root:/bin/bash
	nobody:x:99:99:nobody:/:/bin/false
	`
	psswdFilePath := filepath.Join(rfs.RootFS, "etc", "passwd")
	if err := os.WriteFile(psswdFilePath, []byte(psswd), 0644); err != nil {
		fmt.Printf("fucked up writing passwd file %v /n", err)
	}

	group := `
	root:x:0:
	nobody:x:99:
	`
	groupFilePath := filepath.Join(rfs.RootFS, "etc", "group")
	if err := os.WriteFile(groupFilePath, []byte(group), 0644); err != nil {
		fmt.Printf("fucked up writing group file %v\n", err)
	}
	return nil
}

func (rfs *RootFileSystem) CreateTestFiles() error {
	// testHostFiles := []string{
	// 	// "/Users/ajinkya.singh/Desktop/go-project/conti/runtime/sh/cpu_hog.sh",
	// 	// "/Users/ajinkya.singh/Desktop/go-project/conti/runtime/sh/mem_hog.sh",
	// }

	// for _, file := range testHostFiles {
	// 	destPath := filepath.Join(rfs.RootFS, "root", filepath.Base(file))
	// 	if err := rfs.copyFile(file, destPath); err != nil {
	// 		fmt.Printf("error copying test file %s: %v\n", file, err)
	// 	}
	// }

	return nil
}

func (rfs *RootFileSystem) CreateDeviceFiles() error {
	deviceFiles := map[string][2]uint32{
		"/dev/null":    {1, 3},
		"/dev/zero":    {1, 5},
		"/dev/random":  {1, 8},
		"/dev/urandom": {1, 9},
	}

	for dev, numbers := range deviceFiles {
		devPath := dev
		if _, err := os.Stat(devPath); os.IsNotExist(err) {
			devNum := unix.Mkdev(numbers[0], numbers[1])
			if err := unix.Mknod(devPath, uint32(os.ModeCharDevice|0666), int(devNum)); err != nil {
				fmt.Printf("error creating device file %s: %v\n", devPath, err)
			} else {
				fmt.Printf("created device file: %s\n", devPath)
			}
		} else {
			fmt.Printf("device file already exists: %s\n", devPath)
		}
	}

	return nil
}

func (rfs *RootFileSystem) CreateMinimalRootfs() error {
	// install essential folder and file
	if err := rfs.essentialRootFilesystem(); err != nil {
		return err
	}

	if err := rfs.PermissionEssentialFile(); err != nil {
		return err
	}

	if err := rfs.CreateDeviceFiles(); err != nil {
		fmt.Printf("Error creating device files: %v\n", err)
		return err
	}

	// if err := rfs.CreateTestFiles(); err != nil {
	// 	return err
	// }

	return nil
}

func (rfs *RootFileSystem) IsFileSystemPresent() bool {
	_, err := os.Stat(rfs.RootFS)
	return os.IsExist(err)
}
