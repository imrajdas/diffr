package main

import (
	"fmt"
	"os"

	"github.com/imrajdas/diffr/pkg/cmd/root"
)

var Version string

func main() {
	err := os.Setenv("Version", Version)
	if err != nil {
		fmt.Println("Failed to fetched CLIVersion")
	}

	root.Execute()
}
