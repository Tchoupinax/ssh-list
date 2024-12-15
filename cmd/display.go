package cmd

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
)

func display(
	configs []Config,
	aliasMaxLength *int,
	hostnameMaxLength *int,
	userMaxLength *int,
	identityFileMaxLength *int,
) {
	title := color.New(color.Bold, color.FgWhite).SprintFunc()
	fmt.Println(title("List of SSH services :"))
	fmt.Println()

	for i := 0; i < len(configs); i++ {
		yellow := color.New(color.Bold, color.FgHiGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		pink := color.New(color.FgHiMagenta).SprintFunc()

		index := strconv.Itoa(i)
		if i < 10 {
			index = fmt.Sprintf("%s%d", " ", i)
		}

		fmt.Printf("%s %s %s %s %s \n",
			index,
			yellow(addSpaceToEnd(configs[i].Alias, *aliasMaxLength+1)),
			red(addSpaceToEnd(configs[i].User, *userMaxLength+1)),
			cyan(addSpaceToEnd(configs[i].IdentityFile, *identityFileMaxLength+1)),
			pink(configs[i].Hostname))
	}

	fmt.Println("")
}
