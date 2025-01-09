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

	agentStartCmd.Flags().StringVar(&server, "server", "", "The server address to connect to")
	agentStartCmd.Flags().StringVar(&token, "token", "", "Authentication token for the server")
	agentStartCmd.Flags().StringVar(&workspace, "workspace", "", "Path to the local workspace")
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentStartCmd)

	// Flags for the `scan image` command
	scanImageCmd.Flags().StringVar(&image, "image", "", "Container image to scan")
	scanImageCmd.Flags().StringVar(&policyName, "policy", "", "Policy name for the scan")
	scanImageCmd.Flags().StringVar(&server, "server", "", "The server address to connect to")
	scanImageCmd.Flags().StringVar(&token, "token", "", "Authentication token for the server")

	// Flags for the `scan status` command
	scanStatusCmd.Flags().StringVar(&scanID, "scanID", "", "Scan ID to check status")
	scanStatusCmd.Flags().StringVar(&server, "server", "", "The server address to connect to")
	scanStatusCmd.Flags().StringVar(&token, "token", "", "Authentication token for the server")

	rootCmd.AddCommand(scanCmd)
	scanCmd.AddCommand(scanImageCmd)
	scanCmd.AddCommand(scanStatusCmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
