package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/a-ZINC/conti/container/manager"
	"github.com/a-ZINC/conti/utils"
)

type Client struct {
	ContainerManager *manager.ContainerManager
}

func NewClient(manager *manager.ContainerManager) *Client {
	return &Client{
		ContainerManager: manager,
	}
}

func (c *Client) Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("conti> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		if strings.ToLower(line) == "quit" || strings.ToLower(line) == "exit" {
			return
		}

		c.handleCommand(line)
	}
}

func (c *Client) handleCommand(input string) {
	fmt.Println("Input command:", input)
	lineParts := utils.Tokenize(input)
	fmt.Printf("Tokenized parts: %v\n", lineParts)
	if len(lineParts) < 1  {
		fmt.Println("Usage: container run --name <name> --image <image> --command <command>")
		return
	}
	switch lineParts[0] {
		case "container":
			switch lineParts[1] {
			case "run":
				c.handleRunContainer(lineParts[2:])
			case "stop":
			case "list":
			default:
				fmt.Println("Unknown command:", lineParts[1])
			}

	}
}

func (c *Client) handleRunContainer(lineParts []string) {
	var name, image string
	var command []string
	for i := 0; i < len(lineParts); i++ {
		switch lineParts[i] {
		case "--name", "-n":
			if i+1 < len(lineParts) {
				name = lineParts[i+1]
				i++
			}
		case "--image", "-i":
			if i+1 < len(lineParts) {
				image = lineParts[i+1]
				i++
			}
		case "--command", "-c":
			if i+1 < len(lineParts) {
				command = append(command, lineParts[i+1])
				i++
			}
		}
	}

	if name == "" || image == "" || len(command) == 0 {
		fmt.Printf("Missing required flags:\n")
		if name == "" {
			fmt.Println("  --name <name>")
		}
		if image == "" {
			fmt.Println("  --image <image>")
		}
		if len(command) == 0 {
			fmt.Println("  --command <command>")
		}
		return
	}
	commandStr := strings.Join(command, " ")
	fmt.Printf("Running container with Name: %s, Image: %s, Command: %s\n", name, image, commandStr)

	// Call the container manager to run the container
	c.ContainerManager.CreateContainer(name, image, commandStr)
}
