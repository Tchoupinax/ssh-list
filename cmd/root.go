package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
		home, err := homedir.Dir()
		check(err)
		dat, err := os.ReadFile(fmt.Sprintf("%s%s", home, "/.ssh/config"))
		check(err)
		s := strings.Split(string(dat), "Host ")

		var configs []Config
		var aliasMaxLength = 0
		var hostnameMaxLength = 0
		var userMaxLength = 0

		title := color.New(color.Bold, color.FgWhite).SprintFunc()
		fmt.Println(title("List of SSH services22 :"))
		fmt.Println()

		for _, block := range s {
			config := Config{}

			for _, line := range strings.Split(block, "\n") {
				if strings.Contains(line, "port") {
					config.Port, _ = strconv.ParseInt(strings.Trim(line, " "), 10, 64)
				}

				if strings.Contains(line, "User") {
					config.User = strings.TrimSpace(strings.Replace(line, "User", "", -1))
					if len(config.User) > userMaxLength {
						userMaxLength = len(config.User)
					}
				}

				if strings.Contains(line, "IdentityFile") {
					if !strings.Contains(line, "#") {
						config.IdentityFile = strings.TrimSpace(strings.Replace(line, "IdentityFile", "", -1))
					} else {
						config.IdentityFile = "❌"
					}
				}

				if strings.Contains(line, "HostName") {
					config.Hostname = strings.TrimSpace(strings.Replace(line, "HostName", "", -1))
					if len(config.Hostname) > hostnameMaxLength {
						hostnameMaxLength = len(config.Hostname)
					}
				}
			}

			config.Alias = strings.Split(block, "\n")[0]
			if len(config.Alias) > aliasMaxLength && !strings.HasPrefix(config.Alias, "#") {
				aliasMaxLength = len(config.Alias)
			}

			configs = append(configs, config)
		}

		for i := 0; i < len(configs); i++ {
			if i == 0 {
				continue
			}

			yellow := color.New(color.Bold, color.FgHiGreen).SprintFunc()
			red := color.New(color.FgRed).SprintFunc()
			cyan := color.New(color.FgCyan).SprintFunc()

			index := strconv.Itoa(i)
			if i < 10 {
				index = fmt.Sprintf("%s%d", " ", i)
			}

			fmt.Printf("%s %s %s %s \n",
				index,
				yellow(addSpaceToEnd(configs[i].Alias, aliasMaxLength+1)),
				red(addSpaceToEnd(configs[i].User, userMaxLength+1)),
				cyan(configs[i].IdentityFile))
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
