package compress

import "testing"

func TestCompress(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "remove parenthetical",
			input: "fix the bug (which is critical) now",
			want:  "fix the bug now",
		},
		{
			name:  "remove em-dash aside",
			input: "the server -- as you may know -- is down",
			want:  "the server is down",
		},
		{
			name:  "remove filler clause",
			input: "in order to fix this we need to restart",
			want:  "fix this we need to restart",
		},
		{
			name:  "remove verbose preamble",
			input: "it should be noted that the api is slow",
			want:  "the api is slow",
		},
		{
			name:  "no change needed",
			input: "fix the bug",
			want:  "fix the bug",
		},
		{
			name:  "multiple patterns",
			input: "i think that the server (main one) is basically what we need to fix",
			want:  "the server is we need to fix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compress(tt.input)
			if got != tt.want {
				t.Errorf("Compress(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
