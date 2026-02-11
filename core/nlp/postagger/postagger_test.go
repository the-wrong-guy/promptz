package postagger

import "testing"

func TestTagger_TagSentence(t *testing.T) {
	tagger := New()

	tests := []struct {
		name     string
		sentence string
		checks   map[string]Tag // word -> expected tag
	}{
		{
			name:     "determiners detected",
			sentence: "The big dog runs in the park",
			checks: map[string]Tag{
				"The": TagDeterminer,
			},
		},
		{
			name:     "technical text nouns",
			sentence: "fix the database connection error",
			checks: map[string]Tag{
				"database":   TagNoun,
				"connection": TagNoun,
				"error":      TagNoun,
			},
		},
		{
			name:     "prepositions detected",
			sentence: "the cat is in the house",
			checks: map[string]Tag{
				"in": TagPreposition,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagged := tagger.TagSentence(tt.sentence)

			tagMap := make(map[string]Tag)
			for _, tw := range tagged {
				tagMap[tw.Word] = tw.Tag
			}

			for word, expectedTag := range tt.checks {
				if got, ok := tagMap[word]; !ok {
					t.Errorf("word %q not found in tagged output", word)
				} else if got != expectedTag {
					t.Errorf("word %q = %v, want %v", word, got, expectedTag)
				}
			}
		})
	}
}

func TestTagger_WordsNotLost(t *testing.T) {
	tagger := New()

	// Verify that prose produces tokens for all meaningful words
	sentence := "Hello my world whats up"
	tagged := tagger.TagSentence(sentence)

	if len(tagged) < 4 {
		t.Errorf("expected at least 4 tagged words for %q, got %d", sentence, len(tagged))
		for _, tw := range tagged {
			t.Logf("  %s -> %s", tw.Word, tw.Tag)
		}
	}
}

func TestIsContentWord(t *testing.T) {
	if !IsContentWord(TagNoun) {
		t.Error("Noun should be a content word")
	}
	if !IsContentWord(TagVerb) {
		t.Error("Verb should be a content word")
	}
	if IsContentWord(TagDeterminer) {
		t.Error("Determiner should NOT be a content word")
	}
	if IsContentWord(TagPreposition) {
		t.Error("Preposition should NOT be a content word")
	}
}
