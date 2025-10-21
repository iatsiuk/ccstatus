package formatter

import (
	"ccstatus/internal/calculator"
	"fmt"
	"os"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorCyan   = "\033[36m"
)

// uses os.ModeCharDevice to detect TTY on unix systems (macOS, Linux)
func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// automatically detects TTY and falls back to plain output
func Format(info calculator.ContextInfo, model string) string {
	if !isTerminal(os.Stdout) {
		return FormatPlain(info, model)
	}
	return formatWithColors(info, model)
}

// used internally by Format() and for testing
func formatWithColors(info calculator.ContextInfo, model string) string {
	usageLevel := calculator.GetUsageLevel(info.Percentage)
	color := getColor(usageLevel)

	// format: [ctx: 59261/200000 29.6%] model
	return fmt.Sprintf("%s[ctx: %d/%d %.1f%%]%s %s",
		color,
		info.CurrentTokens,
		info.MaxTokens,
		info.Percentage,
		ColorReset,
		ColorCyan+model+ColorReset,
	)
}

func FormatPlain(info calculator.ContextInfo, model string) string {
	return fmt.Sprintf("[ctx: %d/%d %.1f%%] %s",
		info.CurrentTokens,
		info.MaxTokens,
		info.Percentage,
		model,
	)
}

// automatically detects TTY and falls back to plain output
func FormatError(errorMsg string) string {
	if !isTerminal(os.Stdout) {
		return fmt.Sprintf("[ERROR: %s]", errorMsg)
	}
	return fmt.Sprintf("%s[ERROR: %s]%s",
		ColorRed,
		errorMsg,
		ColorReset,
	)
}

func getColor(level string) string {
	switch level {
	case "green":
		return ColorGreen
	case "yellow":
		return ColorYellow
	case "red":
		return ColorRed
	default:
		return ColorReset
	}
}
