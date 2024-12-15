package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const Version string = "0.3.1"
const BuildDate string = "2024-12-15"

func init() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of ssh-list",
		Long:  "Print the version number of ssh-list",
		Run: func(cmd *cobra.Command, args []string) {
			bold := color.New(color.Bold).SprintFunc()
			italic := color.New(color.Italic).SprintFunc()

			fmt.Println()
			fmt.Println(bold("⚡️ SSH List"))
			fmt.Println()
			fmt.Println("build date: ", bold(BuildDate))
			fmt.Println("version:         ", bold(Version))
			fmt.Println()
			fmt.Println(italic("Need help?"))
			fmt.Println(italic("https://github.com/Tchoupinax/ssh-list/issues"))
			os.Exit(0)
		},
	}

	RootCmd.AddCommand(versionCmd)
}
