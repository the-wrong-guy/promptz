package tfidf

import "testing"

func TestScorer_Score(t *testing.T) {
	docs := []string{
		"the quick brown fox",
		"the lazy brown dog",
		"quick quick quick fox",
	}
	scorer := New(docs)

	// "the" appears in 2/3 docs, "quick" appears in 2/3 docs,
	// "fox" appears in 2/3 docs.
	// In doc[2], "quick" has high TF (3/4), moderate IDF
	quickScore := scorer.Score("quick", docs[2])
	theScore := scorer.Score("the", docs[0])

	// "quick" in doc[2] should score higher than "the" in doc[0]
	// because "quick" has higher TF in doc[2]
	if quickScore <= theScore {
		t.Errorf("expected 'quick' in doc[2] (%.4f) > 'the' in doc[0] (%.4f)", quickScore, theScore)
	}
}

func TestScorer_TopN(t *testing.T) {
	docs := []string{
		"fix the database connection error in production",
		"the server is running fine",
	}
	scorer := New(docs)

	top := scorer.TopN(docs[0], 3)
	if len(top) == 0 {
		t.Fatal("expected at least one top word")
	}

	// Words unique to doc[0] should rank highest
	// "fix", "database", "connection", "error", "production" are unique to doc[0]
	// "the" and "in"/"is" appear in both
	found := false
	for _, w := range top {
		if w == "database" || w == "connection" || w == "error" || w == "fix" || w == "production" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected unique words in top-3, got %v", top)
	}
}

func TestScorer_EmptyDocument(t *testing.T) {
	scorer := New([]string{"hello world", ""})
	score := scorer.Score("hello", "")
	if score != 0 {
		t.Errorf("expected score 0 for empty document, got %f", score)
	}
}
