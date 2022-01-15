package Models

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
	SmsDataTotal      interface{} `json:"sms_data_total"`
	Status            int
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
