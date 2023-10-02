/*
Copyright © 2023 AshutoshPatole
*/
package rcp

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	cs "github.com/AshutoshPatole/ssh-manager/ssh"
	c "github.com/AshutoshPatole/ssh-manager/utils"
	"github.com/TwiN/go-color"
	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

var cRemote string

// reverseCopyCmd represents the reverseCopy command
var ReverseCopyCmd = &cobra.Command{
	Use:     "reverse-copy",
	Short:   "Download file from remote machine",
	Aliases: []string{"rcp"},
	Long:    `Download file from remote machine. Default location for saving is $HOME`,
	Run: func(cmd *cobra.Command, args []string) {
		download()
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reverseCopyCmd.PersistentFlags().String("foo", "", "A help for foo")
	ReverseCopyCmd.Flags().StringVarP(&cRemote, "file", "f", "", "Location of remote file")
	ReverseCopyCmd.MarkFlagRequired("file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reverseCopyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func download() {
	user, server := listServers()
	client, err := cs.EstablishConnection(user, server)
	if err != nil {
		fmt.Println("Could not connect to server", err.Error())
		os.Exit(1)
	}
	home, _ := os.UserHomeDir()
	// Extract the file name and extension from the remote path
	fileName := filepath.Base(cRemote)
	localPath := filepath.Join(home, fileName)

	err = downloadFileSFTP(client, cRemote, localPath)
	if err != nil {
		fmt.Println("Error downloading file : ", err.Error())
		os.Exit(1)
	} else {
		fmt.Println(color.InGreen("The file is copied to " + localPath))
	}

}

func downloadFileSFTP(client *ssh.Client, remotePath, localPath string) error {
	// Create an SFTP session
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	// Open the remote file
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	// Create the local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// Copy the remote file's contents to the local file
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return err
	}

	return nil
}

type serverOpts struct {
	Label    string
	User     string
	HostName string
}

func listServers() (string, string) {
	var config c.Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	selectedHostName := ""
	selectedUserName := ""

	serverOptions := []serverOpts{}

	for _, grp := range config.Groups {
		for _, env := range grp.Environment {
			for _, server := range env.Servers {
				serverOption := serverOpts{
					Label:    fmt.Sprintf("%s (%s)", server.Alias, env.Name),
					User:     grp.User,
					HostName: server.HostName,
				}
				serverOptions = append(serverOptions, serverOption)
			}
		}

	}
	labels := make([]string, len(serverOptions))
	for i, serverOption := range serverOptions {
		labels[i] = serverOption.Label
	}

	prompt := &survey.Select{
		Message: "Select server",
		Options: labels,
	}
	survey.AskOne(prompt, &selectedHostName)

	// Extract environment name from the selected option
	for _, serverOption := range serverOptions {
		if serverOption.Label == selectedHostName {
			selectedUserName = serverOption.User
			selectedHostName = strings.Split(serverOption.HostName, " (")[0]
			break
		}
	}

	fmt.Println(color.InBlue("Trying to connect to " + selectedHostName))
	return selectedUserName, selectedHostName
}
