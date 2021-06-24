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

// projectKeyCmd represents the projectKey command
var projectKeyCmd = &cobra.Command{
	Use:   "key [keyname]",
	Short: "Adds or updates a key for the coding section",
	// ValidArgsFunction: &cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("projectKey called")
	},
}

func init() {
	projectCmd.AddCommand(projectKeyCmd)

	projectKeyCmd.Flags().StringP("namespace", "n", "", "Defines the namespace for the key")
	projectKeyCmd.Flags().String("postfix", "", "Defines the postfix for the key")
	projectKeyCmd.Flags().String("prefix", "", "Defines the prefix for the key")
	projectKeyCmd.Flags().StringArray("import", []string{}, "Defines the prefix for the key")
	projectKeyCmd.Flags().StringArray("inject", []string{}, "Defines the prefix for the key")
}
