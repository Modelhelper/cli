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

// setPortCmd represents the setPort command
var setPortCmd = &cobra.Command{
	Use:   "port",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		apiPort, _ := cmd.Flags().GetInt("api")
		webPort, _ := cmd.Flags().GetInt("web")

		err := config.SetPort(apiPort, webPort)

		if err != nil {
			log.Fatalln("Could not set developer params in config", err)
		}

		fmt.Printf("Port %v set for api and %v set for web in config", apiPort, webPort)
	},
}

func init() {
	setCmd.AddCommand(setPortCmd)

	setPortCmd.Flags().Int("api", 5000, "Sets the port for serving the api.")
	setPortCmd.Flags().Int("web", 8080, "Sets the port for serving the web application.")

}
