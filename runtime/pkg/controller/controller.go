package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Controller struct {
	Name string
	Pid  int
	mem int64
	cpu int
}

func NewController(name string, pid, cpu int, mem int64) *Controller {
	return &Controller{
		Name: name,
		Pid:  pid,
		mem: mem,
		cpu: cpu,
	}
}

func (c *Controller) isCgroupV2() bool {
	groupV2Path := filepath.Join("/sys/fs/cgroup", "cgroup.controllers")
	if _, err := os.Stat(groupV2Path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (c *Controller) SetMemoryLimitV1(bytes int64) error {
	memoryPath := filepath.Join("/sys/fs/cgroup/memory", c.Name)
	if err := os.MkdirAll(memoryPath, 0644); err != nil {
		fmt.Printf("Bro unable to create memory folder for conatiner %s %v", c.Name, err)
		return err
	}
	memory_size_limit := filepath.Join(memoryPath, "memory.limit_in_bytes")
	if err := os.WriteFile(memory_size_limit, []byte(strconv.FormatInt(bytes, 10)), 0644); err != nil {
		fmt.Printf("Bro unable to create memory limit for conatiner %s err: %v", c.Name, err)
		return err
	}

	swap_memeory_limit := filepath.Join(memoryPath, "memory.memsw.limit_int_bytes")
	if err := os.WriteFile(swap_memeory_limit, []byte(strconv.FormatInt(bytes, 10)), 0644); err != nil {
		fmt.Printf("Bro unable to create swap memory limit for conatiner %s", c.Name)
	}

	return nil
}

func (c *Controller) SetMemoryLimitV2(bytes int64) error {
	memoryPath := filepath.Join("/sys/fs/cgroup", c.Name, "memory.max")
	if err := os.WriteFile(memoryPath, []byte(strconv.FormatInt(bytes, 10)), 0644); err != nil {
		fmt.Printf("Bro unable to create memory limit for conatiner %s err: %v", c.Name, err)
		return err
	}
	return nil;
}

// 1-100
func (c *Controller) SetCpuLimitV1(percentage int) error {
	cpuPath := filepath.Join("/sys/fs/cgroup/cpu", c.Name)
	if err := os.MkdirAll(cpuPath, 0644); err != nil {
		fmt.Printf("Bro unable to create cpu folder for conatiner %s", c.Name)
		return err
	}

	cpuPeriod := 100000
	cpuPeriodPath := filepath.Join(cpuPath, "cpu.cfs_period_us")
	if err := os.WriteFile(cpuPeriodPath, []byte(strconv.Itoa(cpuPeriod)), 0644); err != nil {
		fmt.Printf("Bro unable to set cpu period for conatiner %s", c.Name)
		return err
	}

	periodLimit := percentage * cpuPeriod / 100
	cpuPercentagePath := filepath.Join(cpuPath, "cpu.cfs_quota_us")
	if err := os.WriteFile(cpuPercentagePath, []byte(strconv.Itoa(periodLimit)), 0644); err != nil {
		fmt.Printf("Bro unable to set cpu limit for conatiner %s", c.Name)
		return err
	}

	return nil
}

func (c *Controller) SetCpuLimitV2(percentage int) error {
	if (percentage > 100 && percentage <= 0) {
		fmt.Printf("fuck off man no core free")
		return fmt.Errorf("bruh got rizz to ask for more than 100")
	}
	period := 100000
	quota := percentage * period / 100
	cpuPath := filepath.Join("/sys/fs/cgroup", c.Name, "cpu.max")
	if err := os.WriteFile(cpuPath, []byte(strconv.Itoa(quota)), 0644); err != nil {
		fmt.Printf("Bro unable to set cpu period for conatiner %s", c.Name)
		return err
	}

	return nil
}

func (c *Controller) CleanUp() {
	cgroups := []string{"memory", "cpu"}

	for _, cg := range cgroups {
		groupPath := filepath.Join("/sys/fs/cgroup", cg, c.Name)
		if err := os.RemoveAll(groupPath); err != nil {
			fmt.Printf("unable to cleanup %s", groupPath)
		}
	}
}

func (c *Controller) ResourceLookup() {
	resourcePath := filepath.Join("/sys/fs/cgroup/memory", c.Name, "memory.usage_in_bytes")
	buff, err := os.ReadFile(resourcePath); 
	if err != nil {
		fmt.Printf("unable to read memory usage")
		return
	}
	fmt.Printf("memory usage: %s", string(buff))
}

func (c *Controller) setupV2() error {
	cgroupPath := filepath.Join("/sys/fs/cgroup", c.Name)
	if err := os.MkdirAll(cgroupPath, 0644); err != nil {
		fmt.Printf("bro unable to make conatiner cgroup folder %v", cgroupPath)
		return err
	}

	if c.mem > 0 {
		err := c.SetMemoryLimitV2(c.mem)
		if err != nil {
			return err
		}
	}
	if c.cpu > 0 && c.cpu <= 100 {
		err := c.SetCpuLimitV2(c.cpu)
		if err != nil {
			return err
		}
	}

	cgroupProcsPath := filepath.Join(cgroupPath, "cgroup.procs")
	return os.WriteFile(cgroupProcsPath, []byte(strconv.Itoa(c.Pid)), 0644)
}

func (c *Controller) setupV1() error {
	cgroupPath := filepath.Join("/sys/fs/cgroup", c.Name)
	if c.mem > 0 {
		err := c.SetMemoryLimitV1(c.mem)
		if err != nil {
			return err
		}
	}
	if c.cpu > 0 && c.cpu <= 100 {
		err := c.SetCpuLimitV1(c.cpu)
		if err != nil {
			return err
		}
	}
	cgroupProcsPath := filepath.Join(cgroupPath, "cgroup.procs")
	return os.WriteFile(cgroupProcsPath, []byte(strconv.Itoa(c.Pid)), 0644)
}

func (c *Controller) SetupCgroups() {
	switch c.isCgroupV2() {
	case true:
		must(c.setupV2())
	case false:
		must(c.setupV1())
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
