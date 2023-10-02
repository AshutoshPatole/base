package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/TwiN/go-color"
	"golang.org/x/crypto/ssh"

	_ "embed"
)

//go:embed bin/ssh.exe
var sshExe []byte

func Connect(server, user, environment string) {
	home, _ := os.UserHomeDir()

	privKey := home + "/.ssh/id_ed25519"

	_, err := os.ReadFile(privKey)

	if err != nil {
		fmt.Println(color.InRed("Failed to read private key: " + err.Error()))
		return
	}

	var promptColor string
	fmt.Println(promptColor)
	// set bash prompt colors
	if environment == "prd" {
		promptColor = "\\[\\033[1;31m\\]\\[\\033m\\]\\u@\\h \\[\\033[1;36m\\]\\w\\[\\033[0m\\]\\$ "
	} else if environment == "dev" {
		promptColor = "\\[\\033[1;32m\\]\\[\\033m\\]\\u@\\h \\[\\033[1;36m\\]\\w\\[\\033[0m\\]\\$ "
	} else if environment == "uat" {
		promptColor = "\\[\\033[1;34m\\]\\[\\033m\\]\\u@\\h \\[\\033[1;36m\\]\\w\\[\\033[0m\\]\\$ "
	} else if environment == "sit" {
		promptColor = "\\[\\033[1;34m\\]\\[\\033m\\]\\u@\\h \\[\\033[1;36m\\]\\w\\[\\033[0m\\]\\$ "
	} else if environment == "ppd" {
		promptColor = "\\[\\033[1;33m\\]\\[\\033m\\]\\u@\\h \\[\\033[1;36m\\]\\w\\[\\033[0m\\]\\$ "
	}

	var sshCmd *exec.Cmd

	// check which platform
	platform := runtime.GOOS

	if platform == "windows" {

		/*
			Reason why I did this is because there is permission issues while accessing ssh stored in
			C:\Windows\ directory. hence I made a local copy in ssh/bin/ssh.exe folder and embeded it in go build
			HACK!
		*/
		tmpDir, err := os.MkdirTemp("", "embedded_ssh")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(tmpDir)

		// Write the embedded ssh.exe data to a temporary file
		sshExePath := filepath.Join(tmpDir, "ssh.exe")
		err = os.WriteFile(sshExePath, sshExe, 0755)
		if err != nil {
			panic(err)
		}

		// Now you have ssh.exe in a temporary file. You can execute it.
		fmt.Print("SSH version : ")
		cmd := exec.Command(sshExePath, "-V")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		// Construct a single SSH command to set the new PS1 configuration in ~/.bashrc
		exec.Command(
			sshExePath,
			user+"@"+server,
			fmt.Sprintf(`sed -i '/^export PS1=/d' ~/.bashrc && echo 'export PS1="%s"' >> ~/.bashrc`, promptColor),
		).Run()

		// Modify the SSH command to set the prompt colors
		sshCmd = exec.Command(sshExePath, user+"@"+server)
	} else if platform == "linux" {
		/*
			for linux platform it should work just fine!
		*/
		exec.Command(
			"ssh",
			user+"@"+server,
			fmt.Sprintf(`sed -i '/^export PS1=/d' ~/.bashrc && echo 'export PS1="%s"' >> ~/.bashrc`, promptColor),
		).Run()

		// Modify the SSH command to set the prompt colors
		sshCmd = exec.Command("ssh", user+"@"+server)
	}

	// Set the standard input, output, and error streams to the current process's streams
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	ssh_err := sshCmd.Run()
	if ssh_err != nil {
		fmt.Println("Failed to run SSH command:", ssh_err)
	}
}

func EstablishConnection(user, server string) (*ssh.Client, error) {
	home, _ := os.UserHomeDir()

	privKey := home + "/.ssh/id_ed25519"

	privKeyByte, err := os.ReadFile(privKey)

	if err != nil {
		fmt.Println(color.InRed("Failed to read private key: " + err.Error()))
	}

	signer, _ := ssh.ParsePrivateKey(privKeyByte)

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},

		// TODO: fix insecureignore host callback method
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// optional tcp connect timeout
		Timeout: 5 * time.Second,
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", server+":22", config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
