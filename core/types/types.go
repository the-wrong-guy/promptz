package types

// Mode defines the optimization aggressiveness level.
type Mode string

const (
	// ModeConservative aims to remove only completely redundant tokens.
	ModeConservative Mode = "conservative"
	// ModeBalanced aims for a good trade-off between tokne savings and clarity.
	ModeBalanced Mode = "balanced"
	// ModeAggressive aims for maximum token reduction, potentially sacrificing some nuance.
	ModeAggressive Mode = "aggressive"
)

// Message represents a chat message in the conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OptimizeRequest defines the input for the optimization engine.
type OptimizeRequest struct {
	Messages []Message `json:"messages"`
	Mode     Mode      `json:"mode"`
}

// OptimizeResponse defines the output of the optimization engine.
type OptimizeResponse struct {
	Optimized    []Message `json:"optimized"`
	TokensBefore int       `json:"tokens_before"`
	TokensAfter  int       `json:"tokens_after"`
	SavingsRatio float64   `json:"savings_ratio"`
}
