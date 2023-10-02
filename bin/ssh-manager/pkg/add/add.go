/*
Copyright © 2023 AshutoshPatole
*/
package add

import (
	"github.com/spf13/cobra"
)

var allowedEnvironments = []string{"dev", "uat", "sit", "ppd", "prd"}

// AddCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add servers and groups in configurations",
	Long: `Servers and groups can be created in the configuration using add
	`,
	// Run: func(cmd *cobra.Command, args []string) {

	// },
}
