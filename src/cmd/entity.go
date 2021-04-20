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

var skipDescription bool

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
				return
			}
			fmt.Printf("\nEntity: %s.%s", e.Schema, e.Name)
			fmt.Printf("\nCreated: %s\n", "Unknown")

			if len(e.Description) > 0 {
				fmt.Printf("\n\nDescription\n")
				fmt.Printf("------------------\n")
				fmt.Println(e.Description)
				fmt.Println()
			}
			var tbl table.Table
			if skipDescription {
				tbl = table.New("Name", "Type", "Nullable", "IsIdentity", "PK", "FK")
			} else {
				tbl = table.New("Name", "Type", "Nullable", "IsIdentity", "PK", "FK", "Description")
			}

			for _, c := range e.Columns {
				null := "false"
				if c.IsNullable {
					null = "true"
				}
				id := ""
				if c.IsIdentity {
					id = "yes"
				}

				pk := ""
				if c.IsPrimaryKey {
					pk = "PK"
				}

				fk := ""
				if c.IsForeignKey {
					fk = "FK"
				}
				if skipDescription {
					tbl.AddRow(c.Name, c.DataType, null, id, pk, fk)
				} else {
					tbl.AddRow(c.Name, c.DataType, null, id, pk, fk, c.Description)

				}
			}

			fmt.Printf("\n\nColumns\n")
			fmt.Printf("-------------------------------------------\n\n")
			tbl.Print()

			fmt.Printf("\n\nONE TO MANY (.Children)\n")
			fmt.Printf("-------------------------------------------\n\n")
			printChildTable(e.Name, e.ChildRelations)
			fmt.Printf("\n\nONE TO MANY (.Parents)\n")
			fmt.Printf("-------------------------------------------\n\n")
			printParentTable(e.Name, e.ParentRelations)
		} else {
			ents, _ := input.Entities("")

			if ents == nil {
				return
			}

			var tbl table.Table

			if skipDescription {
				tbl = table.New("Name", "Schema", "Alias", "Rows")
				for _, c := range *ents {
					tbl.AddRow(c.Name, c.Schema, c.Alias, c.RowCount)
				}
			} else {
				tbl = table.New("Name", "Schema", "Alias", "Rows", "Description")
				for _, c := range *ents {
					tbl.AddRow(c.Name, c.Schema, c.Alias, c.RowCount, c.Description)
				}

			}

			tbl.Print()
		}

	},
}

func init() {
	rootCmd.AddCommand(entityCmd)

	entityCmd.Flags().BoolVarP(&skipDescription, "skip-description", "", false, "Does not show description")
}

func printChildTable(owner string, relations []input.Relation) {
	var childTable table.Table
	childTable = table.New("Schema", "Name", "FK", "PK", "Constraint")
	for _, ct := range relations {
		fn, pn := "", ""
		if ct.ColumnNullable {
			fn = "*"
		}
		ft := fmt.Sprintf("%s (%v%s)", ct.ColumnName, ct.ColumnType, fn)

		if ct.OwnerColumnNullable {
			pn = "*"
		}

		pt := fmt.Sprintf("%s (%v%s)", ct.OwnerColumnName, ct.OwnerColumnType, pn)
		childTable.AddRow(ct.Schema, ct.Name, ft, pt, ct.ContraintName)
	}

	childTable.Print()
}

func printParentTable(owner string, relations []input.Relation) {
	var tbl table.Table
	tbl = table.New("Schema", "Name", "PK", "FK", "Constraint")
	for _, ct := range relations {
		fn, pn := "", ""
		if ct.ColumnNullable {
			fn = "*"
		}
		ft := fmt.Sprintf("%s (%v%s)", ct.ColumnName, ct.ColumnType, fn)

		if ct.OwnerColumnNullable {
			pn = "*"
		}

		pt := fmt.Sprintf("%s (%v%s)", ct.OwnerColumnName, ct.OwnerColumnType, pn)
		tbl.AddRow(ct.Schema, ct.Name, ft, pt, ct.ContraintName)
	}

	tbl.Print()
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
