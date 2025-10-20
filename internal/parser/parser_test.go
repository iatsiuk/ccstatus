package parser

import (
	"strings"
	"testing"
)

func TestParseTranscriptFromReader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Usage
		wantErr bool
	}{
		{
			name: "valid transcript with single assistant message",
			input: `{"message":{"role":"assistant","usage":{"input_tokens":9,"cache_read_input_tokens":58164,"cache_creation_input_tokens":1097,"output_tokens":2}}}`,
			want: &Usage{
				InputTokens:              9,
				CacheReadInputTokens:     58164,
				CacheCreationInputTokens: 1097,
				OutputTokens:             2,
			},
			wantErr: false,
		},
		{
			name: "valid transcript with multiple messages, returns last",
			input: `{"message":{"role":"assistant","usage":{"input_tokens":5,"cache_read_input_tokens":1000,"cache_creation_input_tokens":500,"output_tokens":10}}}
{"message":{"role":"user","content":"test"}}
{"message":{"role":"assistant","usage":{"input_tokens":9,"cache_read_input_tokens":2000,"cache_creation_input_tokens":600,"output_tokens":20}}}`,
			want: &Usage{
				InputTokens:              9,
				CacheReadInputTokens:     2000,
				CacheCreationInputTokens: 600,
				OutputTokens:             20,
			},
			wantErr: false,
		},
		{
			name: "transcript with malformed json lines - skips them",
			input: `{"invalid json
{"message":{"role":"assistant","usage":{"input_tokens":9,"cache_read_input_tokens":1500,"cache_creation_input_tokens":300,"output_tokens":5}}}
not json at all`,
			want: &Usage{
				InputTokens:              9,
				CacheReadInputTokens:     1500,
				CacheCreationInputTokens: 300,
				OutputTokens:             5,
			},
			wantErr: false,
		},
		{
			name:  "empty transcript",
			input: "",
			want: &Usage{
				InputTokens:              0,
				CacheReadInputTokens:     0,
				CacheCreationInputTokens: 0,
				OutputTokens:             0,
			},
			wantErr: false,
		},
		{
			name:  "no assistant messages",
			input: `{"message":{"role":"user","content":"hello"}}`,
			want: &Usage{
				InputTokens:              0,
				CacheReadInputTokens:     0,
				CacheCreationInputTokens: 0,
				OutputTokens:             0,
			},
			wantErr: false,
		},
		{
			name:  "only malformed json",
			input: `{invalid}`,
			want: &Usage{
				InputTokens:              0,
				CacheReadInputTokens:     0,
				CacheCreationInputTokens: 0,
				OutputTokens:             0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := parseTranscriptFromReader(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseTranscriptFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Fatal("parseTranscriptFromReader() returned nil, want valid Usage")
				}

				if got.InputTokens != tt.want.InputTokens {
					t.Errorf("InputTokens = %v, want %v", got.InputTokens, tt.want.InputTokens)
				}
				if got.CacheReadInputTokens != tt.want.CacheReadInputTokens {
					t.Errorf("CacheReadInputTokens = %v, want %v", got.CacheReadInputTokens, tt.want.CacheReadInputTokens)
				}
				if got.CacheCreationInputTokens != tt.want.CacheCreationInputTokens {
					t.Errorf("CacheCreationInputTokens = %v, want %v", got.CacheCreationInputTokens, tt.want.CacheCreationInputTokens)
				}
				if got.OutputTokens != tt.want.OutputTokens {
					t.Errorf("OutputTokens = %v, want %v", got.OutputTokens, tt.want.OutputTokens)
				}
			}
		})
	}
}

func TestParseTranscript(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "non-existent file",
			path:    "/nonexistent/path/file.jsonl",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTranscript(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTranscript() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseTranscriptZeroUsage(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "no messages with usage returns zero",
			input: `{"message":{"role":"user","content":"test"}}`,
		},
		{
			name:  "messages with all-zero usage returns zero",
			input: `{"message":{"role":"assistant","usage":{"input_tokens":0,"cache_read_input_tokens":0,"cache_creation_input_tokens":0,"output_tokens":0}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := parseTranscriptFromReader(reader)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got == nil {
				t.Fatal("parseTranscriptFromReader() returned nil, want zero Usage")
			}

			if got.InputTokens != 0 || got.CacheReadInputTokens != 0 || got.CacheCreationInputTokens != 0 || got.OutputTokens != 0 {
				t.Errorf("expected zero usage, got %+v", got)
			}
		})
	}
}

// benchmark parsing performance
func BenchmarkParseTranscriptFromReader(b *testing.B) {
	// realistic transcript with 10 messages
	var sb strings.Builder
	for i := 0; i < 10; i++ {
		sb.WriteString(`{"message":{"role":"assistant","usage":{"input_tokens":9,"cache_read_input_tokens":58164,"cache_creation_input_tokens":1097,"output_tokens":2}}}`)
		sb.WriteString("\n")
	}
	transcript := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(transcript)
		_, _ = parseTranscriptFromReader(reader)
	}
}
