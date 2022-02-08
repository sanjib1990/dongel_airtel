/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"dongel/Models"
	"dongel/Service"
	"dongel/config"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var ExecSlackMessage = ""

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Handle Few execs",
	Long:  `Handle login -> stats -> sms -> charge status -> delete SMS`,
	PostRun: func(cmd *cobra.Command, args []string) {
		var slack Models.SlackNotification

		if config.Values.SlackChannelId == "" {
			return
		}
		ExecSlackMessage = "*[Dongel Automation Run]* `" + ExecSlackMessage + "`"

		slack.SendNotification(config.Values.SlackChannelId, ExecSlackMessage)
	},
	Run: func(cmd *cobra.Command, args []string) {
		usr := Models.User{}
		ExecSlackMessage = ""
		usr.UserName = rootCmd.PersistentFlags().Lookup("username").Value.String()
		usr.Password = rootCmd.PersistentFlags().Lookup("password").Value.String()
		minBatteryVlm := rootCmd.PersistentFlags().Lookup("alert-charge").Value.String()
		shoudDelete := rootCmd.PersistentFlags().Lookup("delete").Value.String() == "true"
		viewSms := rootCmd.PersistentFlags().Lookup("view-sms").Value.String() == "true"
		alert := rootCmd.PersistentFlags().Lookup("alert").Value.String() == "true"
		overChargeAlert := rootCmd.PersistentFlags().Lookup("overcharge-alert").Value.String() == "true"
		var minBattery int
		step := 1
		fmt.Print("[" + strconv.Itoa(step) + "] ")
		rsp := Service.Login(&usr)
		ExecSlackMessage += "Login -> "
		if rsp.Status != 1 {
			fmt.Println("Execution Failed")
			ExecSlackMessage += "Fail"
			return
		}
		step++
		fmt.Println("[" + strconv.Itoa(step) + "] Stats ")
		rsp = Service.Stats()
		ExecSlackMessage += "Stats -> "
		if rsp.Status != 1 {
			fmt.Println("Execution Failed")
			ExecSlackMessage += "Fail"
			return
		}

		if minBatteryVlm == "" {
			minBatteryVlm = strconv.Itoa(config.Values.BatteryAlertPercentage)
		}

		minBattery, _ = strconv.Atoi(minBatteryVlm)
		if alert && minBattery > 0 {
			step++
			ExecSlackMessage += "Low Charge Check -> "
			fmt.Println("[" + strconv.Itoa(step) + "] Low Charge Check")
			Service.HandleLowCharge(&rsp, minBattery)
			if rsp.Status != 1 {
				fmt.Println("Execution Failed")
				ExecSlackMessage += "Fail"
				return
			}
		}

		if overChargeAlert && rsp.BatteryCharging == "1" {
			step++
			ExecSlackMessage += "Over Charge Check -> "
			fmt.Println("[" + strconv.Itoa(step) + "] Over Charge Check")
			Service.HandleOverCharge(&rsp)
			if rsp.Status != 1 {
				fmt.Println("Execution Failed")
				ExecSlackMessage += "Fail"
				return
			}
		}

		if viewSms {
			step++
			fmt.Println("[" + strconv.Itoa(step) + "] SMS")
			ExecSlackMessage += "SMS -> "
			rsp = Service.Sms(shoudDelete)
			if rsp.Status != 1 {
				fmt.Println("Execution Failed")
				ExecSlackMessage += "Fail"
				return
			}
		}

		step++
		ExecSlackMessage += "Logout"
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
