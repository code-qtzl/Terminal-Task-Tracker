package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid id %q: %w", args[0], err)
		}
		s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()
		if err := s.Delete(id); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Deleted task %d\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
