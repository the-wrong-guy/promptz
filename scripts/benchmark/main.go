package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"

	"github.com/the-wrong-guy/promptz/core/engine"
	"github.com/the-wrong-guy/promptz/core/types"
)

type BenchmarkSample struct {
	Name     string          `json:"name"`
	Messages []types.Message `json:"messages"`
}

func main() {
	// 1. Load samples
	data, err := ioutil.ReadFile("scripts/benchmark/samples.json")
	if err != nil {
		log.Fatalf("failed to read samples file: %v", err)
	}

	var samples []BenchmarkSample
	if err := json.Unmarshal(data, &samples); err != nil {
		log.Fatalf("failed to unmarshal samples: %v", err)
	}

	// 2. Prepare output table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "SAMPLE NAME\tMODE\tBEFORE\tAFTER\tREDUCTION\t")
	fmt.Fprintln(w, "-----------\t----\t------\t-----\t---------\t")

	// 3. Run benchmarks
	modes := []types.Mode{types.ModeBalanced, types.ModeAggressive}

	for _, sample := range samples {
		for _, mode := range modes {
			req := types.OptimizeRequest{
				Messages: sample.Messages,
				Mode:     mode,
			}

			resp := engine.Optimize(req)

			reduction := resp.SavingsRatio * 100
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%.1f%%\t\n",
				sample.Name,
				mode,
				resp.TokensBefore,
				resp.TokensAfter,
				reduction,
			)
		}
		fmt.Fprintln(w, "\t\t\t\t\t") // spacer
	}

	w.Flush()
}
