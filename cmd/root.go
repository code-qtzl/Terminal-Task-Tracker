package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/code-qtzl/tasks/internal/store"
	"github.com/spf13/cobra"
)

var dataFile string

var rootCmd = &cobra.Command{
	Use:           "tasks",
	Short:         "Manage your todo tasks from the terminal",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dataFile, "file", "f", "", "path to tasks CSV file (defaults to $TASKS_FILE or ~/.tasks.csv)")
}

func resolveFile() (string, error) {
	if dataFile != "" {
		return dataFile, nil
	}
	if env := os.Getenv("TASKS_FILE"); env != "" {
		return env, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".tasks.csv"), nil
}

func openStore() (*store.Store, error) {
	path, err := resolveFile()
	if err != nil {
		return nil, err
	}
	return store.Open(path)
}
