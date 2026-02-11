package engine

import (
	"github.com/the-wrong-guy/promptz/core/normalize"
	"github.com/the-wrong-guy/promptz/core/rewrite"
	"github.com/the-wrong-guy/promptz/core/tokenizer"
	"github.com/the-wrong-guy/promptz/core/types"
)

// Optimize processes the request through the efficiency pipeline.
func Optimize(req types.OptimizeRequest) types.OptimizeResponse {
	// Initialize tokenizer
	tok := tokenizer.NewSimpleTokenizer()

	// 1. Calculate tokens before
	tokensBefore := 0
	for _, msg := range req.Messages {
		tokensBefore += tok.CountTokens(msg.Content)
	}

	// 2. Normalize & Rewrite
	optimizedMsgs := make([]types.Message, 0, len(req.Messages))
	normOpts := normalize.DefaultOptions()

	// Pre-process messages
	// We first normalize text, then deduplicate messages if needed.
	// Actually, strict requirement: "deduplicate consecutive identical messages"
	// and "trim whitespace etc".

	// Let's create a temporary slice of normalized messages first
	normalizedMsgs := make([]types.Message, 0, len(req.Messages))
	for _, msg := range req.Messages {
		cleaned := normalize.NormalizeText(msg.Content, normOpts)
		if cleaned == "" {
			// Option: remove empty messages? or keep them?
			// Generally empty messages are useless tokens.
			continue
		}
		normalizedMsgs = append(normalizedMsgs, types.Message{
			Role:    msg.Role,
			Content: cleaned,
		})
	}

	// Deduplicate
	dedupedMsgs := normalize.DeduplicateMessages(normalizedMsgs)

	// Rewrite
	for _, msg := range dedupedMsgs {
		rewritten := rewrite.RewriteMessage(msg, req.Mode)
		if rewritten == "" {
			continue
		}
		optimizedMsgs = append(optimizedMsgs, types.Message{
			Role:    msg.Role,
			Content: rewritten,
		})
	}

	// 3. Tokens after
	tokensAfter := 0
	for _, msg := range optimizedMsgs {
		tokensAfter += tok.CountTokens(msg.Content)
	}

	// 4. Calculate stats
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
