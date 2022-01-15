// Package cmd /*
package cmd

import (
	"dongel/Models"
	"dongel/Service"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Dongel",
	Long:  `Login to Dongel`,
	Run: func(cmd *cobra.Command, args []string) {
		usr := Models.User{}

		usr.UserName = rootCmd.PersistentFlags().Lookup("username").Value.String()
		usr.Password = rootCmd.PersistentFlags().Lookup("password").Value.String()

		Service.Login(&usr)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
