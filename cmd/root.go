// Package cmd /*
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dongel",
	Short: "Airtel Dongel checks and automation",
	Long: `!! Dongel Automation !!

- Check the status of dongel and notify if the battery is low
- If there are SMS, notify and if required delete them
- Have the following functionalities
	- login to dongel
	- logout from dongel
	- check stats
	- check sms
	- delete sms
		
Command: 	
dongel <command> <flags...> <inputs..>`,
	Run: func(cmd *cobra.Command, args []string) {},
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dongel.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().String("username", "", "Username to login with. Default will be taken from the config / env")
	rootCmd.PersistentFlags().String("password", "", "Password to login with. Default will be taken from the config / env")
	rootCmd.PersistentFlags().BoolP("delete", "d", false, "delete available sms")
	rootCmd.PersistentFlags().BoolP("alert", "a", false, "Alert for low battery")
	rootCmd.PersistentFlags().BoolP("view-sms", "v", false, "View SMS")
	rootCmd.PersistentFlags().String("alert-charge", "", "Charge value below which alert will be triggered")
}
