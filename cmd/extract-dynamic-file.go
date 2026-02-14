package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
)

func extractDynamicFile(
	sshContentFile string,
	aliasMaxLength *int,
	hostnameMaxLength *int,
	userMaxLength *int,
	identityFileMaxLength *int,
) []Config {
	home, err := homedir.Dir()
	check(err)

	var configs []Config

	re := regexp.MustCompile(`Include (.*)`)
	for _, line := range re.FindAllString(sshContentFile, -1) {
		folderPath := strings.NewReplacer(
			"Include ", "",
			"/*", "",
			"~", home,
		).Replace(line)

		entries, err := os.ReadDir(folderPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, entry := range entries {
			sshConfigFile, err := os.ReadFile(folderPath + "/" + entry.Name())
			if err == nil && sshConfigFile != nil {
				configs = append(configs, processConfigsFromFile(string(sshConfigFile), aliasMaxLength, hostnameMaxLength, userMaxLength, identityFileMaxLength)...)
			}
		}
	}

	return configs
}

func extractSSHConfigFile() string {
	home, err := homedir.Dir()
	check(err)

	path, err := homedir.Expand(cfgFile)
	check(err)

	sshConfigFile, err := os.ReadFile(path)
	if len(sshConfigFile) == 0 {
		sshConfigFile, err = os.ReadFile(fmt.Sprintf("%s%s", home, "/.ssh/config"))
	}

	if err != nil {
		title := color.New(color.Bold, color.FgHiRed).SprintFunc()
		fmt.Printf(title("Configuration %s does not exist\n"), path)
		os.Exit(1)
	}

	return string(sshConfigFile)
}

func processConfigsFromFile(
	sshContentFile string,
	aliasMaxLength *int,
	hostnameMaxLength *int,
	userMaxLength *int,
	identityFileMaxLength *int,
) []Config {
	// nolint: prealloc // Dynamic size is unknown
	var configs []Config

	s := strings.Split(sshContentFile, "Host ")

	for _, block := range s {
		if block == "" {
			continue
		}
		if !strings.Contains(block, "Host") {
			continue
		}

		config := Config{}

		for _, line := range strings.Split(block, "\n") {
			if strings.Contains(line, "port") {
				config.Port, _ = strconv.ParseInt(strings.Trim(line, " "), 10, 64)
			}

			if strings.Contains(line, "User") {
				config.User = strings.TrimSpace(strings.ReplaceAll(line, "User", ""))
				if len(config.User) > *userMaxLength {
					*userMaxLength = len(config.User)
				}
			}

			if strings.Contains(line, "IdentityFile") {
				if !strings.Contains(line, "#") {
					config.IdentityFile = strings.TrimSpace(strings.ReplaceAll(line, "IdentityFile", ""))
				} else {
					config.IdentityFile = "❌"
				}

				if len(config.IdentityFile) > *identityFileMaxLength {
					*identityFileMaxLength = len(config.IdentityFile)
				}
			}

			if strings.Contains(line, "HostName") {
				config.Hostname = strings.TrimSpace(strings.ReplaceAll(line, "HostName", ""))
				if len(config.Hostname) > *hostnameMaxLength {
					*hostnameMaxLength = len(config.Hostname)
				}
			}
		}

		config.Alias = strings.Split(block, "\n")[0]
		if len(config.Alias) > *aliasMaxLength && !strings.HasPrefix(config.Alias, "#") {
			*aliasMaxLength = len(config.Alias)
		}

		configs = append(configs, config)
	}

	return configs
}
