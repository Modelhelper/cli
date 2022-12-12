package app

import (
	"fmt"
	"modelhelper/cli/config"

	"github.com/gookit/color"
)

func Initialize(initializer Initializer) error {
	err := initializer.Initialize()
	return err
}

type Initializer interface {
	Initialize() error
}

func (a *Application) Initialize() error {

	PrintWelcomeMessage()

	do := PromptForContinue()

	if do {

		cfg := config.New()

		if err := cfg.Initialize(); err != nil {
			return err
		}

		fmt.Printf(`
Thank you. Your modelhelper config file has been created here:

	%s

%s

If you have questions about the use, history of modelhelper, templates, sources and so on, 
consider becoming a member of my Slack channel: https://model-helper.slack.com and/or 
add issues, bugs and so on in the GitHub repo: https://github.com/Modelhelper/cli

	`, color.Gray.Sprint(config.Location()), color.Green.Sprint("You are now ready to roll :-)"))
	}

	return nil
}
