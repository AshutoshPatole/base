package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type ipinfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Timezone string `json:"timezone"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iplocator",
	Short: "Command to find information about IP address",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			fmt.Println(errors.New("Error: At least one IP address is required."))
			return
		}

		// Create a tablewriter.Writer instance
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(
			[]string{"IP", "Organisation", "City", "Region", "Country", "Location", "TimeZone"},
		)

		for _, ip := range args {
			getIPInfo(ip, table)
		}

		// Render the table after processing all IPs
		table.Render()
	},
}

func getIPInfo(ip string, table *tablewriter.Table) {
	response, err := http.Get("http://ipinfo.io/" + ip)
	if err != nil {
		log.Fatalln("error found with IP " + ip)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln("received broken response from the server")
	}

	var data ipinfo
	err = json.Unmarshal([]byte(body), &data)

	if data.Loc == "" {
		fmt.Println("could not find information about " + ip)
	} else {
		table.Append([]string{data.IP, data.Org, data.City, data.Region, data.Country, data.Loc, data.Timezone})
	}
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.iplocator.yaml)")
	// rootCmd.PersistentFlags().StringVarP(, "version", "1.0", "version of the iplocator")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
