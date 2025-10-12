package filesytem

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func NewRootFileSystemAlpine(rootfs string) *RootFileSystem {
	return &RootFileSystem{
		RootFS:      rootfs,
		TimeoutTime: 2 * time.Minute,
	}
}

func runWithTimeout(timeout time.Duration, name string, args ...string) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (rfs *RootFileSystem) CreateMinimalRootfsAlpine() error {
	if err := os.Mkdir(rfs.RootFS, 0744); err != nil {
		return err
	}

	if err := runWithTimeout(rfs.TimeoutTime, "wget", "-P", "/tmp", "https://dl-cdn.alpinelinux.org/alpine/v3.20/releases/aarch64/alpine-minirootfs-3.20.1-aarch64.tar.gz"); err != nil {
		fmt.Println("error downloading alpine tar")
		return err
	}

	if err := runWithTimeout(rfs.TimeoutTime, "tar", "-xzf", "/tmp/alpine-minirootfs-3.20.1-aarch64.tar.gz", "-C", "/rootfs"); err != nil {
		fmt.Println("error tar extraction")
		return err
	}
	return nil
}
