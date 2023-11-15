package main

import (
	"os"

	"github.com/rarimo/kyc-service-legacy/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
