# Token Efficiency Engine (promptz)

A production-grade Go library and CLI for optimizing LLM prompts. `promptz` reduces token usage while identifying and preserving semantic meaning through deterministic NLP techniques.

## Features

- **Smart Normalization**: Trims whitespace, collapses spans, and removes low-value filler phrases ("in order to", "as a matter of fact").
- **NLP-Enhanced Rewrite**:
  - **POS Tagging**: Uses `jdkato/prose/v2` for accurate Part-of-Speech tagging.
  - **TF-IDF Scoring**: Ranks word importance across the conversation context to preserve key terms.
  - **Sentence Compression**: Strips parenthetical asides and verbose clauses.
- **Semantic Deduplication**: Uses Jaccard similarity to detect and merge near-duplicate messages (e.g., "fix the bug" vs "please fix that bug").
- **Deterministic**: No external AI calls required. Runs entirely locally.

## Modes

- `conservative`: Minimal changes. Safe for all prompts.
- `balanced`: Removes stop words and common filler phrases.
- `aggressive`: Retains only high-value content words (nouns, verbs, adjectives, important adverbs).

## Installation

```bash
go get github.com/the-wrong-guy/promptz
```

## CLI Usage

Build the tool:
```bash
go build -o distill ./cmd/distill
```

Run with JSON input:
```bash
echo '{"messages": [{"role": "user", "content": "Hello my world, whats up"}], "mode": "aggressive"}' | ./distill
```

Output:
```json
{
  "optimized": [
    {
      "role": "user",
      "content": "Hello world whats up"
    }
  ],
  "tokens_before": 6,
  "tokens_after": 4,
  "savings_ratio": 0.33
}
```

## Library Usage

```go
package main

import (
	"fmt"
	"github.com/the-wrong-guy/promptz/core/engine"
	"github.com/the-wrong-guy/promptz/core/types"
)

func main() {
	req := types.OptimizeRequest{
		Messages: []types.Message{
			{Role: "user", Content: "I am facing a database connection error in production"},
		},
		Mode: types.ModeAggressive,
	}

	resp := engine.Optimize(req)

	for _, msg := range resp.Optimized {
		fmt.Printf("%s: %s\n", msg.Role, msg.Content)
	}
	// Output: user: database connection error production
}
```

## Architecture

The engine runs a 7-step optimization pipeline:
1. **Token Count** (Pre-optimization)
2. **Normalize** (Whitespace, filler phrases)
3. **Compress** (Strip parentheticals, verbose patterns)
4. **Similarity Dedup** (Merge near-duplicates)
5. **NLP Rewrite** (POS Tagging + TF-IDF)
6. **Token Count** (Post-optimization)
7. **Metrics Calculation**

## Contributing

We welcome contributions to make `promptz` even better!

**How to contribute:**
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Development:**
- Run tests: `go test -v ./...`
- Ensure code is idiomatic Go and formatted with `gofmt`.

## License

MIT
