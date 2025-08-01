package cmd

import (
	"fmt"
	"os"

	"github.com/shinobistack/gokakashi/internal/experiment"
	"github.com/spf13/cobra"
)

var (
	server       string
	token        string
	workspace    string
	name         string
	id           int
	chidori      bool
	labels       string
	singleStrike bool
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
	rootCmd.PersistentFlags().StringVar(&server, "server", "", "The server address to connect to")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Authentication token for the server")
	rootCmd.PersistentFlags().StringVar(&experiment.Experiments, "experiments", "", "Comma-separated list of experiment flags (e.g., --experiments=exp1,exp2)")

	// To independently process its own configuration file path.
	serverConfigFilePath = serverCmd.Flags().String("config", "", "Path to the config YAML file")
	rootCmd.AddCommand(serverCmd)

	// Flags for the `agent start` command
	agentStartCmd.Flags().StringVar(&name, "name", "", "Unique name for the agent (optional, defaults to agent-<random_suffix>)")
	agentStartCmd.Flags().StringVar(&workspace, "workspace", "", "Workspace for the agent (optional, defaults to /tmp/<agent-name>)")
	agentStartCmd.Flags().StringVar(&labels, "labels", "", "Labels for the agent in key=value format (e.g., --labels=\"key1=value1,key2=value2\")")
	agentStartCmd.Flags().BoolVar(&singleStrike, "single-strike", false, "Run as an ephemeral agent that executes one task and then exits. For ephemeral agents please assign unique labels like github job ID or UUID etc")

	// Flags for the `agent stop` command
	agentStopCmd.Flags().IntVar(&id, "id", 0, "ID of the agent to deregister")
	agentStopCmd.Flags().StringVar(&name, "name", "", "Name of the agent to deregister")
	agentStopCmd.Flags().BoolVar(&chidori, "chidori", false, "Hard delete the agent and its associated tasks")

	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentStartCmd)
	agentCmd.AddCommand(agentStopCmd)

	// Flags for the `scan image` command
	scanImageCmd.Flags().StringVar(&image, "image", "", "Container image to scan")
	scanImageCmd.Flags().StringVar(&policyName, "policy", "", "Policy name for the scan")
	scanImageCmd.Flags().StringVar(&labels, "labels", "", "Labels for the scans in key=value format (e.g., --labels=\"key1=value1,key2=value2\"). Can be used to map agent's task the particular scan via labels")
	scanImageCmd.Flags().StringVar(&scanTimeout, "timeout", "", "Timeout for the scan (e.g., 10m, 30m, 1h)")

	// Flags for the `scan status` command
	scanStatusCmd.Flags().StringVar(&scanID, "scanID", "", "Scan ID to check status")

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
