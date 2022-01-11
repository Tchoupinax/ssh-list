package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	helloCmd := &cobra.Command{
		Use:   "hello",
		Short: "says hello",
		Long:  `This subcommand says hello`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello called")
		},
	}

	RootCmd.AddCommand(helloCmd)
}
