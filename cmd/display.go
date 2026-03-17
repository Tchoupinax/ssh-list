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
	hostnameMaxLength *int,
) {
	header := color.New(color.Bold, color.FgHiWhite).SprintFunc()
	dim := color.New(color.FgHiBlack).SprintFunc()
	styleAlias := color.New(color.Bold, color.FgHiCyan).SprintFunc()
	styleUser := color.New(color.FgHiGreen).SprintFunc()
	styleIdentityFile := color.New(color.FgHiYellow).SprintFunc()
	styleHostname := color.New(color.FgHiWhite).SprintFunc()
	stylePort := color.New(color.FgHiBlack).SprintFunc()

	aliasW := *aliasMaxLength
	userW := *userMaxLength
	idFileW := *identityFileMaxLength
	hostW := *hostnameMaxLength

	indexWidth := len(strconv.Itoa(len(configs) - 1))
	if indexWidth < 1 {
		indexWidth = 1
	}

	fmt.Println()
	fmt.Printf("  %s\n", header("SSH Connections"))
	fmt.Println()

	indexPad := strings.Repeat(" ", indexWidth)
	fmt.Printf("  %s  %s  %s  %s  %s\n",
		dim(indexPad),
		header(addSpaceToEnd("Alias", aliasW+1)),
		header(addSpaceToEnd("User", userW+1)),
		header(addSpaceToEnd("Identity File", idFileW+1)),
		header(addSpaceToEnd("Host", hostW+1)),
	)

	separatorLen := indexWidth + 2 + aliasW + 3 + userW + 3 + idFileW + 3 + hostW + 3
	fmt.Printf("  %s\n", dim(strings.Repeat("─", separatorLen)))

	for i := range configs {
		index := strconv.Itoa(i)
		index = strings.Repeat(" ", indexWidth-len(index)) + index

		aliasVal := trimLine(addSpaceToEnd(configs[i].Alias, aliasW+1))
		userVal := trimLine(addSpaceToEnd(configs[i].User, userW+1))
		idFileVal := trimLine(addSpaceToEnd(configs[i].IdentityFile, idFileW+1))
		hostnameVal := trimLine(configs[i].Hostname)
		portVal := ""
		if configs[i].Port != 0 && configs[i].Port != 22 {
			portVal = fmt.Sprintf(":%d", configs[i].Port)
		}

		if i%2 == 0 {
			fmt.Printf("  %s  %s  %s  %s  %s%s\n",
				dim(index),
				styleAlias(aliasVal),
				styleUser(userVal),
				styleIdentityFile(idFileVal),
				styleHostname(hostnameVal),
				stylePort(portVal),
			)
		} else {
			bg := color.BgHiBlack
			sp := color.New(bg).Sprint("  ")
			fmt.Printf("%s%s%s%s%s%s%s%s%s%s%s\n",
				sp,
				color.New(color.FgWhite, bg).Sprint(index),
				sp,
				color.New(color.Bold, color.FgHiCyan, bg).Sprint(aliasVal),
				sp,
				color.New(color.FgHiGreen, bg).Sprint(userVal),
				sp,
				color.New(color.FgHiYellow, bg).Sprint(idFileVal),
				sp,
				color.New(color.FgHiWhite, bg).Sprint(hostnameVal),
				color.New(color.FgWhite, bg).Sprint(portVal),
			)
		}
	}

	fmt.Println()
}
