package main

import (
	"os"

	"github.com/marv972228/sandbox_judge/cmd/judge/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
