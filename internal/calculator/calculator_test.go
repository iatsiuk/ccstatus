package calculator

import (
	"ccstatus/internal/parser"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name          string
		usage         *parser.Usage
		model         string
		wantTokens    int64
		wantMax       int64
		wantPercentage float64
	}{
		{
			name: "typical usage with corrected formula",
			usage: &parser.Usage{
				InputTokens:              9,
				CacheReadInputTokens:     58164,
				CacheCreationInputTokens: 1097,
				OutputTokens:             2,
			},
			model:         "claude-sonnet-4-5-20250929",
			wantTokens:    58173, // input + cache_read (9 + 58164)
			wantMax:       200000,
			wantPercentage: 29.0865,
		},
		{
			name: "zero usage",
			usage: &parser.Usage{
				InputTokens:              0,
				CacheReadInputTokens:     0,
				CacheCreationInputTokens: 0,
				OutputTokens:             0,
			},
			model:         "claude-3-opus",
			wantTokens:    0,
			wantMax:       200000,
			wantPercentage: 0,
		},
		{
			name: "high usage near limit",
			usage: &parser.Usage{
				InputTokens:              10,
				CacheReadInputTokens:     160000,
				CacheCreationInputTokens: 30000,
				OutputTokens:             100,
			},
			model:         "claude-3-sonnet",
			wantTokens:    160010, // input + cache_read (10 + 160000)
			wantMax:       200000,
			wantPercentage: 80.005,
		},
		{
			name:          "nil usage",
			usage:         nil,
			model:         "claude-3-haiku",
			wantTokens:    0,
			wantMax:       200000,
			wantPercentage: 0,
		},
		{
			name: "claude-2 with 100k limit",
			usage: &parser.Usage{
				InputTokens:              100,
				CacheReadInputTokens:     50000,
				CacheCreationInputTokens: 0,
				OutputTokens:             50,
			},
			model:         "claude-2",
			wantTokens:    50100,
			wantMax:       100000,
			wantPercentage: 50.1,
		},
		{
			name: "over 100% clamped",
			usage: &parser.Usage{
				InputTokens:              50000,
				CacheReadInputTokens:     160000,
				CacheCreationInputTokens: 0,
				OutputTokens:             0,
			},
			model:         "claude-3-opus",
			wantTokens:    210000,
			wantMax:       200000,
			wantPercentage: 100.0, // clamped from 105%
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Calculate(tt.usage, tt.model)

			if got.CurrentTokens != tt.wantTokens {
				t.Errorf("Calculate().CurrentTokens = %v, want %v", got.CurrentTokens, tt.wantTokens)
			}
			if got.MaxTokens != tt.wantMax {
				t.Errorf("Calculate().MaxTokens = %v, want %v", got.MaxTokens, tt.wantMax)
			}

			// allow small floating point difference
			diff := got.Percentage - tt.wantPercentage
			if diff < 0 {
				diff = -diff
			}
			if diff > 0.0001 {
				t.Errorf("Calculate().Percentage = %v, want %v", got.Percentage, tt.wantPercentage)
			}
		})
	}
}

func TestGetUsageLevel(t *testing.T) {
	tests := []struct {
		name       string
		percentage float64
		want       string
	}{
		{
			name:       "zero percent",
			percentage: 0,
			want:       "green",
		},
		{
			name:       "low usage 30%",
			percentage: 30,
			want:       "green",
		},
		{
			name:       "medium usage at boundary 59.9%",
			percentage: 59.9,
			want:       "green",
		},
		{
			name:       "medium usage at boundary 60%",
			percentage: 60,
			want:       "yellow",
		},
		{
			name:       "medium usage 70%",
			percentage: 70,
			want:       "yellow",
		},
		{
			name:       "high usage at boundary 79.9%",
			percentage: 79.9,
			want:       "yellow",
		},
		{
			name:       "high usage at boundary 80%",
			percentage: 80,
			want:       "red",
		},
		{
			name:       "high usage 90%",
			percentage: 90,
			want:       "red",
		},
		{
			name:       "at limit 100%",
			percentage: 100,
			want:       "red",
		},
		{
			name:       "over limit 110%",
			percentage: 110,
			want:       "red",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUsageLevel(tt.percentage)
			if got != tt.want {
				t.Errorf("GetUsageLevel(%v) = %v, want %v", tt.percentage, got, tt.want)
			}
		})
	}
}

// test getModelLimit
func TestGetModelLimit(t *testing.T) {
	tests := []struct {
		name  string
		model string
		want  int64
	}{
		{
			name:  "exact match claude-3-opus",
			model: "claude-3-opus",
			want:  200000,
		},
		{
			name:  "prefix match claude-3-opus-20240229",
			model: "claude-3-opus-20240229",
			want:  200000,
		},
		{
			name:  "claude-2 has 100k limit",
			model: "claude-2",
			want:  100000,
		},
		{
			name:  "unknown model defaults to 200k",
			model: "claude-future-model",
			want:  200000,
		},
		{
			name:  "case insensitive match",
			model: "CLAUDE-3-SONNET",
			want:  200000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getModelLimit(tt.model)
			if got != tt.want {
				t.Errorf("getModelLimit(%q) = %v, want %v", tt.model, got, tt.want)
			}
		})
	}
}

// benchmark calculation performance
func BenchmarkCalculate(b *testing.B) {
	usage := &parser.Usage{
		InputTokens:              9,
		CacheReadInputTokens:     58164,
		CacheCreationInputTokens: 1097,
		OutputTokens:             2,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Calculate(usage, "claude-sonnet-4-5")
	}
}

func BenchmarkGetUsageLevel(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetUsageLevel(75.5)
	}
}
