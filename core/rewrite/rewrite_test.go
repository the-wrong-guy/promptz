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
		messages     []types.Message
		targetIdx    int
		wantContains []string
		wantMissing  []string
	}{
		{
			name:         "conservative no change",
			mode:         types.ModeConservative,
			messages:     []types.Message{{Role: "user", Content: "the quick brown fox"}},
			targetIdx:    0,
			wantContains: []string{"the", "quick", "brown", "fox"},
		},
		{
			name:         "balanced removes function words",
			mode:         types.ModeBalanced,
			messages:     []types.Message{{Role: "user", Content: "the quick brown fox is fast"}},
			targetIdx:    0,
			wantContains: []string{"quick", "brown", "fox", "fast"},
			wantMissing:  []string{"the"},
		},
		{
			name: "aggressive keeps content words",
			mode: types.ModeAggressive,
			messages: []types.Message{
				{Role: "user", Content: "fix the database connection error in production"},
			},
			targetIdx:    0,
			wantContains: []string{"database", "connection", "error", "production"},
			wantMissing:  []string{"the", "in"},
		},
		{
			name: "aggressive preserves negation",
			mode: types.ModeAggressive,
			messages: []types.Message{
				{Role: "user", Content: "do not delete the important files"},
			},
			targetIdx:    0,
			wantContains: []string{"not", "delete", "important", "files"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRewriter(tt.messages)
			got := r.RewriteMessage(tt.messages[tt.targetIdx], tt.mode)

			for _, c := range tt.wantContains {
				if !strings.Contains(strings.ToLower(got), strings.ToLower(c)) {
					t.Errorf("RewriteMessage(%s) = %q, expected to contain %q", tt.mode, got, c)
				}
			}

			for _, m := range tt.wantMissing {
				fields := strings.Fields(got)
				for _, f := range fields {
					if strings.EqualFold(f, m) {
						t.Errorf("RewriteMessage(%s) = %q, should NOT contain %q", tt.mode, got, m)
					}
				}
			}
		})
	}
}
