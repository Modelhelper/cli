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
	"log"

	"github.com/spf13/cobra"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// aboutCmd represents the about command
var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		defer ui.Close()

		table1 := widgets.NewTable()
		table1.Rows = [][]string{
			[]string{"header1", "header2", "header3"},
			[]string{" h", "Go-lang is so cool", "Im working on Ruby"},
			[]string{"2016", "10", "11"},
		}
		table1.TextStyle = ui.NewStyle(ui.ColorWhite)
		table1.SetRect(0, 0, 60, 10)

		ui.Render(table1)

		table2 := widgets.NewTable()
		table2.Rows = [][]string{
			[]string{"header1", "header2", "header3"},
			[]string{"Foundations", "Go-lang is so cool", "Im working on Ruby"},
			[]string{"2016", "11", "11"},
		}
		table2.TextStyle = ui.NewStyle(ui.ColorWhite)
		table2.TextAlignment = ui.AlignCenter
		table2.RowSeparator = false
		table2.SetRect(0, 10, 20, 20)

		ui.Render(table2)

		table3 := widgets.NewTable()
		table3.Rows = [][]string{
			[]string{"header1", "header2", "header3"},
			[]string{"AAA", "BBB", "CCC"},
			[]string{"DDD", "EEE", "FFF"},
			[]string{"GGG", "HHH", "III"},
		}
		table3.TextStyle = ui.NewStyle(ui.ColorWhite)
		table3.RowSeparator = true
		table3.BorderStyle = ui.NewStyle(ui.ColorGreen)
		table3.SetRect(0, 30, 70, 20)
		table3.FillRow = true
		table3.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
		table3.RowStyles[2] = ui.NewStyle(ui.ColorWhite, ui.ColorRed, ui.ModifierBold)
		table3.RowStyles[3] = ui.NewStyle(ui.ColorYellow)

		ui.Render(table3)

	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// aboutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// aboutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
