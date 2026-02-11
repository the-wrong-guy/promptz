package tokenizer

import (
	"regexp"
	"strings"
)

// Tokenizer defines the interface for counting tokens in text.
type Tokenizer interface {
	CountTokens(text string) int
}

// SimpleTokenizer provides a deterministic approximation of token count
// by splitting on whitespace and punctuation.
type SimpleTokenizer struct{}

// NewSimpleTokenizer creates a new instance of SimpleTokenizer.
func NewSimpleTokenizer() *SimpleTokenizer {
	return &SimpleTokenizer{}
}

// CountTokens returns the approximate number of tokens in the text.
// This implementation splits by whitespace and common punctuation roughly
// mimicking how BPE tokenizers might see word boundaries, but much simpler.
func (t *SimpleTokenizer) CountTokens(text string) int {
	if strings.TrimSpace(text) == "" {
		return 0
	}

	// Split by whitespace first
	fields := strings.Fields(text)
	count := 0

	// Regex to match punctuation that often form separate tokens
	// We'll treat a sequence of alphanumeric characters as a token,
	// and sequences of punctuation as separate tokens, roughly.
	// For a simple MVP:
	// 1. Split by whitespace.
	// 2. Further split punctuation from words if they are attached.
	//    e.g. "Hello," -> "Hello" ","

	// A simple meaningful approximation:
	// specific punctuation Often gets its own token: . , ! ? : ; " ' ( ) [ ] { }

	// We basically want to count groups of word-chars and groups of non-word-chars.
	re := regexp.MustCompile(`[\w]+|[^\s\w]+`)

	for _, field := range fields {
		matches := re.FindAllString(field, -1)
		count += len(matches)
	}

	return count
}
