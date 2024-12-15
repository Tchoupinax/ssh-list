package cmd

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

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

var sshConnectionCreated bool

var RootCmd = &cobra.Command{
	Use:   "ssh-list",
	Short: "An example of cobra",
	Long:  `List your SSH configurations easily`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		sshConnectionCreated = false
		defer measurePerformance(&sshConnectionCreated)()

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

		if len(args) > 0 {
			firstArg := args[0]
			if firstArg != "" {
				// Try to detect if it's an integer
				index, err := strconv.Atoi(firstArg)
				if err != nil {
					// Name Case
					configs = filterByRegex(configs, firstArg)

				} else {
					sshConnectionCreated = true
					createSSH(configs[index])
				}
			}
		}

		// If the length is 1, auto activate
		if len(configs) == 1 {
			sshConnectionCreated = true
			createSSH(configs[0])
		}

		if !sshConnectionCreated {
			display(
				configs,
				&aliasMaxLength,
				&userMaxLength,
				&identityFileMaxLength,
			)
		}
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

func filterByRegex(configs []Config, name string) []Config {
	// Compile the regex pattern
	re, err := regexp.Compile(".*" + name + ".*")
	if err != nil {
		return []Config{}
	}

	var filtered []Config
	for _, config := range configs {
		for range re.FindAllString(config.Alias, -1) {
			filtered = append(filtered, config)
		}
	}
	return filtered
}
