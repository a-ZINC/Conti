package container

import (
	"os/exec"

	"github.com/a-ZINC/conti/vm"
	"github.com/google/uuid"
)

type Container struct {
	Name     string
	Pid      int
	Image    string
	Id       uuid.UUID
	Cmd      *exec.Cmd
	Provider vm.VMProvider
	Rootfs   string
	MemoryLimit int64
	CPULimit    int
}

func NewContainer(name, image string) *Container {
	id := uuid.New()
	return &Container{
		Name:    name,
		Image:   image,
		Id:      id,
	}
}

func (c *Container) Stop() error {
	return c.Cmd.Process.Kill()
}


