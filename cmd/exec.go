/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"dongel/Models"
	"dongel/Service"
	"dongel/config"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Handle Few execs",
	Long:  `Handle login -> stats -> sms -> charge status -> delete SMS`,
	Run: func(cmd *cobra.Command, args []string) {
		usr := Models.User{}

		usr.UserName = rootCmd.PersistentFlags().Lookup("username").Value.String()
		usr.Password = rootCmd.PersistentFlags().Lookup("password").Value.String()
		minBatteryVlm := rootCmd.PersistentFlags().Lookup("alert-charge").Value.String()
		shoudDelete := rootCmd.PersistentFlags().Lookup("delete").Value.String() == "true"
		viewSms := rootCmd.PersistentFlags().Lookup("view-sms").Value.String() == "true"
		alert := rootCmd.PersistentFlags().Lookup("alert").Value.String() == "true"
		var minBattery int
		step := 1
		fmt.Print("[" + strconv.Itoa(step) + "] ")
		rs := Service.Login(&usr)
		if rs.Status != 1 {
			fmt.Println("Execution Failed")
			return
		}
		step++
		fmt.Println("[" + strconv.Itoa(step) + "] Stats ")
		rsp := Service.Stats()
		if rsp.Status != 1 {
			fmt.Println("Execution Failed")
			return
		}

		if minBatteryVlm == "" {
			minBatteryVlm = config.Values.BatteryAlertPercentageStr
		}

		minBattery, _ = strconv.Atoi(minBatteryVlm)
		if alert == true && minBattery > 0 {
			step++
			fmt.Print("[" + strconv.Itoa(step) + "] ")
			rsp = Service.HandleCharge(&rsp, minBattery)
		}
		if viewSms == true {
			step++
			fmt.Println("[" + strconv.Itoa(step) + "] SMS")
			rsp = Service.Sms(shoudDelete)
			if rsp.Status != 1 {
				fmt.Println("Execution Failed")
				return
			}
		}

		step++
		fmt.Print("[" + strconv.Itoa(step) + "] ")
		Service.Logout()
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
