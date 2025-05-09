package main

import (
	"fmt"
	"os"

	"github.com/hieptle/gcp-guardrail/pkg/cmd/log-watcher"
)

func main() {
	cmd := log_watcher.GetCommand()
	
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 