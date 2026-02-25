package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// trimLine removes \r and \n so they don't overwrite the line (e.g. index) in the terminal.
func trimLine(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

func display(
	configs []Config,
	aliasMaxLength *int,
	userMaxLength *int,
	identityFileMaxLength *int,
) {
	title := color.New(color.Bold, color.FgWhite).SprintFunc()
	yellow := color.New(color.Bold, color.FgHiGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	pink := color.New(color.FgHiMagenta).SprintFunc()

	fmt.Println(title("List of SSH services :"))
	fmt.Println()

	for i := range configs {
		index := strconv.Itoa(i)
		if i < 10 && len(configs) >= 10 {
			index = fmt.Sprintf("%s%d", " ", i)
		}

		fmt.Printf("%s %s %s %s %s \n",
			index,
			yellow(trimLine(addSpaceToEnd(configs[i].Alias, *aliasMaxLength+1))),
			red(trimLine(addSpaceToEnd(configs[i].User, *userMaxLength+1))),
			cyan(trimLine(addSpaceToEnd(configs[i].IdentityFile, *identityFileMaxLength+1))),
			pink(trimLine(configs[i].Hostname)))
	}

	fmt.Println("")
}
