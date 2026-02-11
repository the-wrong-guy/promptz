# Token Efficiency Engine (promptz)

An open-source library and CLI for optimizing chat prompts by removing low-value tokens and restructuring context while preserving semantic meaning.

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
echo '{"messages": [{"role": "user", "content": "Hello world"}], "mode": "balanced"}' | ./distill
```

## Features

- **Normalization**: Trims whitespace, removes filler phrases, deduplicates messages.
- **Rewrite Modes**:
  - `conservative`: Minimal changes.
  - `balanced`: Removes stop words.
  - `aggressive`: Keywords only.
- **Deterministic**: No AI calls involved.

## Library Usage

```go
import (
    "github.com/the-wrong-guy/promptz/core/engine"
    "github.com/the-wrong-guy/promptz/core/types"
)

req := types.OptimizeRequest{
    Messages: []types.Message{{Role: "user", Content: "Hello world"}},
    Mode:     types.ModeBalanced,
}

resp := engine.Optimize(req)
// resp.Optimized contains the processed messages
```
