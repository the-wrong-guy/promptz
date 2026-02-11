package engine

import (
	"github.com/the-wrong-guy/promptz/core/nlp/compress"
	"github.com/the-wrong-guy/promptz/core/nlp/similarity"
	"github.com/the-wrong-guy/promptz/core/normalize"
	"github.com/the-wrong-guy/promptz/core/rewrite"
	"github.com/the-wrong-guy/promptz/core/tokenizer"
	"github.com/the-wrong-guy/promptz/core/types"
)

// Optimize processes the request through the NLP-enhanced efficiency pipeline.
//
// Pipeline steps:
//  1. Token count before
//  2. Normalize messages (trim, collapse spaces, remove filler phrases)
//  3. Sentence compression (remove parentheticals, verbose clauses)
//  4. Similarity-based deduplication (near-duplicate removal)
//  5. NLP-enhanced rewrite (POS tagging + TF-IDF scored keyword extraction)
//  6. Token count after
//  7. Return metrics
func Optimize(req types.OptimizeRequest) types.OptimizeResponse {
	tok := tokenizer.NewSimpleTokenizer()

	// 1. Calculate tokens before
	tokensBefore := 0
	for _, msg := range req.Messages {
		tokensBefore += tok.CountTokens(msg.Content)
	}

	// 2. Normalize text
	normOpts := normalize.DefaultOptions()
	normalizedMsgs := make([]types.Message, 0, len(req.Messages))
	for _, msg := range req.Messages {
		cleaned := normalize.NormalizeText(msg.Content, normOpts)
		if cleaned == "" {
			continue
		}
		normalizedMsgs = append(normalizedMsgs, types.Message{
			Role:    msg.Role,
			Content: cleaned,
		})
	}

	// 3. Sentence compression (strip parentheticals, filler clauses)
	compressedMsgs := make([]types.Message, 0, len(normalizedMsgs))
	for _, msg := range normalizedMsgs {
		compressed := compress.Compress(msg.Content)
		if compressed == "" {
			continue
		}
		compressedMsgs = append(compressedMsgs, types.Message{
			Role:    msg.Role,
			Content: compressed,
		})
	}

	// 4. Similarity-based deduplication
	simOpts := similarity.DefaultOptions()
	// Use stricter threshold for conservative mode, looser for aggressive
	switch req.Mode {
	case types.ModeConservative:
		simOpts.Threshold = 0.9 // Only near-exact duplicates
	case types.ModeBalanced:
		simOpts.Threshold = 0.7
	case types.ModeAggressive:
		simOpts.Threshold = 0.5 // Catch more near-duplicates
	}
	dedupedMsgs := similarity.DeduplicateBySimilarity(compressedMsgs, simOpts)

	// 5. NLP-enhanced rewrite with full conversation context
	rewriter := rewrite.NewRewriter(dedupedMsgs)
	optimizedMsgs := make([]types.Message, 0, len(dedupedMsgs))
	for _, msg := range dedupedMsgs {
		rewritten := rewriter.RewriteMessage(msg, req.Mode)
		if rewritten == "" {
			continue
		}
		optimizedMsgs = append(optimizedMsgs, types.Message{
			Role:    msg.Role,
			Content: rewritten,
		})
	}

	// 6. Tokens after
	tokensAfter := 0
	for _, msg := range optimizedMsgs {
		tokensAfter += tok.CountTokens(msg.Content)
	}

	// 7. Calculate stats
	savings := 0.0
	if tokensBefore > 0 {
		savings = float64(tokensBefore-tokensAfter) / float64(tokensBefore)
	}

	return types.OptimizeResponse{
		Optimized:    optimizedMsgs,
		TokensBefore: tokensBefore,
		TokensAfter:  tokensAfter,
		SavingsRatio: savings,
	}
}
