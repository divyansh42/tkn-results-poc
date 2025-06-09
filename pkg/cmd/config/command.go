package config

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "set or reset results config",
	}

	configCmd.AddCommand(
		SetCommand(),
	)

	return configCmd
}
