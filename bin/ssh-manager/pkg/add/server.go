/*
Copyright © 2023 AshutoshPatole
*/
package add

import (
	"fmt"
	"os"

	"github.com/AshutoshPatole/ssh-manager/ssh"
	"github.com/TwiN/go-color"
	"github.com/spf13/cobra"
)

var aUser, aGroup, aEnv, aAlias string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use: "server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf(color.InRedOverBlack("Expects server name or ip address"))
		}
		return nil
	},
	Short: "A brief description of your command",
	Long: `
To create a server:
ssm add server host-name -u user -g group -e dev
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		isValid := false
		for _, env := range allowedEnvironments {
			if aEnv == env {
				isValid = true
			}
		}
		if isValid {
			addServer(args[0], aUser)
		} else {
			fmt.Print(color.InYellow("Unknown environment: Allowed values are "))
			fmt.Println(allowedEnvironments)
			os.Exit(1)
		}
	},
}

func init() {
	AddCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&aGroup, "group", "g", "", "Group name in which this server should be added (required)")
	serverCmd.MarkFlagRequired("group")
	serverCmd.Flags().StringVarP(&aUser, "user", "u", "", "User name to connect (required)")
	serverCmd.MarkFlagRequired("user")
	serverCmd.Flags().StringVarP(&aEnv, "env", "e", "dev", "Enviornment name to store this server. Allowed values are [prd, ppd, uat, sit, dev]")
	serverCmd.Flags().StringVarP(&aAlias, "alias", "a", "", "Alias for the server")
	serverCmd.MarkFlagRequired("alias")

}

func addServer(server, user string) {
	password := ssh.AskPass()
	ssh.InitServer(server, user, password, aGroup, aEnv, aAlias)
}
