package Models

import (
	"dongel/config"
	"fmt"
	"github.com/slack-go/slack"
)

type User struct {
	UserName string
	Password string
}

type Response struct {
	Result            string      `json:"result"`
	Messages          []Sms       `json:"messages"`
	BatteryVolPercent string      `json:"battery_vol_percent"`
	BatteryCharging   string      `json:"battery_charging"`
	SmsUnreadNum      string      `json:"sms_unread_num"`
	StaCount          string      `json:"sta_count"`
	Loginfo           string      `json:"loginfo"`
	SubNetworkType    string      `json:"sub_network_type"`
	NetworkType       string      `json:"network_type"`
	UploadRate        string      `json:"realtime_tx_thrpt"`
	DownloadRate      string      `json:"realtime_rx_thrpt"`
	NetworkProvider   string      `json:"network_provider"`
	SignalBar         string      `json:"signalbar"`
	SmsDataTotal      interface{} `json:"sms_data_total"`
	Status            int
	RenderedTxt       string
	RenderedMd        string
}

func (receiver Response) init() {
	receiver.SmsDataTotal = nil
}

type Sms struct {
	Id           string `json:"id"`
	Number       string `json:"number"`
	Content      string `json:"content"`
	Tag          string `json:"tag"`
	Date         string `json:"date"`
	DraftGroupId string `json:"draft_group_id"`
}

type SlackNotification struct {
}

func (sn SlackNotification) getInstance() *slack.Client {
	return slack.New(config.Values.SlackToken, slack.OptionDebug(config.Values.IsDebug))
}

func (sn SlackNotification) SendNotification(channelId string, msg string) int {
	message, _, err := sn.getInstance().PostMessage(channelId,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAsUser(true))

	if err != nil {
		fmt.Println("[Slack Error] ", message, err.Error())
		return 0
	}

	return 1
}
