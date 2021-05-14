/*
Copyright © 2021 Hans-Petter Eitvet

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
	"modelhelper/cli/config"

	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// setDefaultConnectionCmd represents the setDefaultConnection command
var setDefaultConnectionCmd = &cobra.Command{
	Use:     "defaultConnection",
	Aliases: []string{"dc", "default-con", "defcon", "DefaultConnection", "defaulconnection"},
	Short:   "Sets the default connection to be used in code generation",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := ""
		if len(args) == 0 {
			key = promptForKey()
		} else {
			key = args[0]
		}

		err := config.SetDefaultConnection(key)

		if err != nil {
			color.Red.Println(err)
			return
		}

		color.Green.Printf("\nSuccessfully updated the default connection to %s\n", key)
	},
}

func init() {
	setCmd.AddCommand(setDefaultConnectionCmd)

}

func promptForKey() string {
	cfg := config.Load()
	items := []string{}

	for k, _ := range cfg.Connections {
		items = append(items, k)
	}
	prompt := promptui.Select{
		Label: "Select one of the available connections",
		Items: items,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}
