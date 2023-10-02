/*
Copyright © 2023 AshutoshPatole
*/
package list

import (
	"github.com/spf13/cobra"
)

var lGroup string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "list servers from a specific group",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ListServers(lGroup)
	},
}

func init() {
	ListCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")
	serverCmd.Flags().StringVarP(&lGroup, "group", "g", "", "Specify which group servers should be displayed")
	serverCmd.MarkFlagRequired("group")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
