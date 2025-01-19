package main

import (
	"log"
	"os"

	"github.com/divyansh42/tkn-results/pkg/cmd"
)

func main() {
	rootCmd := pkg.Root()

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(1)
	}
}
