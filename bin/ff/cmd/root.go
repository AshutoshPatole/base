package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/spf13/cobra"
)

var filePattern string

var rootCmd = &cobra.Command{
	Use:   "ff",
	Short: "find files",
	Long:  `find files`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(color.InRedOverBlack("Expects a file pattern"))
			os.Exit(1)
		}
		filePattern = args[0]

		// Check if the filePattern contains wildcard '*'
		if strings.Contains(filePattern, "*") {
			traverseWildcard()
		} else {
			traverseExact()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// add flags here
}

func traverseWildcard() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	err = filepath.WalkDir(cwd, visit)
	if err != nil {
		fmt.Println("Error walking directory:", err)
	}
}

func traverseExact() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	err = filepath.WalkDir(cwd, visitExact)
	if err != nil {
		fmt.Println("Error walking directory:", err)
	}
}

func visit(path string, info os.DirEntry, err error) error {
	if err != nil {
		if pErr, ok := err.(*os.PathError); ok && os.IsPermission(pErr.Unwrap()) {
			return nil
		}
		return err
	}
	file := info.Name()

	// Convert the filePattern to a regular expression
	pattern := strings.Replace(filePattern, "*", ".*", -1)
	pattern = "^" + pattern + "$"

	regEx := regexp.MustCompile(pattern)

	if regEx.MatchString(file) {
		fmt.Println(path)
	}
	return nil
}

func visitExact(path string, info os.DirEntry, err error) error {
	if err != nil {
		if pErr, ok := err.(*os.PathError); ok && os.IsPermission(pErr.Unwrap()) {
			return nil
		}
		return err
	}
	file := info.Name()

	if file == filePattern {
		fmt.Println(path)
	}
	return nil
}
