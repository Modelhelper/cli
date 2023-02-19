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

// setTemplateLocationCmd represents the setTemplateLocation command
var setTemplateLocationCmd = &cobra.Command{
	Use:     "templateLocation",
	Aliases: []string{"tl", "loc", "templatelocation", "Templatelocation"},
	Short:   "Sets the template location in the config file",
	Long: `usage:
mh set templatelocation [path]

if path argument is blank/empty, the cli assumes you will add current working dir

if mh is executed from /c/dev/modelhelper/t/templates the following command will set '/c/dev/modelhelper/t/templates' 
as current template location in the config file.

mh set templatelocation
	`,
	Run: func(cmd *cobra.Command, args []string) {

		var path string

		if len(args) == 0 {
			path, _ = os.Getwd()
		} else {
			path = args[0]
		}

		err := config.SetTemplateLocation(path)

		if err != nil {
			log.Fatalln("Could not set template location in config", err)
		}

		fmt.Println(path)
	},
}

func init() {
	setCmd.AddCommand(setTemplateLocationCmd)
}
