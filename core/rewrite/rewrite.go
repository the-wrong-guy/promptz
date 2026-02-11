package rewrite

import (
	"regexp"
	"strings"

	"github.com/the-wrong-guy/promptz/core/types"
)

// RewriteMessage transforms message content based on the optimization mode.
func RewriteMessage(msg types.Message, mode types.Mode) string {
	content := msg.Content

	switch mode {
	case types.ModeConservative:
		// Conservative: minimal changes, mainly relied on normalization which is done before rewrite
		return content

	case types.ModeBalanced:
		// Balanced: remove only very obvious non-keywords, stop words
		return removeStopWords(content)

	case types.ModeAggressive:
		// Aggressive: keep only essential keywords (simple approximation)
		return extractKeywords(content)

	default:
		return content
	}
}

var stopWords = map[string]bool{
	"the": true, "a": true, "an": true,
	"in": true, "on": true, "at": true, "to": true, "for": true, "of": true, "with": true,
	"is": true, "are": true, "was": true, "were": true, "be": true,
	"this": true, "that": true, "these": true, "those": true,
	"it": true, "they": true, "them": true,
}

func removeStopWords(text string) string {
	return processWords(text, func(word string) bool {
		return !stopWords[strings.ToLower(word)]
	})
}

// extractKeywords is an improved heuristic keyword extractor.
func extractKeywords(text string) string {
	return processWords(text, func(word string) bool {
		// Clean punctuation for check
		cleanWord := strings.TrimRight(word, ".,!?:;\"'()[]{}")
		lowerWord := strings.ToLower(cleanWord)

		// 1. Keep if in "Important Short Words" list.
		if importantShortWords[lowerWord] {
			return true
		}

		// 2. Drop if in "Stop Words" list.
		if stopWords[lowerWord] {
			return false
		}

		// 3. Keep if Length > 3 (and not a stop word, handled above).
		if len(cleanWord) > 3 {
			return true
		}

		// 4. Keep if Capitalized (heuristic for entities).
		if len(cleanWord) > 0 && cleanWord[0] >= 'A' && cleanWord[0] <= 'Z' {
			return true
		}

		// 5. Keep if contains digits.
		if regexp.MustCompile(`\d`).MatchString(cleanWord) {
			return true
		}

		return false
	})
}

var importantShortWords = map[string]bool{
	"no": true, "not": true, "nor": true,
	"up": true, "down": true, "on": true, "off": true,
	"in": true, "out": true,
	"why": true, "how": true, "who": true, "what": true, "when": true, "where": true,
	"can": true, "may": true, "must": true,
	"fix": true, "bug": true, "err": true, "log": true, "api": true, "sql": true, "db": true,
	"run": true, "go": true, "py": true, "js": true, "ts": true,
}

func processWords(text string, keepFunc func(string) bool) string {
	// Simple split, filter, join
	// We need to preserve punctuation ideally, but for "rewrite" we often sacrifice grammar.
	// Let's try to preserve basic structure by just filtering the word tokens.

	words := strings.Fields(text)
	var kept []string

	for _, word := range words {
		// Strip punctuation for the check, but keep it in the result?
		// For MVP, simplistic approach: strip punctuation for check.

		cleanWord := strings.TrimRight(word, ".,!?:;\"'()[]{}")
		if keepFunc(cleanWord) {
			kept = append(kept, word)
		}
	}

	return strings.Join(kept, " ")
}
