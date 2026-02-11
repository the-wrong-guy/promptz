package tfidf

import (
	"math"
	"strings"
)

// Scorer computes TF-IDF scores for words across a set of documents (messages).
type Scorer struct {
	// docFreq tracks how many documents contain each word.
	docFreq map[string]int
	// totalDocs is the total number of documents.
	totalDocs int
}

// New creates a Scorer from a slice of documents (message contents).
func New(documents []string) *Scorer {
	s := &Scorer{
		docFreq:   make(map[string]int),
		totalDocs: len(documents),
	}

	for _, doc := range documents {
		seen := make(map[string]bool)
		for _, word := range strings.Fields(strings.ToLower(doc)) {
			clean := strings.TrimRight(word, ".,!?:;\"'()[]{}")
			if clean != "" && !seen[clean] {
				seen[clean] = true
				s.docFreq[clean]++
			}
		}
	}

	return s
}

// Score returns the TF-IDF score for a word within a specific document.
// Higher scores mean the word is more important/unique in context.
func (s *Scorer) Score(word string, document string) float64 {
	lower := strings.ToLower(word)

	// Term Frequency: count of word in document / total words in document
	words := strings.Fields(strings.ToLower(document))
	if len(words) == 0 {
		return 0
	}

	count := 0
	for _, w := range words {
		clean := strings.TrimRight(w, ".,!?:;\"'()[]{}")
		if clean == lower {
			count++
		}
	}
	tf := float64(count) / float64(len(words))

	// Inverse Document Frequency: log(totalDocs / docFreq)
	df := s.docFreq[lower]
	if df == 0 {
		return 0
	}
	idf := math.Log(float64(s.totalDocs+1) / float64(df+1))

	return tf * idf
}

// ScoreWords returns a map of word -> TF-IDF score for all words in a document.
func (s *Scorer) ScoreWords(document string) map[string]float64 {
	scores := make(map[string]float64)
	words := strings.Fields(strings.ToLower(document))

	seen := make(map[string]bool)
	for _, w := range words {
		clean := strings.TrimRight(w, ".,!?:;\"'()[]{}")
		if clean != "" && !seen[clean] {
			seen[clean] = true
			scores[clean] = s.Score(clean, document)
		}
	}

	return scores
}

// TopN returns the top N words by TF-IDF score for a document.
func (s *Scorer) TopN(document string, n int) []string {
	scores := s.ScoreWords(document)

	type wordScore struct {
		word  string
		score float64
	}

	sorted := make([]wordScore, 0, len(scores))
	for w, sc := range scores {
		sorted = append(sorted, wordScore{w, sc})
	}

	// Simple selection sort for small lists (no need for sort package dependency)
	for i := 0; i < len(sorted) && i < n; i++ {
		maxIdx := i
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].score > sorted[maxIdx].score {
				maxIdx = j
			}
		}
		sorted[i], sorted[maxIdx] = sorted[maxIdx], sorted[i]
	}

	result := make([]string, 0, n)
	for i := 0; i < len(sorted) && i < n; i++ {
		result = append(result, sorted[i].word)
	}
	return result
}
