package get

import "github.com/spf13/cobra"

func NewGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get Tekton resources",
	}

	getCmd.AddCommand(
		GetPRCommand(),
		GetTRCommand(),
	)

	return getCmd
}
