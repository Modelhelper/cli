/*
Copyright © 2020 Hans-Petter Eitvet

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"modelhelper/cli/config"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "<not implemented>",
	Run: func(cmd *cobra.Command, args []string) {

		// cfg := modelHelperApp.Configuration
		// fmt.Println(cfg.AppVersion)
		// fmt.Println(cfg.ConfigVersion)
		// fmt.Println(cfg.Languages.Definitions)
		// fmt.Println(cfg.Templates.Location)

		open, _ := cmd.Flags().GetBool("open")
		if open {
			loc := filepath.Join(config.Location(), "config.yaml")
			exe := exec.Command("code", loc)
			if exe.Run() != nil {
				//vim didn't exit with status code 0
			}
		}
		// for _, source := range mhConfig.Sources {
		// 	fmt.Println(source.Name)

		// 	for _, opt := range source.Options {
		// 		fmt.Println(opt)
		// 	}

		// 	for _, grp := range source.Groups {
		// 		for _, itm := range grp.Items {

		// 			fmt.Println(itm)
		// 		}
		// 	}
		// }

		// is := []string{"test_1", "test_2", "test_3"}
		// g := config.SourceGroup{
		// 	Items: is,
		// }

		// gg := make(map[string]config.SourceGroup)
		// gg["basic"] = g

		// s := config.Source{
		// 	Name:   "testing",
		// 	Schema: "dbo",
		// 	Groups: gg,
		// }
		// mhConfig.Sources["testing"] = s

		// viper.Set("AppVersion", 3)
		// viper.Set("Sources", mhConfig.Sources)
		// viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().Bool("open", false, "Opens the config file in VS Code")
}
