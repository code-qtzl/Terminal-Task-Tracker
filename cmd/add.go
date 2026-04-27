package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <description>",
	Short: "Add a new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()
		t := s.Add(args[0])
		fmt.Fprintf(cmd.OutOrStdout(), "Added task %d\n", t.ID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
