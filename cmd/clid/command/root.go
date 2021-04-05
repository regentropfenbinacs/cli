package command

import "github.com/spf13/cobra"

var (
	RootCmd = &cobra.Command{
		Use:   "root",
		Short: "Root Command",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return nil
		},
	}
)
