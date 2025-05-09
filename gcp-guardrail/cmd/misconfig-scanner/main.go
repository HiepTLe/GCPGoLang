package main

import (
	"fmt"
	"os"

	"github.com/hieptle/gcp-guardrail/pkg/cmd/misconfig-scanner"
)

func main() {
	cmd := misconfig_scanner.GetCommand()
	
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 