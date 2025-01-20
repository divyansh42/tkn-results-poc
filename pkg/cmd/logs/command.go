package logs

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	logsCmd := &cobra.Command{
		Use:   "logs [type] [name]",
		Short: "Get logs for Tekton resources",
	}

	logsCmd.AddCommand(
		TrCommand(),
	)

	return logsCmd
}
