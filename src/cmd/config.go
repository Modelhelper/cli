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
	"modelhelper/cli/config"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "<not implemented>",
	Run: func(cmd *cobra.Command, args []string) {

		open, _ := cmd.Flags().GetBool("open")
		if open {
			var editor string
			cfg := config.Load()

			if len(cfg.DefaultEditor) > 0 {
				editor = cfg.DefaultEditor
			} else {
				editor = promptForEditorKey("Please select editor to open the config")
			}

			loc := filepath.Join(config.Location(), "config.yaml")
			exe := exec.Command(editor, loc)
			if exe.Run() != nil {
				//vim didn't exit with status code 0
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().Bool("open", false, "Opens the config file in VS Code")
}
