package manager

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/a-ZINC/conti/container/container"
	"github.com/a-ZINC/conti/vm"
	"github.com/google/uuid"
)

type ContainerManager struct {
	Shell *vm.Shell

	containers map[uuid.UUID]*container.Container
	mu         sync.RWMutex
}

func NewContainerManager(shell *vm.Shell) *ContainerManager {
	return &ContainerManager{
		containers: make(map[uuid.UUID]*container.Container),
		Shell:      shell,
	}
}

func (cm *ContainerManager) AddContainer(c *container.Container) {
	defer cm.mu.Unlock()
	cm.mu.Lock()
	cm.containers[c.Id] = c
}

func (cm *ContainerManager) RemoveContainer(id string) {
	uuidId := uuid.MustParse(id)
	defer cm.mu.Unlock()
	cm.mu.Lock()
	delete(cm.containers, uuidId)
}

func (cm *ContainerManager) GetContainer(id string) (*container.Container, bool) {
	uuidId := uuid.MustParse(id)
	defer cm.mu.RUnlock()
	cm.mu.RLock()
	c, exists := cm.containers[uuidId]
	return c, exists
}

func (cm *ContainerManager) ListContainers() []*container.Container {
	defer cm.mu.RUnlock()
	cm.mu.RLock()
	containers := make([]*container.Container, 0, len(cm.containers))
	for _, c := range cm.containers {
		containers = append(containers, c)
	}
	return containers
}

func (cm *ContainerManager) Run(id string) error {
	uuidId := uuid.MustParse(id)
	defer cm.mu.Unlock()
	cm.mu.Lock()
	c, ok := cm.containers[uuidId]
	if !ok {
		return fmt.Errorf("container not found")
	}

	_, err := cm.Shell.ExecuteCommandInVM(c.Image)
	if err != nil {
		return err
	}

	return nil
}

func (cm *ContainerManager) Stop(id string) error {
	uuidId := uuid.MustParse(id)
	defer cm.mu.Unlock()
	cm.mu.Lock()
	_, ok := cm.containers[uuidId]
	if !ok {
		return fmt.Errorf("container not found")
	}
	return nil
}

func (cm *ContainerManager) CreateContainer(name, image, command string) *container.Container {
	cont := container.NewContainer(name, image)
	cm.AddContainer(cont)
	cmd := exec.Command("limactl", "shell", "conti", "--", "/usr/local/runtime", "run", command)
	cmd.Env = append(os.Environ(),
		"CONTI_CONTAINER_ID="+cont.Id.String(),
		"CONTI_CONTAINER_NAME="+cont.Name,
		"CONTI_CONTAINER_IMAGE="+cont.Image,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return nil
	}
	return cont
}
