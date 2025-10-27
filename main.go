package main

import (
	"os"

	"github.com/beyachad-maan/auth-service/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
