/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"dongel/Service"
	"github.com/spf13/cobra"
)

// smsCmd represents the sms command
var smsCmd = &cobra.Command{
	Use:   "sms",
	Short: "Check SMS",
	Long:  `Check available SMS and do necessary actions if required`,
	Run: func(cmd *cobra.Command, args []string) {
		shoudDelete := rootCmd.PersistentFlags().Lookup("delete").Value.String() == "true"
		Service.Sms(shoudDelete)
	},
}

func init() {
	rootCmd.AddCommand(smsCmd)
}
