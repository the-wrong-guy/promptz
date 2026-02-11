package engine

import (
	"testing"

	"github.com/the-wrong-guy/promptz/core/types"
)

func TestOptimize_BasicDedup(t *testing.T) {
	req := types.OptimizeRequest{
		Mode: types.ModeBalanced,
		Messages: []types.Message{
			{Role: "user", Content: "   Hello world   "},
			{Role: "user", Content: "Hello world"},
			{Role: "assistant", Content: "Hi there, how can I help you today?"},
		},
	}

	resp := Optimize(req)

	if len(resp.Optimized) != 2 {
		t.Errorf("Expected 2 messages after optimization, got %d", len(resp.Optimized))
		for i, m := range resp.Optimized {
			t.Logf("  [%d] %s: %q", i, m.Role, m.Content)
		}
	}

	if resp.TokensAfter >= resp.TokensBefore {
		t.Errorf("Expected token reduction. Before: %d, After: %d", resp.TokensBefore, resp.TokensAfter)
	}

	if resp.SavingsRatio <= 0 {
		t.Errorf("Expected positive savings ratio, got %f", resp.SavingsRatio)
	}
}

func TestOptimize_NearDuplicateDedup(t *testing.T) {
	req := types.OptimizeRequest{
		Mode: types.ModeBalanced,
		Messages: []types.Message{
			{Role: "user", Content: "fix the bug in the database"},
			{Role: "user", Content: "please fix the bug in the database now"},
			{Role: "assistant", Content: "I'll investigate the issue"},
		},
	}

	resp := Optimize(req)

	// The two user messages are near-duplicates and should be merged
	userMsgs := 0
	for _, m := range resp.Optimized {
		if m.Role == "user" {
			userMsgs++
		}
	}

	if userMsgs != 1 {
		t.Errorf("Expected 1 user message after near-duplicate dedup, got %d", userMsgs)
		for i, m := range resp.Optimized {
			t.Logf("  [%d] %s: %q", i, m.Role, m.Content)
		}
	}
}

func TestOptimize_Compression(t *testing.T) {
	req := types.OptimizeRequest{
		Mode: types.ModeAggressive,
		Messages: []types.Message{
			{Role: "user", Content: "I think that the server (which is the main one) is down due to the fact that the database connection failed"},
		},
	}

	resp := Optimize(req)

	if resp.SavingsRatio <= 0.2 {
		t.Errorf("Expected significant savings ratio for verbose input, got %f", resp.SavingsRatio)
	}

	t.Logf("Before: %d tokens, After: %d tokens, Savings: %.1f%%",
		resp.TokensBefore, resp.TokensAfter, resp.SavingsRatio*100)
	if len(resp.Optimized) > 0 {
		t.Logf("Optimized: %q", resp.Optimized[0].Content)
	}
}

func TestOptimize_ConservativeMinimalChange(t *testing.T) {
	req := types.OptimizeRequest{
		Mode: types.ModeConservative,
		Messages: []types.Message{
			{Role: "user", Content: "please help me fix the database error"},
		},
	}

	resp := Optimize(req)

	// Conservative should make minimal changes (mostly normalization)
	if len(resp.Optimized) != 1 {
		t.Errorf("Expected 1 message, got %d", len(resp.Optimized))
	}
}
