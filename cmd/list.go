package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

var listAll bool

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List tasks (uncompleted by default)",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}
		defer s.Close()

		tasks := s.List(listAll)
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 4, ' ', 0)
		if listAll {
			fmt.Fprintln(w, "ID\tTask\tCreated\tDone")
			for _, t := range tasks {
				fmt.Fprintf(w, "%d\t%s\t%s\t%t\n", t.ID, t.Description, timediff.TimeDiff(t.CreatedAt), t.IsComplete)
			}
		} else {
			fmt.Fprintln(w, "ID\tTask\tCreated")
			for _, t := range tasks {
				fmt.Fprintf(w, "%d\t%s\t%s\n", t.ID, t.Description, timediff.TimeDiff(t.CreatedAt))
			}
		}
		return w.Flush()
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "show completed tasks too")
	rootCmd.AddCommand(listCmd)
}
