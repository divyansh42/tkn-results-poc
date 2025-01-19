package main

import (
	"fmt"
	"os"

	"github.com/divyansh42/tkn-results/pkg/cmd"
)

func main() {
	rootCmd := pkg.Root()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
