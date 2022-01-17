package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

var Values *config

func init() {
	envPath := os.Getenv("ENV_PATH")

	if envPath != "" {
		envPath += "/"
	}

	viper.SetConfigFile(envPath + ".env")

	var result map[string]interface{}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&result)
	if err != nil {
		fmt.Printf("Unable to decode into map, %v", err)
	}

	decErr := mapstructure.Decode(result, &Values)

	if decErr != nil {
		fmt.Println("error decoding")
	}

	Values.MaxRetries, _ = strconv.Atoi(Values.MaxRetriesStr)
	Values.BatteryAlertPercentage, _ = strconv.Atoi(Values.BatteryAlertPercentageStr)

	Values.IsDebug = Values.DebugStr == "1"

	Values.Uri.SetCommand = "/goform/goform_set_cmd_process"
	Values.Uri.GetCommand = "/goform/goform_get_cmd_process"
}

type config struct {
	UserName                  string `mapstructure:"user_name"`
	Password                  string `mapstructure:"password"`
	BaseUrl                   string `mapstructure:"base_url"`
	MaxRetriesStr             string `mapstructure:"max_retries"`
	BatteryAlertPercentageStr string `mapstructure:"battery_alert_percentage"`
	SlackToken                string `mapstructure:"slack_token"`
	DebugStr                  string `mapstructure:"debug"`
	SlackChannelId            string `mapstructure:"slack_channel_id"`
	IsDebug                   bool
	Uri                       uri
	MaxRetries                int
	BatteryAlertPercentage    int
}

type uri struct {
	SetCommand string
	GetCommand string
}
