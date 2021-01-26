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
	"modelhelper/cli/input"
	table "modelhelper/cli/ui"

	_ "github.com/gookit/color"
	"github.com/spf13/cobra"
)

// entityCmd represents the entity command
var entityCmd = &cobra.Command{
	Use:     "entity",
	Aliases: []string{"e"},

	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("entity called")
		src := source

		if len(source) == 0 {
			src = getSourceName()
		}

		input := input.GetSource(src, mhConfig)

		if len(args) > 0 {
			en := args[0]
			e, err := input.Entity(en)
			if err != nil {

			}

			if e == nil {
				fmt.Println("The entity could not be found")
			}

			fmt.Println(src, e.Name, e.Type, e.Schema, e.Description)
			// sss := c.Green.String()
			// fmt.Println(sss)
			// headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
			// columnFmt := color.New(color.FgYellow).SprintfFunc()

			tbl := table.New("Name", "Type", "Description")
			//tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			for _, c := range e.Columns {
				tbl.AddRow(c.Name, c.DataType, c.Description)
			}

			tbl.Print()
		} else {
			ents, _ := input.Entities("")
			tbl := table.New("Name", "Type", "Description")
			//tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			for _, c := range *ents {
				tbl.AddRow(c.Name, c.Schema, c.Description)
			}
			tbl.Print()
		}

	},
}

func init() {
	rootCmd.AddCommand(entityCmd)
}

func getSourceName() string {
	defaultSource := mhConfig.DefaultSource

	if len(defaultSource) == 0 {
		if len(mhConfig.Sources) == 0 {
			defaultSource = ""
		} else {
			for _, s := range mhConfig.Sources {

				defaultSource = s.Name
				break
			}
		}

	}

	return defaultSource
}
