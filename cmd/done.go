package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "Mark a task as done",
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
		if err := s.Complete(id); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Marked task %d as done\n", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
