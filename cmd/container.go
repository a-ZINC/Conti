package cmd

import "github.com/spf13/cobra"

var RunArgs struct {
	Name    string
	Image   string
	Command string
}

var ContainerCmd = &cobra.Command{
	Use:   "container",
	Short: "Run the container command",
	Long:  `This command allows you to run the container application`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var RunContainerCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a container",
	Long:  `This command allows you to run a container with specified name, image, and command`,
	Run: func(cmd *cobra.Command, args []string) {
		RunArgs.Name = cmd.Flag("name").Value.String()
		RunArgs.Image = cmd.Flag("image").Value.String()
		RunArgs.Command = cmd.Flag("command").Value.String()

		
	},
}

func init() {
	ContainerCmd.AddCommand(RunContainerCmd)

	RunContainerCmd.Flags().StringVarP(&RunArgs.Name, "name", "n", "", "Name of the container")
	RunContainerCmd.Flags().StringVarP(&RunArgs.Image, "image", "i", "", "Image to use for the container")
	RunContainerCmd.Flags().StringVarP(&RunArgs.Command, "command", "c", "", "Command to run in the container")

	RunContainerCmd.MarkFlagRequired("name")
	RunContainerCmd.MarkFlagRequired("image")
	RunContainerCmd.MarkFlagRequired("command")
}