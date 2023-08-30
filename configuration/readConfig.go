package configuration

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetConfig(path string) (error, string) {
	viper.SetConfigType("json")
	viper.AddConfigPath("../")
	viper.SetConfigFile("./config.json")
	fmt.Printf("Using configuration: %s\n", viper.ConfigFileUsed())
	_ = viper.ReadInConfig()
	url := viper.GetString(path)
	return nil, url
}
