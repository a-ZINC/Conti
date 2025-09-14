# Conti

A Go-based container and virtual machine management platform that provides utilities for running, managing, and orchestrating containerized applications and virtual machines with enhanced isolation and control.

## Features

### 🐳 Container Management
- **Create and Deploy:** Spin up containers with custom configurations and isolated environments
- **Lifecycle Management:** Start, stop, restart, and monitor container states
- **Custom Namespaces:** Advanced namespace isolation for enhanced security and resource separation
- **Resource Control:** Manage CPU, memory, and storage allocations for containers

### 🖥️ Virtual Machine Support
- **VM Provisioning:** Launch and control virtual machines with flexible configurations
- **Shell Access:** Interactive shell access to VMs for debugging and administration
- **Resource Management:** Dynamic allocation and monitoring of VM resources
- **State Management:** Save, restore, and migrate VM states

### 💻 Client Interface
- **Unified CLI:** Single command-line interface for managing both containers and VMs
- **Interactive Commands:** Streamlined operations with intuitive command structure
- **Real-time Monitoring:** Live status updates and resource monitoring
- **Batch Operations:** Execute multiple operations efficiently

### 🗂️ Filesystem Management
- **Root Filesystem Control:** Manage and customize root filesystems for containers and VMs
- **Mount Point Management:** Dynamic mounting and unmounting of storage volumes
- **Overlay Networks:** Support for complex filesystem overlays and unions
- **Storage Backends:** Multiple storage driver support for flexibility

### 🔒 Security & Isolation
- **Namespace Runner:** Execute processes in completely isolated namespaces
- **Security Policies:** Configurable security contexts and access controls
- **Network Isolation:** Separate network namespaces for enhanced security
- **User Namespace Support:** Run containers with unprivileged user mappings

### 🎨 Developer Experience
- **Colorful CLI Output:** Enhanced terminal output with color-coded messages for better readability
- **Extensible Architecture:** Modular design allowing easy addition of plugins and extensions
- **Rich Logging:** Comprehensive logging with multiple verbosity levels
- **Error Handling:** Detailed error messages and recovery suggestions

## Quick Start

### Prerequisites
- Go 1.19 or higher
- Linux operating system (for namespace support)
- Root privileges (for container/VM management)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/conti.git
cd conti

# Build the project
make build

# Run the application
make run
```

### Basic Usage

```bash
# Start a container
./conti container create --image ubuntu:latest --name my-container

# Launch a VM
./conti vm create --memory 2G --cpu 2 --name my-vm

# List running containers and VMs
./conti list

# Access container shell
./conti container shell my-container

# Access VM shell
./conti vm shell my-vm
```

## Project Structure

```
conti/
├── host/                    # Main application code
│   ├── main.go             # Host application entry point
│   ├── client/             # Client interface logic
│   ├── container/          # Container management modules
│   ├── utils/              # Utility functions and helpers
│   ├── vm/                 # VM management and shell utilities
│   ├── conti               # Compiled binary
│   ├── Makefile           # Build and run commands
│   └── go.mod             # Module dependencies
├── runtime/                # Runtime execution environment
│   ├── main.go            # Runtime entry point
│   ├── pkg/filesystem/    # Filesystem management
│   ├── pkg/runner/        # Namespace and process runner
│   └── go.mod            # Runtime dependencies
└── README.md             # This file
```

## Building from Source

```bash
# Development build
make build

# Production build with optimizations
make build-prod

# Clean build artifacts
make clean

# Run tests
make test

# Install dependencies
go mod tidy
```





---

**Note:** Conti is under active development. APIs and features may change between versions. Please check the [CHANGELOG](CHANGELOG.md) for breaking changes and updates.
