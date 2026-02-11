package normalize

import (
	"regexp"
	"strings"

	"github.com/the-wrong-guy/promptz/core/types"
)

// Options configuration for normalization.
type Options struct {
	RemoveFillerPhrases bool
	AdditionalFillers   []string
}

// DefaultOptions returns default normalization options.
func DefaultOptions() Options {
	return Options{
		RemoveFillerPhrases: true,
	}
}

// NormalizeText applies a series of text cleanup operations.
func NormalizeText(text string, opts Options) string {
	text = trimWhitespace(text)
	text = collapseRepeatedSpaces(text)
	if opts.RemoveFillerPhrases {
		text = removeFillerPhrases(text, opts.AdditionalFillers)
	}
	// Re-trim in case removal left edges
	return trimWhitespace(text)
}

func trimWhitespace(text string) string {
	return strings.TrimSpace(text)
}

func collapseRepeatedSpaces(text string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(text, " ")
}

var commonFillers = []string{
	"please",
	"can you",
	"could you",
	"would you",
	"kindly",
	"i want to",
	"i need to",
	"help me",
	"just",
	"basically",
	"actually",
	"literally",
	"um",
	"uh",
	"like",
}

func removeFillerPhrases(text string, additional []string) string {
	fillers := append(commonFillers, additional...)

	// Case insensitive replacement
	// We need to be careful not to remove parts of words.
	// e.g. "please" in "displease" should typically stay, but " please " should go.
	// For simplicity, we'll use word boundaries \b

	for _, filler := range fillers {
		// regex: (?i)\bfiller\b
		pattern := `(?i)\b` + regexp.QuoteMeta(filler) + `\b`
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllString(text, "")
	}

	// Clean up any double spaces left behind
	return collapseRepeatedSpaces(text)
}

// DeduplicateMessages removes consecutive messages with identical content and role.
func DeduplicateMessages(messages []types.Message) []types.Message {
	if len(messages) == 0 {
		return messages
	}

	result := make([]types.Message, 0, len(messages))
	result = append(result, messages[0])

	for i := 1; i < len(messages); i++ {
		curr := messages[i]
		prev := result[len(result)-1]

		if curr.Role == prev.Role && curr.Content == prev.Content {
			continue
		}
		result = append(result, curr)
	}

	return result
}
