package get

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get Tekton resources",
	}

	getCmd.AddCommand(
		PrCommand(),
		TrCommand(),
	)

	return getCmd
}
