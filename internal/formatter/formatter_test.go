package formatter

import (
	"ccstatus/internal/calculator"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name         string
		info         calculator.ContextInfo
		model        string
		wantContains []string
		wantColor    string
	}{
		{
			name: "green usage",
			info: calculator.ContextInfo{
				CurrentTokens: 50000,
				MaxTokens:     200000,
				Percentage:    25.0,
			},
			model:        "claude-sonnet-4-5",
			wantContains: []string{"ctx:", "50000", "200000", "25.0%", "claude-sonnet-4-5"},
			wantColor:    ColorGreen,
		},
		{
			name: "yellow usage",
			info: calculator.ContextInfo{
				CurrentTokens: 140000,
				MaxTokens:     200000,
				Percentage:    70.0,
			},
			model:        "claude-sonnet-4-5",
			wantContains: []string{"ctx:", "140000", "200000", "70.0%"},
			wantColor:    ColorYellow,
		},
		{
			name: "red usage",
			info: calculator.ContextInfo{
				CurrentTokens: 180000,
				MaxTokens:     200000,
				Percentage:    90.0,
			},
			model:        "claude-sonnet-4-5",
			wantContains: []string{"ctx:", "180000", "200000", "90.0%"},
			wantColor:    ColorRed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.info, tt.model)

			if !strings.Contains(got, tt.wantColor) {
				t.Errorf("Format() does not contain expected color %q, got %q", tt.wantColor, got)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("Format() does not contain %q, got %q", want, got)
				}
			}

			if !strings.Contains(got, ColorReset) {
				t.Errorf("Format() does not contain color reset")
			}
		})
	}
}

func TestFormatCompact(t *testing.T) {
	tests := []struct {
		name      string
		info      calculator.ContextInfo
		wantColor string
	}{
		{
			name: "compact green",
			info: calculator.ContextInfo{
				CurrentTokens: 50000,
				MaxTokens:     200000,
				Percentage:    25.0,
			},
			wantColor: ColorGreen,
		},
		{
			name: "compact yellow",
			info: calculator.ContextInfo{
				CurrentTokens: 140000,
				MaxTokens:     200000,
				Percentage:    70.0,
			},
			wantColor: ColorYellow,
		},
		{
			name: "compact red",
			info: calculator.ContextInfo{
				CurrentTokens: 180000,
				MaxTokens:     200000,
				Percentage:    90.0,
			},
			wantColor: ColorRed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatCompact(tt.info)

			if !strings.Contains(got, tt.wantColor) {
				t.Errorf("FormatCompact() does not contain expected color %q", tt.wantColor)
			}

			if !strings.Contains(got, "[") || !strings.Contains(got, "%]") {
				t.Errorf("FormatCompact() format incorrect, got %q", got)
			}
		})
	}
}

func TestFormatPlain(t *testing.T) {
	tests := []struct {
		name         string
		info         calculator.ContextInfo
		model        string
		wantContains []string
	}{
		{
			name: "plain format",
			info: calculator.ContextInfo{
				CurrentTokens: 59261,
				MaxTokens:     200000,
				Percentage:    29.6305,
			},
			model:        "claude-sonnet-4-5",
			wantContains: []string{"ctx:", "59261", "200000", "29.6%", "claude-sonnet-4-5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatPlain(tt.info, tt.model)

			// should not contain any color codes
			if strings.Contains(got, "\033[") {
				t.Errorf("FormatPlain() contains color codes, got %q", got)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("FormatPlain() does not contain %q, got %q", want, got)
				}
			}
		})
	}
}

func TestGetColor(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  string
	}{
		{
			name:  "green level",
			level: "green",
			want:  ColorGreen,
		},
		{
			name:  "yellow level",
			level: "yellow",
			want:  ColorYellow,
		},
		{
			name:  "red level",
			level: "red",
			want:  ColorRed,
		},
		{
			name:  "unknown level defaults to reset",
			level: "unknown",
			want:  ColorReset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getColor(tt.level)
			if got != tt.want {
				t.Errorf("getColor(%q) = %q, want %q", tt.level, got, tt.want)
			}
		})
	}
}

// benchmark formatter performance
func BenchmarkFormat(b *testing.B) {
	info := calculator.ContextInfo{
		CurrentTokens: 59261,
		MaxTokens:     200000,
		Percentage:    29.6305,
	}
	model := "claude-sonnet-4-5"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Format(info, model)
	}
}

func BenchmarkFormatCompact(b *testing.B) {
	info := calculator.ContextInfo{
		CurrentTokens: 59261,
		MaxTokens:     200000,
		Percentage:    29.6305,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatCompact(info)
	}
}
