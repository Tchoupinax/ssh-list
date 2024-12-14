package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	sshConfigFile, err := os.ReadFile(fmt.Sprintf("%s%s", home, "/.ssh/config"))
	check(err)

	return string(sshConfigFile)
}

func processConfigsFromFile(
	sshContentFile string,
	aliasMaxLength *int,
	hostnameMaxLength *int,
	userMaxLength *int,
	identityFileMaxLength *int,
) []Config {
	fmt.Println("processConfigsFromFile")
	var configs []Config

	s := strings.Split(sshContentFile, "Host ")
	fmt.Println(len(s))

	for _, block := range s {
		if block == "" {
			continue
		}
		config := Config{}

		for _, line := range strings.Split(block, "\n") {
			if strings.Contains(line, "port") {
				config.Port, _ = strconv.ParseInt(strings.Trim(line, " "), 10, 64)
			}

			if strings.Contains(line, "User") {
				config.User = strings.TrimSpace(strings.Replace(line, "User", "", -1))
				if len(config.User) > *userMaxLength {
					*userMaxLength = len(config.User)
				}
			}

			if strings.Contains(line, "IdentityFile") {
				if !strings.Contains(line, "#") {
					config.IdentityFile = strings.TrimSpace(strings.Replace(line, "IdentityFile", "", -1))
				} else {
					config.IdentityFile = "❌"
				}

				if len(config.IdentityFile) > *identityFileMaxLength {
					*identityFileMaxLength = len(config.IdentityFile)
				}
			}

			if strings.Contains(line, "HostName") {
				config.Hostname = strings.TrimSpace(strings.Replace(line, "HostName", "", -1))
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
