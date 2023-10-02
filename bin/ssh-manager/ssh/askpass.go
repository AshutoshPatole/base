package ssh

import (
	"fmt"
	"os"
	"syscall"

	"github.com/TwiN/go-color"
	"golang.org/x/term"
)

func AskPass() string {
	fmt.Print("Password: ")
	bytepwd, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Print(color.InRedOverBlack(err))
		os.Exit(1)
	}
	fmt.Println("")
	password := string(bytepwd)
	return password
}
