package normalize

import (
	"testing"

	"github.com/the-wrong-guy/promptz/core/types"
)

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "basic trim",
			input: "  hello world  ",
			want:  "hello world",
		},
		{
			name:  "collapse spaces",
			input: "hello   world",
			want:  "hello world",
		},
		{
			name:  "remove fillers",
			input: "can you please help me fix this code",
			want:  "fix this code",
		},
		{
			name:  "case insensitive fillers",
			input: "PLEASE fix this",
			want:  "fix this",
		},
		{
			name:  "keep partial words",
			input: "unpleasant",
			want:  "unpleasant", // should not remove 'please' from 'unpleasant'
		},
	}

	opts := DefaultOptions()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeText(tt.input, opts); got != tt.want {
				t.Errorf("NormalizeText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDeduplicateMessages(t *testing.T) {
	msgs := []types.Message{
		{Role: "user", Content: "hi"},
		{Role: "user", Content: "hi"},
		{Role: "assistant", Content: "hello"},
		{Role: "user", Content: "hi"},
	}

	got := DeduplicateMessages(msgs)
	if len(got) != 3 {
		t.Errorf("DeduplicateMessages() length = %d, want 3", len(got))
	}

	if got[0].Content != "hi" || got[1].Content != "hello" || got[2].Content != "hi" {
		t.Errorf("DeduplicateMessages() content mismatch")
	}
}
