package rewrite

import (
	"strings"

	"github.com/the-wrong-guy/promptz/core/nlp/postagger"
	"github.com/the-wrong-guy/promptz/core/nlp/tfidf"
	"github.com/the-wrong-guy/promptz/core/types"
)

// Rewriter performs NLP-enhanced message rewriting.
type Rewriter struct {
	tagger *postagger.Tagger
	scorer *tfidf.Scorer
}

// NewRewriter creates a Rewriter with NLP context from all messages.
// The TF-IDF scorer uses all messages as the document corpus.
func NewRewriter(messages []types.Message) *Rewriter {
	docs := make([]string, len(messages))
	for i, m := range messages {
		docs[i] = m.Content
	}
	return &Rewriter{
		tagger: postagger.New(),
		scorer: tfidf.New(docs),
	}
}

// RewriteMessage transforms message content based on the optimization mode.
// Uses prose-based POS tagging for balanced and aggressive modes.
func (r *Rewriter) RewriteMessage(msg types.Message, mode types.Mode) string {
	content := msg.Content

	switch mode {
	case types.ModeConservative:
		return content

	case types.ModeBalanced:
		return r.rewriteBalanced(content)

	case types.ModeAggressive:
		return r.rewriteAggressive(content)

	default:
		return content
	}
}

// rewriteBalanced uses prose POS tagging to keep content words and important adverbs,
// dropping determiners, prepositions, and conjunctions.
func (r *Rewriter) rewriteBalanced(text string) string {
	tagged := r.tagger.TagSentence(text)
	var kept []string

	for _, tw := range tagged {
		switch tw.Tag {
		case postagger.TagDeterminer, postagger.TagPreposition, postagger.TagConjunction:
			continue
		case postagger.TagPronoun:
			lower := strings.ToLower(tw.Word)
			if lower == "i" || lower == "you" || lower == "we" {
				kept = append(kept, tw.Word)
			}
		default:
			kept = append(kept, tw.Word)
		}
	}

	return strings.Join(kept, " ")
}

// rewriteAggressive drops all function words (determiners, prepositions,
// conjunctions, non-essential pronouns, non-essential adverbs).
// Unknown words are kept by default.
func (r *Rewriter) rewriteAggressive(text string) string {
	tagged := r.tagger.TagSentence(text)
	var kept []string

	for _, tw := range tagged {
		lower := strings.ToLower(tw.Word)

		switch tw.Tag {
		case postagger.TagDeterminer:
			continue
		case postagger.TagPreposition:
			continue
		case postagger.TagConjunction:
			continue
		case postagger.TagPronoun:
			if lower == "i" || lower == "you" || lower == "we" {
				kept = append(kept, tw.Word)
			}
			continue
		case postagger.TagAdverb:
			if importantAdverbs[lower] {
				kept = append(kept, tw.Word)
			}
			continue
		default:
			kept = append(kept, tw.Word)
		}
	}

	return strings.Join(kept, " ")
}

var importantAdverbs = map[string]bool{
	"not": true, "never": true, "always": true,
	"up": true, "down": true, "out": true, "off": true,
	"only": true, "also": true, "still": true,
	"how": true, "why": true, "when": true, "where": true,
}

// --- Legacy API for backward compatibility ---

// RewriteMessage is the legacy function that creates a one-off Rewriter.
// Prefer using NewRewriter + r.RewriteMessage for better TF-IDF context.
func RewriteMessage(msg types.Message, mode types.Mode) string {
	r := NewRewriter([]types.Message{msg})
	return r.RewriteMessage(msg, mode)
}
