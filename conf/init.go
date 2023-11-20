package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

var (
	DSN         string
	ServiceAddr string
)

func Init() {
	// Define configuration file destination and type
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf")

	// Read config file, panic when fail
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	//Set database destination
	DSN = strings.Join([]string{
		viper.GetString("Database.User"),
		":",
		viper.GetString("Database.Password"),
		"@tcp(",
		viper.GetString("Database.Host"),
		":",
		viper.GetString("Database.Port"),
		")/",
		viper.GetString("Database.Name"),
		"?charset=utf8mb4&parseTime=True&loc=Local",
	}, "")

	//Set service address
	ServiceAddr = viper.GetString("Service.IP") + ":" + viper.GetString("Service.Port")
}
