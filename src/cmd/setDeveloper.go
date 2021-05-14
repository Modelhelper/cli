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
	"log"
	"modelhelper/cli/config"

	"github.com/spf13/cobra"
)

// setDeveloperCmd represents the setDeveloper command
var setDeveloperCmd = &cobra.Command{
	Use:     "developer",
	Aliases: []string{"d", "dev"},
	Short:   "Sets developer and/or email",
	Run: func(cmd *cobra.Command, args []string) {
		devName, _ := cmd.Flags().GetString("name")
		devEMail, _ := cmd.Flags().GetString("email")

		err := config.SetDeveloper(devName, devEMail)

		if err != nil {
			log.Fatalln("Could not set developer params in config", err)
		}

		fmt.Printf("Developer: %s with %s is set in config", devName, devEMail)
	},
}

func init() {
	setCmd.AddCommand(setDeveloperCmd)

	setDeveloperCmd.Flags().String("name", "", "Sets the developer name in config.")
	setDeveloperCmd.Flags().String("email", "", "Sets the developer email in config.")

}
