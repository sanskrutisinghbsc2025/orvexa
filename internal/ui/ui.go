package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// PrintHeader prints a styled header.
func PrintHeader(title string) {
	fmt.Println()
	color.New(color.Bold, color.FgCyan).Print("━━━ " + strings.ToUpper(title) + " ")
	lineLen := 40 - len(title)
	if lineLen < 5 {
		lineLen = 5
	}
	color.New(color.FgHiBlack).Println(strings.Repeat("━", lineLen))
}

// PrintSuccess prints a success message with an icon.
func PrintSuccess(msg string) {
	color.New(color.FgGreen).Printf("✔ %s\n", msg)
}

// PrintInfo prints an info message with an icon.
func PrintInfo(msg string) {
	color.New(color.FgBlue).Printf("ℹ %s\n", msg)
}

// PrintWarning prints a warning message with an icon.
func PrintWarning(msg string) {
	color.New(color.FgYellow).Printf("⚠ %s\n", msg)
}

// PrintError prints an error message with an icon.
func PrintError(msg string) {
	color.New(color.FgRed, color.Bold).Printf("✘ ERROR: %s\n", msg)
}

// PrintKV prints a key-value pair with styling.
func PrintKV(k, v string) {
	color.New(color.FgHiBlack).Printf("  %-15s ", k+":")
	fmt.Println(v)
}

// PrintBranch prints a branch reference with current indicator.
func PrintBranch(name string, active bool) {
	if active {
		color.New(color.FgGreen, color.Bold).Printf("* %s (active)\n", name)
	} else {
		fmt.Printf("  %s\n", name)
	}
}
