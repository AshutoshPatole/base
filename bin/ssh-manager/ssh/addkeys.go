package ssh

import (
	"fmt"
	"os"

	"github.com/TwiN/go-color"
	"golang.org/x/crypto/ssh"
)

func AddPubKeysToServer(session *ssh.Session) bool {
	home, _ := os.UserHomeDir()

	pubKeyPath := home + "/.ssh/id_ed25519.pub"

	pubKey, err := os.ReadFile(pubKeyPath)

	if err != nil {
		fmt.Println(color.InRed("Could not read public key " + pubKeyPath))
		return false
	}

	command := fmt.Sprintf("mkdir -p ~/.ssh/; chmod 700 -R ~/.ssh; echo '%s' >> ~/.ssh/authorized_keys; chmod 600 ~/.ssh/authorized_keys", pubKey)
	if err := session.Run(command); err != nil {
		fmt.Println(color.InRed("Failed to add public key " + err.Error()))
		return false
	} else {
		fmt.Println("Added public key to server")
	}

	defer session.Close()
	return true

}
