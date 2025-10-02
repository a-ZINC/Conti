//go:build linux

package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/a-ZINC/conti/runtime/pkg/controller"
	"github.com/a-ZINC/conti/runtime/pkg/filesytem"
	"github.com/a-ZINC/conti/runtime/pkg/network"
)

type Runner struct {
	RootFS string
}

func NewRunner(rootfs string) *Runner {
	return &Runner{
		RootFS: rootfs,
	}
}
func (r *Runner) setupFilesytem() error {
	filesys := filesytem.NewRootFileSystem(r.RootFS)
	if !filesys.IsFileSystemPresent() {
		err := filesys.CreateMinimalRootfs()
		if err != nil {
			fmt.Printf("Error creating minimal root filesystem: %v\n", err)
			return err
		}
	}
	if err := syscall.Mount("", "/", "", uintptr(syscall.MS_REC|syscall.MS_PRIVATE), ""); err != nil {
        fmt.Printf("warning: could not make mounts private: %v\n", err)
    }
	if err := syscall.Mount(r.RootFS, r.RootFS, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		fmt.Printf("Error mounting root filesystem: %v\n", err)
		return err
	}

	putOld := filepath.Join(r.RootFS, ".put_old")
	if err := os.MkdirAll(putOld, 0700); err != nil {
		return fmt.Errorf("creating put_old: %w", err)
	}

	dent, _ := os.ReadDir(putOld)
    if len(dent) != 0 {
        return fmt.Errorf("put_old must be empty, contains %d entries", len(dent))
    }

	if err := syscall.Mount(putOld, putOld, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		fmt.Printf("Error mounting put_old: %v\n", err)
		return err
	}

	if err := os.Chdir(r.RootFS); err != nil {
		return fmt.Errorf("chdir to new root: %w", err)
	}

	if err := syscall.PivotRoot(".", ".put_old"); err != nil {
		fmt.Printf("Error performing pivot root: %v\n", err)
		if data, rerr := os.ReadFile("/proc/self/mountinfo"); rerr == nil {
            fmt.Printf("mountinfo:\n%s\n", string(data))
        }
        return fmt.Errorf("pivot_root: %w", err)
	}
	if err := os.Chdir("/"); err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		return err
	}
	if err := syscall.Unmount("/.put_old", syscall.MNT_DETACH); err != nil {
		fmt.Printf("Error unmounting put_old: %v\n", err)
		return err
	}

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Printf("Error mounting proc: %v\n", err)
		return err
	}
	if err := syscall.Mount("sysfs", "/sys", "sysfs", 0, ""); err != nil {
		fmt.Printf("Error mounting sysfs: %v\n", err)
		return err
	}
	// if err := syscall.Mount("tmpfs", "/dev", "tmpfs", 0, ""); err != nil {
	// 	fmt.Printf("Error mounting tmpfs: %v\n", err)
	// 	return err
	// }
	return nil
}

func (r *Runner) ExecuteContainerProcess() {
	fmt.Printf("Container process: %d\n", os.Getpid())
	if err := syscall.Sethostname([]byte("container")); err != nil {
		fmt.Printf("Error setting hostname: %v\n", err)
		return
	}
	if err := r.setupFilesytem(); err != nil {
		fmt.Printf("Error setting file system %v\n", err)
	}
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running command: %v\n", err)
		return
	}
	fmt.Printf("Command %s finished\n", os.Args[2])
}

func (r *Runner) CreateContainerProcess() {
	cpu := 50
	mem := int64(100 * 1024 * 1024)
	containerName := os.Getenv("CONTI_CONTAINER_NAME")
	fmt.Printf("Creating container %s with PID %d\n", containerName, os.Getpid())

	networkmanager := network.NewNetworkManager("br0", containerName)
	if err := networkmanager.SetupNetworking(); err != nil {
		return
	}
	cmd := exec.Command("/proc/self/exe", append([]string{"init"}, os.Args[2:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	// fmt.Printf("Creating container process: %d\n", cmd.Process.Pid)

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}
	err := networkmanager.SetupContainerNetwork(cmd.Process.Pid)
	if err != nil {
		fmt.Printf("Error setting up container network: %v\n", err)
		return
	}
	name := fmt.Sprintf("container-%d", cmd.Process.Pid)
	controller := controller.NewController(name, cmd.Process.Pid, cpu, mem)
	controller.SetupCgroups()
	// go controller.LoopResourceLookup()
	fmt.Printf("Started process with PID %d (parent PID: %d)\n", cmd.Process.Pid, os.Getpid())

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for command: %v\n", err)
		return
	}
	fmt.Printf("Process %d exited\n", cmd.Process.Pid)
}
