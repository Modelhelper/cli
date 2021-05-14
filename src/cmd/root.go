/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"modelhelper/cli/app"
	"os"

	"github.com/spf13/cobra"
)

var modelHelperApp *app.Application

var cfgFile string

// var mhConfig config.Config
// var source string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mh",
	Short: "Shows information about the ModelHelper CLI",

	Run: func(cmd *cobra.Command, args []string) {

		printLogoInfo()
		// rootExists := config.LocationExists()

		// if rootExists == false {
		// 	ex, err := os.Executable()
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	// Path to executable file
		// 	fmt.Println(ex)

		// 	// Resolve the direcotry
		// 	// of the executable
		// 	exPath := filepath.Dir(ex)
		// 	fmt.Println("Executable path :" + exPath)

		// 	// Use EvalSymlinks to get
		// 	// the real path.
		// 	realPath, err := filepath.EvalSymlinks(exPath)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	fmt.Println("Symlink evaluated:" + realPath)

		// 	dir, err := os.UserHomeDir()
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	fmt.Println(dir)
		// }
	},
}

func SetApplication(app *app.Application) {
	modelHelperApp = app
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
}

// func GetConfig() *config.Config {
// 	return &mhConfig
// }
func init() {
	// cobra.OnInitialize(initConfig)
	//rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "Sets the source")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// configPath := config.Location()
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

	// err = viper.Unmarshal(&mhConfig)
	// if err != nil {
	// 	// t.Fatalf("unable to decode into struct, %v", err)
	// }
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}

	// 	// Search config in home directory with name ".cli" (without extension).
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(".cli")
	// }

	// viper.AutomaticEnv() // read in environment variables that match

	// // If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }

	// defaultSource := mhConfig.DefaultSource

	// if len(defaultSource) == 0 {
	// 	if len(mhConfig.Sources) == 0 {
	// 		defaultSource = ""
	// 	} else {
	// 		for _, s := range mhConfig.Sources {

	// 			defaultSource = s.Name
	// 			break
	// 		}
	// 	}

	// }

	//app.SetConfig(mhConfig)

}
