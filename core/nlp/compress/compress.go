package compress

import (
	"regexp"
	"strings"
)

// Compress removes low-value clauses and redundant structures from text.
func Compress(text string) string {
	text = removeParentheticals(text)
	text = removeFillerClauses(text)
	text = collapseSpaces(text)
	return strings.TrimSpace(text)
}

// removeParentheticals strips content within parentheses and dashes.
// e.g. "fix the bug (which is critical) now" → "fix the bug now"
func removeParentheticals(text string) string {
	// Remove (...) content
	re := regexp.MustCompile(`\s*\([^)]*\)\s*`)
	text = re.ReplaceAllString(text, " ")

	// Remove -- ... -- content (em-dash asides)
	re2 := regexp.MustCompile(`\s*--\s*[^-]*--\s*`)
	text = re2.ReplaceAllString(text, " ")

	return text
}

// removeFillerClauses removes common verbose clause patterns.
var fillerClauses = []string{
	"which is",
	"that is",
	"in order to",
	"as a matter of fact",
	"it should be noted that",
	"it is important to note that",
	"what i mean is",
	"the thing is",
	"as you know",
	"as we know",
	"to be honest",
	"in my opinion",
	"i think that",
	"i believe that",
	"i feel like",
	"it seems like",
	"it seems that",
	"basically what",
	"essentially",
	"at this point in time",
	"at the end of the day",
	"due to the fact that",
	"in the event that",
	"for the purpose of",
	"with regard to",
	"with respect to",
	"in terms of",
	"on the other hand",
	"having said that",
	"that being said",
	"needless to say",
}

func removeFillerClauses(text string) string {
	for _, clause := range fillerClauses {
		pattern := `(?i)\b` + regexp.QuoteMeta(clause) + `\b`
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllString(text, "")
	}
	return text
}

func collapseSpaces(text string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(text, " ")
}
