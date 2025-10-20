package formatter

import (
	"ccstatus/internal/calculator"
	"fmt"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
	ColorCyan   = "\033[36m"
)

// Format creates colored status line output
// Returns formatted string with ANSI colors based on usage level
func Format(info calculator.ContextInfo, model string) string {
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

// FormatCompact creates shorter status line for limited space
func FormatCompact(info calculator.ContextInfo) string {
	usageLevel := calculator.GetUsageLevel(info.Percentage)
	color := getColor(usageLevel)

	// format: [29.6%]
	return fmt.Sprintf("%s[%.1f%%]%s",
		color,
		info.Percentage,
		ColorReset,
	)
}

// FormatPlain creates output without colors
func FormatPlain(info calculator.ContextInfo, model string) string {
	return fmt.Sprintf("[ctx: %d/%d %.1f%%] %s",
		info.CurrentTokens,
		info.MaxTokens,
		info.Percentage,
		model,
	)
}

// FormatError creates error output with red color
func FormatError(errorMsg string) string {
	return fmt.Sprintf("%s[ERROR: %s]%s",
		ColorRed,
		errorMsg,
		ColorReset,
	)
}

// getColor returns ANSI color code based on usage level
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
