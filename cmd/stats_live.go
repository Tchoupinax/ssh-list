package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"
)

// displayStatsLive prints the table immediately with "…" in CPU/RAM, then refreshes in the
// alternate screen buffer as each host responds. Falls back to blocking fetch when stdout is not a TTY.
func displayStatsLive(
	configs []Config,
	aliasMaxLength *int,
	userMaxLength *int,
	identityFileMaxLength *int,
	hostnameMaxLength *int,
) {
	if len(configs) == 0 {
		display(configs, aliasMaxLength, userMaxLength, identityFileMaxLength, hostnameMaxLength, nil)
		return
	}
	if os.Getenv("TERM") == "dumb" || !term.IsTerminal(int(os.Stdout.Fd())) {
		stats := fetchAllServerStats(configs)
		display(configs, aliasMaxLength, userMaxLength, identityFileMaxLength, hostnameMaxLength, stats)
		return
	}

	started := time.Now()

	stats := make([]ServerStats, len(configs))
	var mu sync.Mutex

	const maxConcurrent = 8
	sem := make(chan struct{}, maxConcurrent)
	doneCh := make(chan struct{}, len(configs))

	for i := range configs {
		i := i
		go func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			s := fetchServerStats(configs[i])
			mu.Lock()
			stats[i] = s
			mu.Unlock()
			doneCh <- struct{}{}
		}()
	}

	leaveAlternate := func() {
		fmt.Print("\033[?1049l\033[?25h")
	}
	// If we panic while in the alternate buffer, restore the main screen.
	leftAlternate := false
	defer func() {
		if !leftAlternate {
			leaveAlternate()
		}
	}()

	fmt.Print("\033[?1049h\033[H\033[?25l")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		leaveAlternate()
		fmt.Fprintln(os.Stderr, "^C")
		os.Exit(130)
	}()

	spinFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧"}
	received := 0
	spin := 0

	dim := terminalMutedColor().SprintFunc()

	redraw := func() {
		mu.Lock()
		rec := received
		sp := spin
		snap := make([]ServerStats, len(stats))
		copy(snap, stats)
		mu.Unlock()

		body := renderTableToString(
			configs,
			aliasMaxLength,
			userMaxLength,
			identityFileMaxLength,
			hostnameMaxLength,
			snap,
		)
		var footer string
		if rec >= len(configs) {
			footer = fmt.Sprintf("  %s\n", dim(fmt.Sprintf("Metrics %d/%d · done", len(configs), len(configs))))
		} else {
			footer = fmt.Sprintf("  %s\n", dim(fmt.Sprintf("Fetching metrics %d/%d %s", rec, len(configs), spinFrames[sp%len(spinFrames)])))
		}
		fmt.Print("\033[H")
		fmt.Print(body)
		fmt.Print(footer)
		fmt.Print("\033[J")
	}

	ticker := time.NewTicker(180 * time.Millisecond)
	defer ticker.Stop()

	redraw()

	for received < len(configs) {
		select {
		case <-doneCh:
			mu.Lock()
			received++
			mu.Unlock()
		case <-ticker.C:
			mu.Lock()
			spin++
			mu.Unlock()
		}
		redraw()
	}

	// Leaving the alternate screen discards its contents — print the final table on the main buffer
	// so it stays in the scrollback. Timing is appended without clearing the table.
	mu.Lock()
	snap := make([]ServerStats, len(stats))
	copy(snap, stats)
	mu.Unlock()

	leftAlternate = true
	leaveAlternate()

	body := renderTableToString(
		configs,
		aliasMaxLength,
		userMaxLength,
		identityFileMaxLength,
		hostnameMaxLength,
		snap,
	)
	fmt.Print(body)
	elapsed := time.Since(started)
	fmt.Printf("  %s\n\n", dim(fmt.Sprintf("Metrics fetched in %v", elapsed.Round(time.Millisecond))))
}
