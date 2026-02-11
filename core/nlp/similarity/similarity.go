package similarity

import (
	"strings"

	"github.com/the-wrong-guy/promptz/core/types"
)

// DefaultThreshold is the default Jaccard similarity threshold for near-duplicate detection.
const DefaultThreshold = 0.7

// Options controls similarity-based deduplication behavior.
type Options struct {
	// Threshold is the Jaccard similarity threshold (0.0 to 1.0).
	// Messages with similarity >= threshold are considered near-duplicates.
	Threshold float64
}

// DefaultOptions returns default similarity options.
func DefaultOptions() Options {
	return Options{
		Threshold: DefaultThreshold,
	}
}

// JaccardSimilarity computes the Jaccard similarity coefficient between two texts.
// Returns a value between 0.0 (completely different) and 1.0 (identical word sets).
func JaccardSimilarity(a, b string) float64 {
	setA := wordSet(a)
	setB := wordSet(b)

	if len(setA) == 0 && len(setB) == 0 {
		return 1.0 // both empty = identical
	}

	intersection := 0
	for word := range setA {
		if setB[word] {
			intersection++
		}
	}

	union := len(setA)
	for word := range setB {
		if !setA[word] {
			union++
		}
	}

	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

// DeduplicateBySimlarity removes consecutive near-duplicate messages from same role.
// When near-duplicates are found, the shorter message is kept (more concise).
func DeduplicateBySimilarity(messages []types.Message, opts Options) []types.Message {
	if len(messages) == 0 {
		return messages
	}

	result := make([]types.Message, 0, len(messages))
	result = append(result, messages[0])

	for i := 1; i < len(messages); i++ {
		curr := messages[i]
		prev := result[len(result)-1]

		// Only deduplicate consecutive messages from the same role
		if curr.Role == prev.Role {
			sim := JaccardSimilarity(curr.Content, prev.Content)
			if sim >= opts.Threshold {
				// Keep the shorter one (replace previous if current is shorter)
				if len(curr.Content) < len(prev.Content) {
					result[len(result)-1] = curr
				}
				continue
			}
		}

		result = append(result, curr)
	}

	return result
}

func wordSet(text string) map[string]bool {
	set := make(map[string]bool)
	for _, word := range strings.Fields(strings.ToLower(text)) {
		clean := strings.TrimRight(word, ".,!?:;\"'()[]{}")
		if clean != "" {
			set[clean] = true
		}
	}
	return set
}
