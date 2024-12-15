package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func measurePerformance(sshConnectionCreated *bool) func() {
	italic := color.New(color.Italic).SprintFunc()
	start := time.Now()

	return func() {
		if !*sshConnectionCreated {
			fmt.Printf(italic("performed in %v\n"), time.Since(start))
		}
	}
}
