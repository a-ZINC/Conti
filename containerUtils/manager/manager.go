package manager

import (
	"fmt"
	"sync"

	"github.com/a-ZINC/conti/containerUtils/container"
	"github.com/a-ZINC/conti/vm"
	"github.com/google/uuid"
)

type ContainerManager struct {
	shell *vm.Shell

	containers map[uuid.UUID]*container.Container
	mu sync.RWMutex
}

func NewContainerManager(shell *vm.Shell) *ContainerManager {
	return &ContainerManager{
		containers: make(map[uuid.UUID]*container.Container),
		shell:     shell,
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
	if (!ok) {
		return fmt.Errorf("container not found")
	}

	output, err := cm.shell.ExecuteCommand(c.Command)
	if err != nil {
		return err
	}
	fmt.Printf("output: %s", output)

	return nil
}

func (cm *ContainerManager) Stop(id string) error {
	uuidId := uuid.MustParse(id)
	defer cm.mu.Unlock()
	cm.mu.Lock()
	c, ok := cm.containers[uuidId]
	if (!ok) {
		return fmt.Errorf("container not found")
	}
	killCmd := fmt.Sprintf("kill %d", c.Pid)
	output, err := cm.shell.ExecuteCommand(killCmd)
	if err != nil {
		return err
	}
	fmt.Printf("output: %s", output)
	
	return nil
} 