package similarity

import (
	"testing"

	"github.com/the-wrong-guy/promptz/core/types"
)

func TestJaccardSimilarity(t *testing.T) {
	tests := []struct {
		name   string
		a, b   string
		wantGT float64 // want score > this
		wantLT float64 // want score < this
	}{
		{
			name:   "identical",
			a:      "fix the bug",
			b:      "fix the bug",
			wantGT: 0.99,
			wantLT: 1.01,
		},
		{
			name:   "near duplicate",
			a:      "fix the bug",
			b:      "please fix that bug",
			wantGT: 0.3,
			wantLT: 0.8,
		},
		{
			name:   "completely different",
			a:      "hello world",
			b:      "database connection error",
			wantGT: -0.1,
			wantLT: 0.1,
		},
		{
			name:   "both empty",
			a:      "",
			b:      "",
			wantGT: 0.99,
			wantLT: 1.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JaccardSimilarity(tt.a, tt.b)
			if got <= tt.wantGT || got >= tt.wantLT {
				t.Errorf("JaccardSimilarity(%q, %q) = %.4f, want between %.2f and %.2f",
					tt.a, tt.b, got, tt.wantGT, tt.wantLT)
			}
		})
	}
}

func TestDeduplicateBySimilarity(t *testing.T) {
	msgs := []types.Message{
		{Role: "user", Content: "fix the bug in the database"},
		{Role: "user", Content: "please fix the bug in the database"},
		{Role: "assistant", Content: "I'll look into the database bug"},
		{Role: "user", Content: "something completely different"},
	}

	opts := Options{Threshold: 0.7}
	got := DeduplicateBySimilarity(msgs, opts)

	// First two user messages are near-duplicates → keep shorter
	if len(got) != 3 {
		t.Errorf("expected 3 messages after dedup, got %d", len(got))
		for i, m := range got {
			t.Logf("  [%d] %s: %s", i, m.Role, m.Content)
		}
	}

	// The kept message should be the shorter one
	if len(got) > 0 && got[0].Content != "fix the bug in the database" {
		t.Errorf("expected shorter message to be kept, got %q", got[0].Content)
	}
}
