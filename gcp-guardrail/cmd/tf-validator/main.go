package main

import (
	"fmt"
	"os"

	"github.com/hieptle/gcp-guardrail/pkg/cmd/tf-validator"
)

func main() {
	cmd := tf_validator.GetCommand()
	
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 