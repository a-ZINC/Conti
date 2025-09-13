package filesytem

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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
		fmt.Printf("error copying file bro check it")
	}

	ctx, cancel := context.WithTimeout(context.Background(), rfs.TimeoutTime)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ldd", binaryPath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("copied binary and its dependencies")
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
	return nil

}

func (rfs *RootFileSystem) essentialRootFilesystem() error {
	dirList := []string{
		"bin", "etc", "lib", "lib64", "usr", "proc", "sys", "dev", "tmp", "var", "home",
	}

	for _, dir := range dirList {
		dirPath := filepath.Join(rfs.RootFS, dir)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Printf("error occured creating %s", dir)
		}
	}

	essentailBinary := []string{
		"/bin/bash",
		"/bin/ls",
		"/bin/cat",
		"/bin/mount",
		"/bin/umount",
		"/bin/ps",
	}
	for _, binary := range essentailBinary {
		err := rfs.copyBinary(binary)
		if err != nil {
			fmt.Printf("bro somewhere in creating FS we fucked up")
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
		fmt.Printf("fucked up writing passwd file")
	}

	group := `
	root:x:0:
	nobody:x:99:
	`
	groupFilePath := filepath.Join(rfs.RootFS, "etc", "group")
	if err := os.WriteFile(groupFilePath, []byte(group), 0644); err != nil {
		fmt.Printf("fucked up writing group file")
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

	return nil
}

func (rfs *RootFileSystem) IsFileSystemPresent() bool {
	_, err := os.Stat(rfs.RootFS)
	return os.IsExist(err)
}
