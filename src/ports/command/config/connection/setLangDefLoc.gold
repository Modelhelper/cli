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
	"os"

	"github.com/spf13/cobra"
)

// setLangDefLocCmd represents the setLangDefLoc command
var setLangDefLocCmd = &cobra.Command{
	Use:     "languageLocation",
	Aliases: []string{"ll", "lang", "languagelocation", "Languagelocation"},
	Short:   "Sets the path to where to find language definition files",
	Run: func(cmd *cobra.Command, args []string) {
		var path string

		if len(args) == 0 {
			path, _ = os.Getwd()
		} else {
			path = args[0]
		}

		err := config.SetLangDefLocation(path)

		if err != nil {
			log.Fatalln("Could not set language location path in config", err)
		}

		fmt.Println("Successfully updated language def path to ", path)
	},
}

func init() {
	setCmd.AddCommand(setLangDefLocCmd)
}
