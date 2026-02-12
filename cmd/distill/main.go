package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/the-wrong-guy/promptz/core/engine"
	"github.com/the-wrong-guy/promptz/core/types"
)

func main() {
	// Read from Stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	// Parse Request
	var req types.OptimizeRequest
	if err := json.Unmarshal(input, &req); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Optimize
	resp := engine.Optimize(req)

	// Output Response
	output, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
