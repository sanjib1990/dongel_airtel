package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var Values *config

func init() {
	viper.SetConfigFile(".env")

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

	Values.Uri.SetCommand = "/goform/goform_set_cmd_process"
	Values.Uri.GetCommand = "/goform/goform_get_cmd_process"

}

type config struct {
	UserName string `mapstructure:"user_name"`
	Password string `mapstructure:"password"`
	BaseUrl  string `mapstructure:"base_url"`
	Uri      uri
}

type uri struct {
	SetCommand string
	GetCommand string
}
