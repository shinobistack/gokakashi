package cmd

import (
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a container image",
	Run:   scanImage,
}

func scanImage(cmd *cobra.Command, args []string) {
}
