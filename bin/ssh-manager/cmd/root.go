/*
Copyright © 2023 AshutoshPatole
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/AshutoshPatole/ssh-manager/pkg/add"
	"github.com/AshutoshPatole/ssh-manager/pkg/connect"
	"github.com/AshutoshPatole/ssh-manager/pkg/importer"
	"github.com/AshutoshPatole/ssh-manager/pkg/list"
	"github.com/AshutoshPatole/ssh-manager/pkg/rcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssm",
	Short: "CLI for managing SSH connections efficiently",
	Long: `A CLI utility for managing SSH servers and keys efficiently:

This utility is designed to streamline SSH key handling by categorizing servers into groups. 
To get started, users can run the 'ssm import template' command to obtain a predefined template for server configuration

	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(add.AddCmd)
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(importer.ImportCmd)
	rootCmd.AddCommand(connect.ConnectCmd)
	rootCmd.AddCommand(rcp.ReverseCopyCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ssh-manager.json)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pkg" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".ssh-manager")
		viper.SafeWriteConfig()

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		// do nothing as of now
		fmt.Println("")

	}
}
