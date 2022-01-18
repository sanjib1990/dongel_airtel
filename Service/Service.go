package Service

import (
	"context"
	"dongel/Models"
	"dongel/config"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/slack-go/slack"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func appendBatteryDetails(writer table.Writer, stats Models.Response) {
	charging := "No"
	if stats.BatteryCharging == "1" {
		charging = "Yes"
	}
	writer.AppendRow([]interface{}{"Is Battery Charging", charging})
	writer.AppendSeparator()

	writer.AppendRow([]interface{}{"Battery Percentage", stats.BatteryVolPercent})
	writer.AppendSeparator()
}

func appendHasUnreadSMS(writer table.Writer, stats Models.Response) {
	writer.AppendRow([]interface{}{"Unread message count", stats.SmsUnreadNum})
	writer.AppendSeparator()
}

func appendTotalConnections(writer table.Writer, stats Models.Response) {
	writer.AppendRow([]interface{}{"Total number of connections", stats.StaCount})
	writer.AppendSeparator()
}

func DisplayStats(stats *Models.Response) {
	t := table.NewWriter()

	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Value"})

	appendBatteryDetails(t, *stats)
	appendHasUnreadSMS(t, *stats)
	appendTotalConnections(t, *stats)

	t.Render()
}

func Login(user *Models.User) Models.Response {
	username := getUserName(&user.UserName)
	password := getPassword(&user.Password)

	encodedUserName := base64.StdEncoding.EncodeToString([]byte(*username))
	encodedPasswordName := base64.StdEncoding.EncodeToString([]byte(*password))

	data := map[string]string{
		"isTest":   "false",
		"goformId": "LOGIN",
		"username": encodedUserName,
		"password": encodedPasswordName,
	}

	return makeSetCommandCalls(data, "0", "Login: ")
}

func Stats() Models.Response {
	data := map[string]string{
		"multi_data": "1",
		"cmd":        "battery_vol_percent,battery_charging,sms_unread_num,sta_count,loginfo",
	}
	var rsp Models.Response

	makeGetCommandCalls(data, &rsp)

	if rsp.Loginfo != "ok" {
		fmt.Println("[X] Please Login...")
		return rsp
	}

	DisplayStats(&rsp)

	return rsp
}

func HandleCharge(resp *Models.Response, minBattery int) Models.Response {
	if minBattery <= 0 {
		minBattery = config.Values.BatteryAlertPercentage
	}
	var rs Models.Response
	currBttrVlm, err := strconv.Atoi(resp.BatteryVolPercent)
	if err != nil {
		fmt.Println("Cannot handle charge alert, as the value could not be determined")
		return rs
	}
	if currBttrVlm > minBattery {
		fmt.Println("Battery is above thrashold")
		return rs
	}

	// send Slack notification
	api := slack.New(config.Values.SlackToken, slack.OptionDebug(config.Values.IsDebug))
	channelMessage := "<@sanjib> Charge For your Airtel Modem is below " + strconv.Itoa(minBattery) + ". Please put it for charging."

	if config.Values.SlackChannelId == "" {
		return rs
	}

	message, _, err := api.PostMessage(config.Values.SlackChannelId,
		slack.MsgOptionText(channelMessage, false),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		fmt.Println("[Slack Error] ", message, err.Error())
		return rs
	}

	rs.Status = 1

	return rs
}

func Sms(delete bool) Models.Response {
	data := map[string]string{
		"isTest":        "false",
		"cmd":           "sms_data_total",
		"page":          "0",
		"data_per_page": "500",
		"mem_store":     "1",
		"tags":          "10",
		"order_by":      "order+by+id+desc",
	}

	var rsp Models.Response
	makeGetCommandCalls(data, &rsp)

	if rsp.SmsDataTotal == "" {
		fmt.Println("Checking Stats..")
		Stats()
		return rsp
	}

	DisplaySmsList(rsp.Messages)

	if delete == true {
		DeleteSmsByIds(rsp.Messages)
	}

	return rsp
}

func DisplaySmsList(smsList []Models.Sms) {
	writer := table.NewWriter()

	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"SMS ID", "From Contact Number", "Content"})

	for _, sms := range smsList {
		cont, err := hex.DecodeString(sms.Content)
		if err != nil {
			continue
		}

		writer.AppendRow([]interface{}{sms.Id, sms.Number, string(cont)})
		writer.AppendSeparator()
	}

	writer.Render()
}

func DeleteSmsByIds(smsList []Models.Sms) {
	ids := ""

	for _, sms := range smsList {
		ids += sms.Id + ";"
	}

	data := map[string]string{
		"isTest":        "false",
		"goformId":      "DELETE_SMS",
		"msg_id":        ids,
		"data_per_page": "500",
		"notCallback":   "true",
	}

	fmt.Println("Deleting SMS by Id: " + ids)

	makeSetCommandCalls(data, "success", "SMS Deleted: ")
}

func getContexedRequest(urlStr string, method string, body io.Reader) (*http.Request, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	return req, cancel, err
}

func makeApiCall(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	go func() {
		select {
		case <-time.After(100 * time.Millisecond):
			fmt.Print("Timed Out --> ")
		case <-request.Context().Done():
		}
	}()
	return client.Do(request)
}

func makeGetCommandCalls(queryData map[string]string, obj *Models.Response) {
	apiUrl := config.Values.BaseUrl
	resource := config.Values.Uri.GetCommand

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()
	req, cancel, _ := getContexedRequest(urlStr, http.MethodGet, nil)
	defer cancel()
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	query := req.URL.Query()
	for key, value := range queryData {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := makeApiCall(req)

	if err != nil {
		fmt.Println("Request Failed")

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error occured !", err)
		return
	}

	err = json.Unmarshal(bodyBytes, &obj)

	if err != nil {
		fmt.Println("Invalid response")
		return
	}

	obj.Status = 1
}

func makeSetCommandCalls(formData map[string]string, resultCheck string, message string) Models.Response {
	apiUrl := config.Values.BaseUrl
	resource := config.Values.Uri.SetCommand
	var obj Models.Response

	data := url.Values{}
	for key, value := range formData {
		data.Set(key, value)
	}

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	r, cancel, _ := getContexedRequest(urlStr, http.MethodPost, strings.NewReader(data.Encode()))
	defer cancel()

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := makeApiCall(r)

	if err != nil {
		fmt.Println(message + "Fail")
		return obj
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("Error occured !", err)
			return obj
		}

		err = json.Unmarshal(bodyBytes, &obj)

		if err != nil {
			fmt.Println("Invalid response")
			return obj
		}

		status := "Success"

		if obj.Result != resultCheck {
			status = "Fail"
		}

		fmt.Println(message + status)
	}

	obj.Status = 1

	return obj
}

func Logout() Models.Response {
	data := map[string]string{
		"isTest":   "false",
		"goformId": "LOGOUT",
	}

	return makeSetCommandCalls(data, "success", "Logout: ")
}

func getUserName(username *string) *string {
	if username == nil || len(strings.TrimSpace(*username)) == 0 {
		return &config.Values.UserName
	}
	return username
}

func getPassword(password *string) *string {
	if password == nil || len(strings.TrimSpace(*password)) == 0 {
		return &config.Values.Password
	}

	return password
}
