package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
)

// trimLine removes \r and \n so they don't overwrite the line (e.g. index) in the terminal.
func trimLine(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

// padCell pads s to exactly width terminal columns (handles wide runes vs byte length).
func padCell(s string, width int) string {
	w := runewidth.StringWidth(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

type terminalTheme struct {
	title       *color.Color
	header      *color.Color
	separator   *color.Color
	idx         *color.Color
	alias       *color.Color
	user        *color.Color
	identity    *color.Color
	host        *color.Color
	port        *color.Color
	stripeBg    *color.Color
	stripeIdx   *color.Color
	stripeAlias *color.Color
	stripeUser  *color.Color
	stripeID    *color.Color
	stripeHost  *color.Color
	stripePort  *color.Color
}

func newTerminalTheme() terminalTheme {
	t := terminalTheme{
		title:     color.New(color.Bold, color.FgWhite),
		header:    color.New(color.Bold, color.FgHiWhite),
		separator: color.New(color.Bold, color.FgWhite),
		idx:       color.New(color.Bold, color.FgHiBlue),
		alias:     color.New(color.Bold, color.FgHiBlue),
		user:      color.New(color.Bold, color.FgHiBlue),
		identity:  color.New(color.Bold, color.FgHiBlue),
		host:      color.New(color.Bold, color.FgHiBlue),
		port:      color.New(color.Bold, color.Bold, color.FgHiBlue),
	}

	t.stripeBg = color.New(color.Bold, color.FgWhite)
	t.stripeIdx = color.New(color.Bold, color.FgWhite)
	t.stripeAlias = color.New(color.Bold, color.FgWhite)
	t.stripeUser = color.New(color.Bold, color.FgWhite)
	t.stripeID = color.New(color.Bold, color.FgWhite)
	t.stripeHost = color.New(color.Bold, color.FgWhite)
	t.stripePort = color.New(color.Bold, color.FgWhite)

	return t
}

func (t terminalTheme) printRow(
	index, aliasVal, userVal, idFileVal, hostCell string,
	striped bool,
	cpuVal, ramVal string,
	withStats bool,
) {
	sp := "  "
	if striped {
		spPad := t.stripeBg.Sprint(sp)
		if withStats {
			fmt.Printf("%s%s%s%s%s%s%s%s%s%s%s%s%s%s\n",
				spPad,
				t.stripeIdx.Sprint(index),
				spPad,
				t.stripeAlias.Sprint(aliasVal),
				spPad,
				t.stripeUser.Sprint(userVal),
				spPad,
				t.stripeID.Sprint(idFileVal),
				spPad,
				t.stripeHost.Sprint(hostCell),
				spPad,
				t.stripeHost.Sprint(cpuVal),
				spPad,
				t.stripeHost.Sprint(ramVal),
			)
			return
		}
		fmt.Printf("%s%s%s%s%s%s%s%s%s%s\n",
			spPad,
			t.stripeIdx.Sprint(index),
			spPad,
			t.stripeAlias.Sprint(aliasVal),
			spPad,
			t.stripeUser.Sprint(userVal),
			spPad,
			t.stripeID.Sprint(idFileVal),
			spPad,
			t.stripeHost.Sprint(hostCell),
		)
		return
	}

	if withStats {
		fmt.Printf("  %s  %s  %s  %s  %s  %s  %s\n",
			t.idx.Sprint(index),
			t.alias.Sprint(aliasVal),
			t.user.Sprint(userVal),
			t.identity.Sprint(idFileVal),
			t.host.Sprint(hostCell),
			t.host.Sprint(cpuVal),
			t.host.Sprint(ramVal),
		)
		return
	}

	fmt.Printf("  %s  %s  %s  %s  %s\n",
		t.idx.Sprint(index),
		t.alias.Sprint(aliasVal),
		t.user.Sprint(userVal),
		t.identity.Sprint(idFileVal),
		t.host.Sprint(hostCell),
	)
}

// columnWidth is max content width + 1 for a trailing space (same idea as previous col()).
func columnWidth(maxContent int) int {
	if maxContent < 0 {
		return 0
	}
	return maxContent + 1
}

func display(
	configs []Config,
	aliasMaxLength *int,
	userMaxLength *int,
	identityFileMaxLength *int,
	hostnameMaxLength *int,
	stats []ServerStats,
) {
	th := newTerminalTheme()

	// Recompute max display widths so columns align (byte len from config parse is wrong for UTF-8 / em dash).
	aliasW := *aliasMaxLength
	userW := *userMaxLength
	idFileW := *identityFileMaxLength
	hostW := *hostnameMaxLength

	maxInt := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	for _, c := range configs {
		if w := runewidth.StringWidth(trimLine(c.Alias)); w > aliasW {
			aliasW = w
		}
		if w := runewidth.StringWidth(trimLine(c.User)); w > userW {
			userW = w
		}
		if w := runewidth.StringWidth(trimLine(c.IdentityFile)); w > idFileW {
			idFileW = w
		}
		hp := trimLine(c.Hostname)
		if c.Port != 0 && c.Port != 22 {
			hp += fmt.Sprintf(":%d", c.Port)
		}
		if w := runewidth.StringWidth(hp); w > hostW {
			hostW = w
		}
	}

	aliasW = maxInt(aliasW, runewidth.StringWidth("Alias"))
	userW = maxInt(userW, runewidth.StringWidth("User"))
	idFileW = maxInt(idFileW, runewidth.StringWidth("Identity file"))
	hostW = maxInt(hostW, runewidth.StringWidth("Host"))

	withStats := stats != nil && len(stats) == len(configs)
	cpuW, ramW := runewidth.StringWidth("CPU"), runewidth.StringWidth("RAM")
	if withStats {
		for i := range configs {
			if w := runewidth.StringWidth(trimLine(stats[i].CPUString())); w > cpuW {
				cpuW = w
			}
			if w := runewidth.StringWidth(trimLine(stats[i].RAMString())); w > ramW {
				ramW = w
			}
		}
	}

	aw := columnWidth(aliasW)
	uw := columnWidth(userW)
	iw := columnWidth(idFileW)
	hw := columnWidth(hostW)
	cw := columnWidth(cpuW)
	rw := columnWidth(ramW)

	indexWidth := len(strconv.Itoa(len(configs) - 1))
	if indexWidth < 1 {
		indexWidth = 1
	}

	gap := 2 // two spaces between columns, matching fmt "  %s  %s"
	separatorLen := gap + indexWidth + gap + aw + gap + uw + gap + iw + gap + hw
	if withStats {
		separatorLen += gap + cw + gap + rw
	}

	fmt.Println()
	fmt.Printf("  %s\n", th.title.Sprint("SSH connections"))
	fmt.Println()

	indexPad := strings.Repeat(" ", indexWidth)
	if withStats {
		fmt.Printf("  %s  %s  %s  %s  %s  %s  %s\n",
			th.separator.Sprint(indexPad),
			th.header.Sprint(padCell("Alias", aw)),
			th.header.Sprint(padCell("User", uw)),
			th.header.Sprint(padCell("Identity file", iw)),
			th.header.Sprint(padCell("Host", hw)),
			th.header.Sprint(padCell("CPU", cw)),
			th.header.Sprint(padCell("RAM", rw)),
		)
	} else {
		fmt.Printf("  %s  %s  %s  %s  %s\n",
			th.separator.Sprint(indexPad),
			th.header.Sprint(padCell("Alias", aw)),
			th.header.Sprint(padCell("User", uw)),
			th.header.Sprint(padCell("Identity file", iw)),
			th.header.Sprint(padCell("Host", hw)),
		)
	}

	fmt.Printf("  %s\n", th.separator.Sprint(strings.Repeat("─", separatorLen)))

	for i := range configs {
		index := strconv.Itoa(i)
		index = strings.Repeat(" ", indexWidth-len(index)) + index

		aliasVal := padCell(trimLine(configs[i].Alias), aw)
		userVal := padCell(trimLine(configs[i].User), uw)
		idFileVal := padCell(trimLine(configs[i].IdentityFile), iw)

		hostPort := trimLine(configs[i].Hostname)
		if configs[i].Port != 0 && configs[i].Port != 22 {
			hostPort += fmt.Sprintf(":%d", configs[i].Port)
		}
		hostCell := padCell(hostPort, hw)

		cpuVal := ""
		ramVal := ""
		if withStats {
			cpuVal = padCell(trimLine(stats[i].CPUString()), cw)
			ramVal = padCell(trimLine(stats[i].RAMString()), rw)
		}

		th.printRow(index, aliasVal, userVal, idFileVal, hostCell, i%2 == 1, cpuVal, ramVal, withStats)
	}

	fmt.Println()
}
