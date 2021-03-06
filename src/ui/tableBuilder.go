package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type TableConverter interface {
	Rows() [][]string
	Header() []string
}

func RenderTable(tc TableConverter) {

	table := CreateTableWithColumnSeparator(tc.Header())
	table.AppendBulk(tc.Rows())

	table.Render()
}

func CreateStandardTable(header []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false) // Set Border to false
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(false)
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("-")
	table.SetHeaderLine(true)

	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("-")

	return table
}

func CreateTableWithColumnSeparator(header []string) *tablewriter.Table {
	table := CreateStandardTable(header)
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("+")
	table.SetColumnSeparator("|")

	return table
}

func CreateTableWithColumnAndRowSeparator(header []string) *tablewriter.Table {
	table := CreateStandardTable(header)

	table.SetRowLine(true)
	table.SetBorder(true) // Set Border to false

	table.SetCenterSeparator("+")
	table.SetColumnSeparator("|")

	return table
}

func PrintConsoleTitle(title string) {
	t := ConsoleTitle(title)

	fmt.Printf("\n\n%v\n\n", t)
}

func ConsoleTitle(title string) string {
	ct, l := strings.ToUpper(title), 50

	return PadRight(ct, " ", l)
}

func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

func Sequence(str, pad string, length int) string {
	sl := len(str)

	v := ""
	for i := 0; i < length-sl; i++ {
		v += pad
	}

	return v
}
