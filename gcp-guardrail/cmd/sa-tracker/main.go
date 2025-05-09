package main

import (
	"fmt"
	"os"

	"github.com/hieptle/gcp-guardrail/pkg/cmd/sa-tracker"
)

func main() {
	cmd := sa_tracker.GetCommand()
	
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 