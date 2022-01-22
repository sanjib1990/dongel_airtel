package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"sync"
)

var Values *config
var confMu = &sync.Mutex{}
var cfg map[string]string

func get(key string, _default string) string {
	confMu.Lock()
	defer confMu.Unlock()
	if val, ok := cfg[key]; ok {
		return val
	}
	val, ok := os.LookupEnv(key)
	if !ok {
		return _default
	}
	return val
}

func init() {
	env, _ := godotenv.Read()

	if env == nil {
		env = make(map[string]string)
	}

	// Assigining env to config
	cfg = env
	Values = &config{}
	Values.UserName = get("user_name", "")
	Values.Password = get("password", "")
	Values.BaseUrl = get("base_url", "")
	Values.SlackToken = get("slack_token", "")
	Values.SlackChannelId = get("slack_channel_id", "")
	Values.MaxRetries, _ = strconv.Atoi(get("max_retries", "5"))
	Values.BatteryAlertPercentage, _ = strconv.Atoi(get("battery_alert_percentage", "20"))
	Values.IsDebug = get("debug", "0") == "1"

	Values.Uri.SetCommand = "/goform/goform_set_cmd_process"
	Values.Uri.GetCommand = "/goform/goform_get_cmd_process"
}

type config struct {
	UserName               string `mapstructure:"user_name"`
	Password               string `mapstructure:"password"`
	BaseUrl                string `mapstructure:"base_url"`
	SlackToken             string `mapstructure:"slack_token"`
	SlackChannelId         string `mapstructure:"slack_channel_id"`
	IsDebug                bool
	Uri                    uri
	MaxRetries             int
	BatteryAlertPercentage int
}

type uri struct {
	SetCommand string
	GetCommand string
}
