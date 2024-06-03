package main

import (
	"fmt"
	"os"

	cmd "github.com/DoraTiger/safeline-utils/cmd/app/commands"
)

func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.VersionCmd,
		cmd.CertCmd,
	)
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v", err)
		os.Exit(1)
	}
}
