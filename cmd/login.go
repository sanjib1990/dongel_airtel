// Package cmd /*
package cmd

import (
	"dongel/Service"
	"dongel/config"
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Dongel",
	Long:  `Login to Dongel`,
	Run: func(cmd *cobra.Command, args []string) {
		Service.Login()

		fmt.Println("Login Has been called", config.Values.UserName)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().String("username", "", "Username to login with. Default will be taken from the config / env")
	loginCmd.PersistentFlags().String("password", "", "Password to login with. Default will be taken from the config / env")
}
