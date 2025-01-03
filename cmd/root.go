package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gokakashi",
	Short: "GoKakashi - The Container image vulnerability management platform",
	Run: func(cmd *cobra.Command, args []string) {
		// Display help if no subcommand is provided
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// To independently process its own configuration file path.
	serverConfigFilePath = serverCmd.Flags().String("config", "", "Path to the config YAML file")
	rootCmd.AddCommand(serverCmd)

	rootCmd.AddCommand(scanCmd)

	agentStartCmd.Flags().StringVar(&server, "server", "", "The server address to connect to")
	agentStartCmd.Flags().StringVar(&token, "token", "", "Authentication token for the server")
	agentStartCmd.Flags().StringVar(&workspace, "workspace", "", "Path to the local workspace")
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentStartCmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
