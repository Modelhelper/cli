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
	"modelhelper/cli/config"
	"modelhelper/cli/project"
	"modelhelper/cli/vault"
	"path/filepath"

	"github.com/spf13/cobra"
)

// secretGetCmd represents the secretGet command
var secretGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in the vault",

	Run: func(cmd *cobra.Command, args []string) {

		secretsPath := filepath.Join(config.Location(), ".secrets")
		encodingKey, _ := cmd.Flags().GetString("key")
		parent := cmd.Parent()

		if parent != nil {
			// root
			if parent.Use == "project" {
				secretsPath = filepath.Join(project.DefaultLocation(), ".secrets")
			}

		}

		v := vault.File(encodingKey, secretsPath)
		key := args[0]
		value, err := v.Get(key)
		if err != nil {
			fmt.Println("no value set")
			return
		}
		fmt.Printf("%s = %s\n", key, value)
	},
}

func init() {
	secretCmd.AddCommand(secretGetCmd)

}
