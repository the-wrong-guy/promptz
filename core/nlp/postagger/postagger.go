// Package postagger wraps jdkato/prose/v2 to provide POS tagging for the engine.
package postagger

import (
	"strings"

	"github.com/jdkato/prose/v2"
)

// Tag represents a simplified POS category derived from Penn Treebank tags.
type Tag string

const (
	TagNoun        Tag = "NOUN"
	TagVerb        Tag = "VERB"
	TagAdjective   Tag = "ADJ"
	TagAdverb      Tag = "ADV"
	TagPronoun     Tag = "PRON"
	TagDeterminer  Tag = "DET"
	TagPreposition Tag = "PREP"
	TagConjunction Tag = "CONJ"
	TagNumber      Tag = "NUM"
	TagParticle    Tag = "PART"
	TagUnknown     Tag = "UNK"
)

// TaggedWord holds a word token and its simplified POS tag.
type TaggedWord struct {
	Word string
	Tag  Tag
}

// Tagger wraps prose's document-level POS tagger.
type Tagger struct{}

// New creates a new Tagger instance.
func New() *Tagger {
	return &Tagger{}
}

// TagSentence tags all words in a sentence using prose's trained POS model.
func (t *Tagger) TagSentence(sentence string) []TaggedWord {
	doc, err := prose.NewDocument(sentence,
		prose.WithSegmentation(false),
		prose.WithExtraction(false),
	)
	if err != nil {
		// Fallback: return each word as unknown
		words := strings.Fields(sentence)
		result := make([]TaggedWord, len(words))
		for i, w := range words {
			result[i] = TaggedWord{Word: w, Tag: TagUnknown}
		}
		return result
	}

	tokens := doc.Tokens()
	result := make([]TaggedWord, 0, len(tokens))
	for _, tok := range tokens {
		// Skip pure punctuation tokens
		if isPunctuation(tok.Text) {
			continue
		}
		result = append(result, TaggedWord{
			Word: tok.Text,
			Tag:  mapPennTag(tok.Tag),
		})
	}
	return result
}

// TagWord tags a single word by creating a minimal document.
func (t *Tagger) TagWord(word string) Tag {
	tagged := t.TagSentence(word)
	if len(tagged) > 0 {
		return tagged[0].Tag
	}
	return TagUnknown
}

// IsContentWord returns true if the tag is a noun, verb, or adjective.
func IsContentWord(tag Tag) bool {
	switch tag {
	case TagNoun, TagVerb, TagAdjective:
		return true
	default:
		return false
	}
}

// mapPennTag converts Penn Treebank tags (used by prose) to our simplified tags.
func mapPennTag(penn string) Tag {
	switch penn {
	// Nouns
	case "NN", "NNS", "NNP", "NNPS":
		return TagNoun

	// Verbs
	case "VB", "VBD", "VBG", "VBN", "VBP", "VBZ", "MD":
		return TagVerb

	// Adjectives
	case "JJ", "JJR", "JJS":
		return TagAdjective

	// Adverbs
	case "RB", "RBR", "RBS", "WRB":
		return TagAdverb

	// Pronouns
	case "PRP", "PRP$", "WP", "WP$":
		return TagPronoun

	// Determiners
	case "DT", "WDT", "PDT":
		return TagDeterminer

	// Prepositions / subordinating conjunctions
	case "IN", "TO":
		return TagPreposition

	// Coordinating conjunctions
	case "CC":
		return TagConjunction

	// Numbers
	case "CD":
		return TagNumber

	// Particles
	case "RP":
		return TagParticle

	default:
		return TagUnknown
	}
}

func isPunctuation(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}
