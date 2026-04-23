package cmd

import (
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

// terminalUsesLightBackground reports whether stdout looks like a light-background terminal
// (OSC 11, COLORFGBG, or similar). When false, use the dark-terminal palette.
// The result is cached for the process (OSC query runs at most once).
func terminalUsesLightBackground() bool {
	return terminalLightBackgroundCached()
}

var terminalLightBackgroundCached = sync.OnceValue(func() bool {
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}
	out := termenv.NewOutput(os.Stdout)
	return !out.HasDarkBackground()
})

// terminalMutedColor returns a faint color suitable for secondary text on the current background.
func terminalMutedColor() *color.Color {
	if terminalUsesLightBackground() {
		return color.New(color.Faint, color.FgHiBlack)
	}
	return color.New(color.Faint, color.FgHiWhite)
}
