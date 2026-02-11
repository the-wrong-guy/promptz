package rewrite

import (
	"strings"
	"testing"

	"github.com/the-wrong-guy/promptz/core/types"
)

func TestRewriteMessage(t *testing.T) {
	tests := []struct {
		name         string
		mode         types.Mode
		input        string
		wantContains []string
		wantMissing  []string
	}{
		{
			name:         "conservative no change",
			mode:         types.ModeConservative,
			input:        "the quick brown fox",
			wantContains: []string{"the", "quick", "brown", "fox"},
		},
		{
			name:         "balanced removes stop words",
			mode:         types.ModeBalanced,
			input:        "the quick brown fox is fast",
			wantContains: []string{"quick", "brown", "fox", "fast"},
			wantMissing:  []string{"the", "is"},
		},
		{
			name:         "aggressive keywords context",
			mode:         types.ModeAggressive,
			input:        "Hello my world, whats up",
			wantContains: []string{"Hello", "world", "whats", "up"},
			wantMissing:  []string{"my"}, // "my" should likely be removed as a stop word/possessive not critical
		},
		{
			name:         "aggressive keywords technical",
			mode:         types.ModeAggressive,
			input:        "fix the api bug in prod",
			wantContains: []string{"fix", "api", "bug", "prod"}, // "in", "the" gone. "prod" > 3 chars. "api", "fix", "bug" in important list.
			wantMissing:  []string{"the"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := types.Message{Content: tt.input}
			got := RewriteMessage(msg, tt.mode)

			for _, c := range tt.wantContains {
				if !strings.Contains(got, c) {
					t.Errorf("RewriteMessage(%s) = %q, expected to contain %q", tt.mode, got, c)
				}
			}

			for _, m := range tt.wantMissing {
				// Naive check: ensure word is not present as a distinct token
				// e.g. "my" shouldn't be found in "myCode"
				// But strings.Contains("myCode", "my") is true.
				// Better check: fields
				fields := strings.Fields(got)
				for _, f := range fields {
					clean := strings.TrimRight(f, ".,!?:;\"'()[]{}")
					if clean == m {
						t.Errorf("RewriteMessage(%s) = %q, should NOT contain %q", tt.mode, got, m)
					}
				}
			}
		})
	}
}
