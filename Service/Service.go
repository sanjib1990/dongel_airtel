package Service

import (
	"dongel/Models"
	"dongel/config"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Login(user *Models.User) {
	apiUrl := config.Values.BaseUrl
	resource := config.Values.Uri.SetCommand

	data := url.Values{}

	username := getUserName(user.UserName)
	password := getPassword(user.Password)

	encodedUserName := base64.StdEncoding.EncodeToString([]byte(username))
	encodedPasswordName := base64.StdEncoding.EncodeToString([]byte(username))

	data.Set("isTest", "false")
	data.Set("goformId", "LOGIN")
	data.Set("username", encodedUserName)
	data.Set("password", encodedPasswordName)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}

	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	r.Header.Add("Sec-GPC", "1")
	/**
		Host:192.168.1.1
	Origin:http://192.168.1.1
	Referer:http://192.168.1.1/index.html
	Sec-GPC:1
	User-Agent:Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36
	X-Requested-With:XMLHttpRequest
	*/
	r.Header.Add("Origin", "http://192.168.1.1")
	r.Header.Add("Host", "192.168.1.1")
	r.Header.Add("Referer", "http://192.168.1.1/index.html")
	r.Header.Add("X-Requested-With", "XMLHttpRequest")
	r.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")

	resp, _ := client.Do(r)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("Error occured !", err)
			return
		}

		str := string(bodyBytes)
		fmt.Println(str)
	}
}

func getUserName(username *string) *string {
	if username == nil {
		return *config.Values.UserName
	}

	return username
}

func getPassword(password string) string {
	return ""
}
