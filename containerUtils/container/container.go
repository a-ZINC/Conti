package container

import "github.com/google/uuid"

type Container struct {
	Name string
	Pid int
	Image string
	Command string
	Id uuid.UUID
}

func NewContainer(name, image, command string, pid int) *Container {
	id := uuid.New()
	return &Container{
		Name: name,
		Image: image,
		Command: command,
		Pid: pid,
		Id: id,
	}
}