package main

import (
	"ccstatus/internal/calculator"
	"ccstatus/internal/formatter"
	"ccstatus/internal/parser"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ModelInfo represents model information from Claude Code
type ModelInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// StatusInput represents JSON input from Claude Code via stdin
type StatusInput struct {
	SessionID      string    `json:"session_id"`
	Cwd            string    `json:"cwd"`
	Model          ModelInfo `json:"model"`
	TranscriptPath string    `json:"transcript_path"`
}

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// run is the main logic, separated for testing
func run(stdin io.Reader, stdout io.Writer) error {
	// read JSON input from stdin
	var input StatusInput
	decoder := json.NewDecoder(stdin)
	if err := decoder.Decode(&input); err != nil {
		return fmt.Errorf("failed to decode input: %w", err)
	}

	// validate input
	if input.TranscriptPath == "" {
		return fmt.Errorf("transcript_path is empty")
	}

	// parse transcript to get usage
	usage, err := parser.ParseTranscript(input.TranscriptPath)
	if err != nil {
		// show explicit error instead of silent degradation
		fmt.Fprint(stdout, formatter.FormatError(fmt.Sprintf("parse error: %v", err)))
		return err
	}

	// extract model name
	model := input.Model.ID
	if model == "" {
		model = "claude"
	}

	// calculate context info with model-specific limits
	info := calculator.Calculate(usage, model)

	// format and output
	output := formatter.Format(info, model)
	fmt.Fprint(stdout, output)

	return nil
}
