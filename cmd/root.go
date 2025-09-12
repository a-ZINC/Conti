package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "conti",
	Short: "Container management CLI",
	Long:  `This is a CLI tool for managing containers`,
}

func init() {
	RootCmd.AddCommand(ContainerCmd)
}