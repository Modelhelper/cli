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
	"modelhelper/cli/source"

	"github.com/gookit/color"

	"github.com/spf13/cobra"
)

// setConnectionCmd represents the setConnection command
var setConnectionCmd = &cobra.Command{
	Use:   "connection [key]",
	Short: "Adds or updates a connection",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		key := args[0]
		conType, _ := cmd.Flags().GetString("type")
		conString, _ := cmd.Flags().GetString("constr")
		conDesc, _ := cmd.Flags().GetString("description")
		conSchema, _ := cmd.Flags().GetString("schema")
		merge, _ := cmd.Flags().GetBool("merge")
		def, _ := cmd.Flags().GetBool("default")

		if len(conType) == 0 {
			conType = "mssql"
		}

		if !source.IsConnectionTypeValid(conType) {
			log.Fatalf("The type: %s is not a valid type\n", conType)
		}

		if len(conString) == 0 {

			color.Red.Println("NB !!")

			fmt.Printf("\nThe --constr | -c option is empty. Consider to use a valid connection string for %s\n", conType)
			fmt.Printf("Use mh build constr %s to build a valid ConnectionString for a %[1]s\n", conType)
			// fmt.Println("Copy the generated connection string and copy and paste in this")
			fmt.Println("\nUse the option --key <keynam> to update the connection with the new connection string")
		}

		c := source.Connection{
			Name:             key,
			Description:      conDesc,
			Schema:           conSchema,
			Type:             conType,
			ConnectionString: conString,
		}

		err := config.SetConnection(key, &c, def, merge)

		if err != nil {
			log.Fatalln("Could not add or update connection ", err)
		}

		fmt.Printf("Succesfully updated connection list with %s (%s)", key, conType)

	},
}

func init() {
	setCmd.AddCommand(setConnectionCmd)

	setConnectionCmd.Flags().StringP("type", "t", "mssql", "The type of connection to add, default mssql")
	setConnectionCmd.Flags().StringP("constr", "c", "", "Sets the connection string for the type. For ms sql use this format \n'sqlserver://<user>:<password>@<server>?database=<databasename>'")
	setConnectionCmd.Flags().StringP("description", "d", "", "Sets a description")
	setConnectionCmd.Flags().String("schema", "", "Sets the schema or owner for the collection of entities")
	// setConnectionCmd.Flags().BoolP("groups", "d", false, "If true, the cli will ask for groups to be added")
	// setConnectionCmd.Flags().Bool("merge", false, "If true and the connection exists, empty properties will be replaced with the incoming.")
	setConnectionCmd.Flags().Bool("default", false, "If true set this connection to the default connection (default value is false)")

	setConnectionCmd.MarkFlagRequired("type")
	setConnectionCmd.MarkFlagRequired("constr")
}
