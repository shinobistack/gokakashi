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
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	serverConfigFilePath = serverCmd.Flags().String("config", "", "Path to the config YAML file")
	rootCmd.AddCommand(serverCmd)

	rootCmd.AddCommand(scanCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
