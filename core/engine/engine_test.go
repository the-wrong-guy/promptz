package engine

import (
	"testing"

	"github.com/the-wrong-guy/promptz/core/types"
)

func TestOptimize(t *testing.T) {
	req := types.OptimizeRequest{
		Mode: types.ModeBalanced,
		Messages: []types.Message{
			{Role: "user", Content: "   Hello world   "},
			{Role: "user", Content: "Hello world"},
			{Role: "assistant", Content: "Hi there, how can I help you today?"},
		},
	}

	resp := Optimize(req)

	// User messages should be deduplicated
	// "Hi there..." should be rewritten

	if len(resp.Optimized) != 2 {
		t.Errorf("Expected 2 messages after optimization, got %d", len(resp.Optimized))
	}

	// Verify token count reduction
	if resp.TokensAfter >= resp.TokensBefore {
		t.Errorf("Expected token reduction. Before: %d, After: %d", resp.TokensBefore, resp.TokensAfter)
	}

	if resp.SavingsRatio <= 0 {
		t.Errorf("Expected positive savings ratio, got %f", resp.SavingsRatio)
	}
}
