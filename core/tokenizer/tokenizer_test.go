package tokenizer

import (
	"testing"
)

func TestSimpleTokenizer_CountTokens(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{
			name: "empty string",
			text: "",
			want: 0,
		},
		{
			name: "whitespace only",
			text: "   \t\n",
			want: 0,
		},
		{
			name: "single word",
			text: "hello",
			want: 1,
		},
		{
			name: "multiple words",
			text: "hello world",
			want: 2,
		},
		{
			name: "punctuation splitting",
			text: "hello, world!",
			want: 4, // "hello", ",", "world", "!"
		},
		{
			name: "mixed punctuation",
			text: "test (bracket) [square]",
			want: 7, // "test", "(", "bracket", ")", "[", "square", "]"  -- logic might group punctuation together or split, let's verify regex
			// regex `[\w]+|[^\s\w]+`
			// "test" -> 1
			// "(bracket)" -> "(", "bracket", ")" -> 3
			// "[square]" -> "[", "square", "]" -> 3
			// Total: 1 + 3 + 3 = 7
		},
		{
			name: "symbols",
			text: "a=b+c",
			want: 5, // "a", "=", "b", "+", "c"
		},
	}

	tokenizer := NewSimpleTokenizer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenizer.CountTokens(tt.text); got != tt.want {
				t.Errorf("SimpleTokenizer.CountTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
