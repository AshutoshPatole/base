/*
Copyright © 2023 AshutoshPatole
*/
package importer

import (
	"fmt"
	"os"

	"github.com/AshutoshPatole/ssh-manager/ssh"
	"github.com/TwiN/go-color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var filePath, iGroup string

type Group struct {
	Name        string        `yaml:"name"`
	User        string        `yaml:"user"`
	Environment []Environment `yaml:"env"`
}

type Environment struct {
	Name    string   `yaml:"name"`
	Servers []Server `yaml:"servers"`
}

type Server struct {
	HostName string `yaml:"hostname"`
	Alias    string `yaml:"alias"`
}

// importCmd represents the import command
var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import file",
	Long: `
You can import your groups and server configurations using YAML file (or)
You can add servers individually using ssm add server command. 

To import config in bulk:
ssm import -f $HOME/import.yaml --all

To import specific group:
ssm import -f $HOME/import.yaml -g groupName

To download a template YAML file. This will save the template in ~/Downloads folder
ssm import template 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		isImportAll, _ := cmd.Flags().GetBool("all")
		if filePath != "" {
			readFile(filePath, iGroup, isImportAll)
		}
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")
	ImportCmd.Flags().StringVarP(&filePath, "file", "f", "", "Specify YAML file path to import (required)")
	ImportCmd.MarkFlagRequired("file")
	ImportCmd.Flags().StringVarP(&iGroup, "group", "g", "", "Specify which group to include from the import")
	ImportCmd.Flags().BoolP("all", "a", false, "Import all groups")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func readFile(filepath, group string, allGroup bool) {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}
	var data []Group
	if err := yaml.Unmarshal(yamlFile, &data); err != nil {
		fmt.Printf("Error unmarshaling YAML: %s\n", err)
		return
	}

	if !allGroup && group == "" {
		fmt.Println(color.InRedOverBlack("Either specify group name or use --all flag"))
		return
	}

	if allGroup {
		for _, grp := range data {
			fmt.Println("Importing " + grp.Name)
			password := ssh.AskPass()
			for _, env := range grp.Environment {
				for _, host := range env.Servers {
					ssh.InitServer(host.HostName, grp.User, password, grp.Name, env.Name, host.Alias)
				}
			}
		}
	}
	if group != "" {
		for _, grp := range data {
			if grp.Name == group {
				fmt.Println("Importing " + grp.Name)
				password := ssh.AskPass()
				for _, env := range grp.Environment {
					for _, host := range env.Servers {
						ssh.InitServer(host.HostName, grp.User, password, grp.Name, env.Name, host.Alias)
					}
				}
			}
		}
	}
}
