/*
Copyright © 2023 AshutoshPatole
*/
package connect

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AshutoshPatole/ssh-manager/ssh"
	c "github.com/AshutoshPatole/ssh-manager/utils"
	"github.com/TwiN/go-color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cEnv string

// connectCmd represents the connect command
var ConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to the servers",
	Long: `
To connect to the servers use:
ssm connect group-name

You can also specify which environments to list:
ssm connect group-name -e ppd
	`,
	Aliases: []string{"c", "con"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 1 {
			fmt.Println(color.InYellow("Usage: ssm connect group-name\nYou can also pass environment using -e (optional)"))
			os.Exit(1)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ListToConnectServers(args[0], cEnv)
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	ConnectCmd.Flags().StringVarP(&cEnv, "env", "e", "", "Specify which environments servers should be displayed")
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type ServerOption struct {
	Label       string
	Environment string
	HostName    string
}

func ListToConnectServers(group, environment string) {
	var config c.Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	selectedEnvName := ""
	selectedHostName := ""

	serverOptions := []ServerOption{}
	user := ""

	for _, grp := range config.Groups {
		if grp.Name == group {
			user = grp.User
			for _, env := range grp.Environment {
				if environment != "" {
					if environment == env.Name {
						for _, server := range env.Servers {
							serverOption := ServerOption{
								Label:       fmt.Sprintf("%s (%s)", server.Alias, env.Name),
								Environment: env.Name,
								HostName:    server.HostName,
							}
							serverOptions = append(serverOptions, serverOption)
						}
					}
				} else {
					for _, server := range env.Servers {
						serverOption := ServerOption{
							Label:       fmt.Sprintf("%s (%s)", server.Alias, env.Name),
							Environment: env.Name,
							HostName:    server.HostName,
						}
						serverOptions = append(serverOptions, serverOption)
					}
				}
			}
		}
	}
	labels := make([]string, len(serverOptions))
	for i, serverOption := range serverOptions {
		labels[i] = serverOption.Label
	}

	fmt.Println(color.InGreen(selectedEnvName))
	prompt := &survey.Select{
		Message: "Select server",
		Options: labels,
	}
	survey.AskOne(prompt, &selectedHostName)

	// Extract environment name from the selected option
	for _, serverOption := range serverOptions {
		if serverOption.Label == selectedHostName {
			selectedEnvName = serverOption.Environment
			selectedHostName = strings.Split(serverOption.HostName, " (")[0]
			break
		}
	}
	if selectedHostName != "" && user != "" && selectedEnvName != "" {
		fmt.Println(color.InGreen("Host : " + selectedHostName))
		fmt.Println(color.InGreen("User : " + user))
		fmt.Println(color.InGreen("Environment : " + selectedEnvName))

		ssh.Connect(selectedHostName, user, selectedEnvName)
	} else {
		fmt.Println(color.InRed("Aborted!!"))
	}
}
