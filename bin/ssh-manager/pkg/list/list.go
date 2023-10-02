/*
Copyright © 2023 AshutoshPatole
*/
package list

import (
	"fmt"
	"log"
	"os"

	c "github.com/AshutoshPatole/ssh-manager/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List groups and servers",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("list called")
	// },
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ListGroups() {
	var config c.Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID.", "Group Name", "Environments", "Env Name", "Server(s)"})
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1, 2})
	table.SetRowLine(true)
	env := 0
	for idx, group := range config.Groups {
		env = len(group.Environment)
		for _, envs := range group.Environment {
			table.Append([]string{fmt.Sprint(idx + 1), group.Name, fmt.Sprint(env), envs.Name, fmt.Sprint(len(envs.Servers))})
		}
	}
	table.Render()
}

func ListServers(group string) {
	var config c.Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Environment", "Server", "IP", "Key status"})
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	// ANSI escape code for colors
	const (
		green = "\033[32m"
		red   = "\033[31m"
		reset = "\033[0m"
	)

	// Checkmark and cross characters
	const (
		checkmark = "✓"
		cross     = "✗"
	)

	for _, grp := range config.Groups {
		if grp.Name == group {
			for _, env := range grp.Environment {
				for _, server := range env.Servers {
					status := ""
					if server.KeyAuth {
						status = green + checkmark + reset
					} else {
						status = red + cross + reset
					}
					table.Append([]string{env.Name, server.HostName, server.IP, status})
				}

			}
		}
	}
	table.Render()
}
