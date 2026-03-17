package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/t2sc0m/aix/adapter"
	"github.com/t2sc0m/aix/runner"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check adapter installation and auth status",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	adapters := []adapter.Adapter{
		adapter.NewCodexAdapter(runner.NewExecRunner()),
	}

	for _, a := range adapters {
		installed := "not installed"
		if a.IsInstalled() {
			installed = "installed"
		}

		auth := a.AuthStatus()
		authStr := auth.Detail
		if auth.Authenticated {
			authStr = "authenticated (" + auth.Detail + ")"
		}

		fmt.Printf("%-10s  %-15s  %s\n", a.Name(), installed, authStr)
	}

	return nil
}
