package serve

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the modelhelper web server",
		Long:  `With this command you can spin up a server`,
		Run: func(cmd *cobra.Command, args []string) {
			// mhApp := app.New()
			// ctx := mhApp.CreateContext()
			// a := app.NewModelhelperCli()
			port, _ := cmd.Flags().GetInt("port")
			// open, _ := cmd.Flags().GetBool("open")
			message := fmt.Sprintf(`
	
	%s		
	
	ModelHelper website is now running on http://localhost:%v
	You may also access the ModelHelper API here: http://localhost:%v/api
	
	And read the API documentation here: http://localhost:%v/api/docs.
	
	To exit and stop the service, press ctrl + c`, "a.Logo()", port, port, port)

			fmt.Println(message)

			// server.Serve(port, open)
		},
	}
	serveCmd.Flags().IntP("port", "p", 8080, "The port to serve")
	serveCmd.Flags().BoolP("open", "o", false, "Opens a browser and the modelhelper website ")

	return serveCmd
}
