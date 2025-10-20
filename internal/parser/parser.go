package parser

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Usage represents token usage statistics from Claude API
type Usage struct {
	InputTokens            int64                  `json:"input_tokens"`
	CacheReadInputTokens   int64                  `json:"cache_read_input_tokens"`
	CacheCreationInputTokens int64                `json:"cache_creation_input_tokens"`
	OutputTokens           int64                  `json:"output_tokens"`
	CacheCreation          map[string]int         `json:"cache_creation,omitempty"`
}

// Message represents a single message in the JSONL transcript
type Message struct {
	Message struct {
		Role  string `json:"role"`
		Usage Usage  `json:"usage"`
	} `json:"message"`
}

// ParseTranscript reads a JSONL transcript file and returns the last message usage data
// Returns error if file cannot be read or no messages with usage found
func ParseTranscript(transcriptPath string) (*Usage, error) {
	if transcriptPath == "" {
		return nil, errors.New("transcript path is empty")
	}

	// resolve to absolute path to prevent path traversal
	absPath, err := filepath.Abs(transcriptPath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// check for suspicious patterns (parent directory references)
	if strings.Contains(filepath.ToSlash(absPath), "..") {
		return nil, errors.New("invalid path: contains parent directory references")
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open transcript: %w", err)
	}
	defer file.Close()

	return parseTranscriptFromReader(file)
}

// parseTranscriptFromReader parses transcript from io.Reader
// separated for testing purposes
func parseTranscriptFromReader(r io.Reader) (*Usage, error) {
	var lastUsage *Usage
	scanner := bufio.NewScanner(r)

	// increase buffer size for large lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			// skip malformed lines
			continue
		}

		// accept any message with usage data, regardless of role
		// this catches user prompts and tool calls that may have usage info
		if msg.Message.Role != "" && hasValidUsage(&msg.Message.Usage) {
			// copy to avoid pointer to loop variable issue
			usageCopy := msg.Message.Usage
			lastUsage = &usageCopy
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading transcript: %w", err)
	}

	if lastUsage == nil {
		// return zero usage instead of error for empty transcripts
		return &Usage{
			InputTokens:              0,
			CacheReadInputTokens:     0,
			CacheCreationInputTokens: 0,
			OutputTokens:             0,
		}, nil
	}

	return lastUsage, nil
}

// hasValidUsage checks if usage struct contains meaningful data
func hasValidUsage(usage *Usage) bool {
	if usage == nil {
		return false
	}
	// usage is valid if at least one token field is non-zero
	return usage.InputTokens > 0 ||
		usage.CacheReadInputTokens > 0 ||
		usage.CacheCreationInputTokens > 0 ||
		usage.OutputTokens > 0
}
