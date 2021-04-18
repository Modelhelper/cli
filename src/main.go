package main

import (
	// "fmt"

	"modelhelper/cli/app"
	"modelhelper/cli/cmd"
	"modelhelper/cli/config"
	// "github.com/spf13/viper"
)

var AppConfig string = "Something"

func main() {

	rootExists := config.LocationExists()

	if rootExists == false {
		cfg := config.New()

		// initialize config
		err := app.Initialize(cfg)

		if err != nil {
			panic(err)
		}
	} else {
		// configPath := common.ConfigFolder()
		// viper.SetConfigName("config")
		// viper.SetConfigType("yaml")

		// viper.AddConfigPath(configPath) // optionally look for config in the working directory
		// err := viper.ReadInConfig()     // Find and read the config file
		// if err != nil {                 // Handle errors reading the config file
		// 	panic(fmt.Errorf("Fatal error config file: %s \n", err))
		// }

		// if err := viper.ReadInConfig(); err != nil {
		// 	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// 		// Config file not found; ignore error if desired
		// 	} else {
		// 		// Config file was found but another error was produced
		// 	}
		// }

		cmd.Execute()
	}

}
