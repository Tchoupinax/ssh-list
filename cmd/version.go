package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version   string
	buildDate string
	commit    string
)

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
			fmt.Println("build date: ", bold(version))
			fmt.Println("version:    ", bold(buildDate))
			fmt.Println("commit:     ", bold(commit))
			fmt.Println()
			fmt.Println(italic("Need help?"))
			fmt.Println(italic("https://github.com/Tchoupinax/ssh-list/issues"))
			os.Exit(0)
		},
	}

	RootCmd.AddCommand(versionCmd)
}
