package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func addSpaceToEnd(s string, size int) string {
	var diff = size - len(s)

	for i := 0; i < diff; i++ {
		s = fmt.Sprintf("%s%s", s, " ")
	}

	return s
}

type Config struct {
	Alias        string
	Hostname     string
	User         string
	Port         int64
	IdentityFile string
}

var RootCmd = &cobra.Command{
	Use:   "ssh-list",
	Short: "An example of cobra",
	Long:  `List your SSH configurations easily`,
	Run: func(cmd *cobra.Command, args []string) {
		var configs []Config
		var aliasMaxLength = 0
		var hostnameMaxLength = 0
		var userMaxLength = 0
		var identityFileMaxLength = 0

		var content = extractSSHConfigFile()

		configs = append(
			configs,
			processConfigsFromFile(
				content,
				&aliasMaxLength,
				&hostnameMaxLength,
				&userMaxLength,
				&identityFileMaxLength,
			)...,
		)

		configs = append(
			configs,
			extractDynamicFile(
				content,
				&aliasMaxLength,
				&hostnameMaxLength,
				&userMaxLength,
				&identityFileMaxLength,
			)...,
		)

		sort.Slice(configs, func(i, j int) bool {
			return configs[i].Alias < configs[j].Alias
		})

		title := color.New(color.Bold, color.FgWhite).SprintFunc()
		fmt.Println(title("List of SSH services :"))
		fmt.Println()

		for i := 0; i < len(configs); i++ {
			if i == 0 {
				continue
			}

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
				yellow(addSpaceToEnd(configs[i].Alias, aliasMaxLength+1)),
				red(addSpaceToEnd(configs[i].User, userMaxLength+1)),
				cyan(addSpaceToEnd(configs[i].IdentityFile, identityFileMaxLength+1)),
				pink(configs[i].Hostname))
		}

		fmt.Println("")

		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "~/.ssh/config", "SSH config file")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "version", "", "display version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra-example" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra-example")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
