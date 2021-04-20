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
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "<not implemented>",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(mhConfig.AppVersion)
		fmt.Println(mhConfig.ConfigVersion)
		fmt.Println(mhConfig.Languages.Definitions)
		fmt.Println(mhConfig.Templates.Location)

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
}
