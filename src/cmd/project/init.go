package project

import (
	"fmt"
	"modelhelper/cli/config"
	"modelhelper/cli/project"
	"modelhelper/cli/ui"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func NewProjectInitCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new project in the current working directory",

		Run: func(cmd *cobra.Command, args []string) {

			p := project.NewModelhelperProject()

			init := true

			if project.Exists(project.DefaultLocation()) {
				color.Red.Println("NB!!")
				color.Red.Println("A project already exists in this location")
				init = ui.PromptForYesNo("Overwrite current project file? [y/N]", "n")
			}

			if init {
				p.Config.Version = "3.0"
				p.Config.Name = ui.PromptForString("Enter the name of the project")
				// p.DefaultSource = promptForConnectionKey()
				p.Config.Language = ui.PromptForLanguage("Select the primary code language")
				p.Config.Options = make(map[string]string)
				// p.OwnerName = promptForString("Enter the owner (company name) for this project")

				if ui.PromptForYesNo("Clone connections from config? [Y/n]", "y") {
					cfg := config.Load()
					p.Config.Connections = cfg.Connections
					// clone
				}

				err := p.Save()

				if err != nil {
					fmt.Println(err)
				}

				open, _ := cmd.Flags().GetBool("open")

				if project.Exists(project.DefaultLocation()) {

					if open {
						// openProjectInEditor()
					}
				}
			}
		},
	}

	cmd.Flags().Bool("open", false, "Opens the project file in default editor")

	return cmd
}
