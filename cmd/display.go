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

type terminalTheme struct {
	title      *color.Color
	header     *color.Color
	separator  *color.Color
	idx        *color.Color
	alias      *color.Color
	user       *color.Color
	identity   *color.Color
	host       *color.Color
	port       *color.Color
	stripeBg   *color.Color
	stripeIdx  *color.Color
	stripeAlias *color.Color
	stripeUser *color.Color
	stripeID   *color.Color
	stripeHost *color.Color
	stripePort *color.Color
}

func newTerminalTheme() terminalTheme {
	t := terminalTheme{
		title:    color.New(color.Bold, color.FgWhite),
		header:   color.New(color.Bold, color.FgHiWhite),
		separator: color.New(color.Bold, color.FgWhite),
		idx: color.New(color.Bold, color.FgHiBlue),
		alias: color.New(color.Bold, color.FgHiBlue),
		user: color.New(color.Bold, color.FgHiBlue),
		identity: color.New(color.Bold, color.FgHiBlue),
		host: color.New(color.Bold, color.FgHiBlue),
		port: color.New(color.Bold, color.Bold, color.FgHiBlue),
	}

	// --- Zebra stripe: subtle slate panel + light text (readable on the panel) ---
	const sr, sg, sb = 45, 48, 55
	t.stripeBg = color.New(color.Bold, color.FgWhite)
	t.stripeIdx =color.New(color.Bold, color.FgWhite)
	t.stripeAlias = color.New(color.Bold, color.FgWhite)
	t.stripeUser = color.New(color.Bold, color.FgWhite)
	t.stripeID = color.New(color.Bold, color.FgWhite)
	t.stripeHost = color.New(color.Bold, color.FgWhite)
	t.stripePort = color.New(color.Bold, color.FgWhite)

	return t
}

func (t terminalTheme) printRow(
	index, aliasVal, userVal, idFileVal, hostnameVal, portVal string,
	striped bool,
) {
	sp := "  "
	if striped {
		spPad := t.stripeBg.Sprint(sp)
		fmt.Printf("%s%s%s%s%s%s%s%s%s%s%s\n",
			spPad,
			t.stripeIdx.Sprint(index),
			spPad,
			t.stripeAlias.Sprint(aliasVal),
			spPad,
			t.stripeUser.Sprint(userVal),
			spPad,
			t.stripeID.Sprint(idFileVal),
			spPad,
			t.stripeHost.Sprint(hostnameVal),
			t.stripePort.Sprint(portVal),
		)
		return
	}

	fmt.Printf("  %s  %s  %s  %s  %s%s\n",
		t.idx.Sprint(index),
		t.alias.Sprint(aliasVal),
		t.user.Sprint(userVal),
		t.identity.Sprint(idFileVal),
		t.host.Sprint(hostnameVal),
		t.port.Sprint(portVal),
	)
}

func display(
	configs []Config,
	aliasMaxLength *int,
	userMaxLength *int,
	identityFileMaxLength *int,
	hostnameMaxLength *int,
) {
	th := newTerminalTheme()

	aliasW := *aliasMaxLength
	userW := *userMaxLength
	idFileW := *identityFileMaxLength
	hostW := *hostnameMaxLength

	// +1 keeps column padding aligned with cell content below
	col := func(w int) int { return w + 1 }

	indexWidth := len(strconv.Itoa(len(configs) - 1))
	if indexWidth < 1 {
		indexWidth = 1
	}

	fmt.Println()
	fmt.Printf("  %s\n", th.title.Sprint("SSH connections"))
	fmt.Println()

	indexPad := strings.Repeat(" ", indexWidth)
	fmt.Printf("  %s  %s  %s  %s  %s\n",
		th.separator.Sprint(indexPad),
		th.header.Sprint(addSpaceToEnd("Alias", col(aliasW))),
		th.header.Sprint(addSpaceToEnd("User", col(userW))),
		th.header.Sprint(addSpaceToEnd("Identity file", col(idFileW))),
		th.header.Sprint(addSpaceToEnd("Host", col(hostW))),
	)

	separatorLen := indexWidth + 2 + col(aliasW) + 3 + col(userW) + 3 + col(idFileW) + 3 + col(hostW) + 3
	fmt.Printf("  %s\n", th.separator.Sprint(strings.Repeat("─", separatorLen)))

	for i := range configs {
		index := strconv.Itoa(i)
		index = strings.Repeat(" ", indexWidth-len(index)) + index

		aliasVal := trimLine(addSpaceToEnd(configs[i].Alias, col(aliasW)))
		userVal := trimLine(addSpaceToEnd(configs[i].User, col(userW)))
		idFileVal := trimLine(addSpaceToEnd(configs[i].IdentityFile, col(idFileW)))
		hostnameVal := trimLine(configs[i].Hostname)
		portVal := ""
		if configs[i].Port != 0 && configs[i].Port != 22 {
			portVal = fmt.Sprintf(":%d", configs[i].Port)
		}

		th.printRow(index, aliasVal, userVal, idFileVal, hostnameVal, portVal, i%2 == 1)
	}

	fmt.Println()
}
