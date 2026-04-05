package main

import (
	"fmt"
	"os"

	"github.com/mswiente/pkm.ai/internal/cli"
	"github.com/mswiente/pkm.ai/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pkm: config error: %v\n", err)
		os.Exit(1)
	}

	root := cli.NewRootCommand(cfg)
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "pkm: %v\n", err)
		os.Exit(1)
	}
}
