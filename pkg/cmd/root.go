package pkg

import (
	"fmt"
	"github.com/divyansh42/tkn-results/pkg/cmd/get"
	"github.com/divyansh42/tkn-results/pkg/cmd/logs"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	namespace string
	server    string
	token     string
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tkn",
		Short: "CLI to interact with Tekton results",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if server == "" {
				return fmt.Errorf("server URL must be provided via --server flag")
			}
			if token == "" {
				return fmt.Errorf("token must be provided")
			}
			if namespace == "" {
				ns, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
				if err == nil {
					namespace = string(ns)
				} else {
					namespace = "default"
				}
			}
			if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
				server = "https://" + server
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Namespace to use")
	rootCmd.PersistentFlags().StringVar(&server, "server", "", "API server address")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "API token")

	rootCmd.AddCommand(
		get.Command(),
		logs.Command(),
	)

	return rootCmd
}
