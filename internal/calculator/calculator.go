package calculator

import (
	"ccstatus/internal/parser"
	"strings"
)

// model context limits in tokens
var modelLimits = map[string]int64{
	"claude-3-opus":           200000,
	"claude-3-sonnet":         200000,
	"claude-3-haiku":          200000,
	"claude-3-5-sonnet":       200000,
	"claude-3-5-haiku":        200000,
	"claude-3-5-opus":         200000,
	"claude-sonnet-4":         200000,
	"claude-sonnet-4-5":       200000,
	"claude-2.1":              200000,
	"claude-2":                100000,
	"claude-instant-1.2":      100000,
	"claude-instant-1":        100000,
}

const (
	// DefaultContextTokens is fallback context window limit
	DefaultContextTokens = 200000
)

// ContextInfo contains calculated context usage information
type ContextInfo struct {
	CurrentTokens int64
	MaxTokens     int64
	Percentage    float64
}

// Calculate computes context usage from parsed usage data
// Formula (corrected): current_context = input_tokens + cache_read_input_tokens
// This properly accounts for both new tokens and cached tokens
func Calculate(usage *parser.Usage, model string) ContextInfo {
	maxTokens := getModelLimit(model)

	if usage == nil {
		return ContextInfo{
			CurrentTokens: 0,
			MaxTokens:     maxTokens,
			Percentage:    0,
		}
	}

	// correct formula: input_tokens includes all non-cached tokens
	// cache_read_input_tokens includes all cached tokens being read
	currentTokens := usage.InputTokens + usage.CacheReadInputTokens
	percentage := (float64(currentTokens) / float64(maxTokens)) * 100

	// clamp percentage to avoid >100% display issues
	if percentage > 100.0 {
		percentage = 100.0
	}

	return ContextInfo{
		CurrentTokens: currentTokens,
		MaxTokens:     maxTokens,
		Percentage:    percentage,
	}
}

// getModelLimit returns context window limit for given model
func getModelLimit(model string) int64 {
	// try exact match first
	if limit, ok := modelLimits[model]; ok {
		return limit
	}

	// try prefix match (e.g., "claude-3-opus-20240229" matches "claude-3-opus")
	modelLower := strings.ToLower(model)
	for prefix, limit := range modelLimits {
		if strings.HasPrefix(modelLower, prefix) {
			return limit
		}
	}

	// fallback to default
	return DefaultContextTokens
}

// GetUsageLevel returns usage level based on percentage
// green: 0-60%, yellow: 60-80%, red: 80-100%
func GetUsageLevel(percentage float64) string {
	switch {
	case percentage < 60:
		return "green"
	case percentage < 80:
		return "yellow"
	default:
		return "red"
	}
}
