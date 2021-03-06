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

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Root command to set parameter in the config",
	Run: func(cmd *cobra.Command, args []string) {

		m := `This is the root command for setting values to the config file. 
		
Please use one of the following sub commands to set correct information:

templatelocation | tl    : to set a template location
developer | dev          : to set developer params
 
samples:

-- Template location
usage:
mh set templatelocation [path]

if path argument is blank/empty, the cli assumes you will add current working dir

if mh is executed from /c/dev/modelhelper/t/templates the following command will set '/c/dev/modelhelper/t/templates' 
as current template location in the config file.

mh set templatelocation

-- Developer info
usage:
mh set developer [options]

mh set developer --name your name --email your@email.com
mh set dev --name name --email email

		`

		fmt.Println(m)
		fmt.Println("")
		fmt.Println("Please use one of the sub commands to set correct information")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	configCmd.AddCommand(setCmd)
}
