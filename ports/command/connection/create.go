package connection

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

type createConnectionHandler struct {
	connectionService modelhelper.ConnectionService
	// configService     modelhelper.ConfigService
	// config            *models.Config
	// program           *tea.Program
}

func NewCreateConnectionCommand(cs modelhelper.ConnectionService) *cobra.Command {
	handler := &createConnectionHandler{
		connectionService: cs,
	}

	cmd := &cobra.Command{
		Example: "mh connection create mssql 'name'",
		Use:     "create [mssql|file] {name}",
		Aliases: []string{"c"},
		RunE:    handler.handleCmd,

		Args:  cobra.MaximumNArgs(1),
		Short: "Create a connections, use sub commands for each specific type of connection",
		// 		Long: `
		// This is the root command for creating a connection. You must combine this
		// command with a sub command of either 'mssql' or 'file', followed by the name of the connection

		// Use the various options for each connection type to give details.`,
		// Run:
	}

	cmd.Flags().String("clone", "", "Clone the connection from the given name")

	return cmd
}

func (h *createConnectionHandler) handleCmd(cmd *cobra.Command, args []string) error {
	var (
		dbType    string
		csType    string
		conName   string
		conString string
		desc      string
		save      bool
	)

	if len(args) > 0 {
		conName = args[0]
	}

	// if cmd.Flags().GetBool("clone") {

	// }

	nameIsNotUsed := func(name string) error {
		if len(name) == 0 {
			return fmt.Errorf("The name cannot be empty")
		}

		if strings.Contains(name, " ") {
			return fmt.Errorf("The name cannot contain spaces")
		}

		cons, _ := h.connectionService.Connections()
		if _, exists := cons[name]; exists {
			return fmt.Errorf("A connection with the name already exists")
		}

		return nil
	}

	dbTypeSelect := huh.NewSelect[string]().
		Title("Select Database Engine").
		Options(
			huh.NewOption("Postgres", "postgres"),
			huh.NewOption("Microsoft SQL Server", "mssql"),
			huh.NewOption("Filebased database description in YAML", "file"),
		).
		Value(&dbType)
	connectionStringSelect := huh.NewSelect[string]().
		Title("Select Database Engine").
		Options(
			huh.NewOption("Enter full connection string", "full"),
			huh.NewOption("Select enviornment variable", "env"),
			huh.NewOption("Build and use as environment var", "buildEnv"),
			huh.NewOption("Build the connection string", "build"),
		).
		Value(&csType)

	descInput := huh.NewText().
		Title("Description").
		Placeholder("What is this description for").
		Value(&desc)

	connNameInput := huh.NewInput().
		Title("Enter connection name").
		Prompt("name: ").
		Validate(nameIsNotUsed).
		Value(&conName)

	confirm := huh.NewConfirm().
		Title("Save connection details?").
		Description("Thats all we need for now ").
		Affirmative("Save").
		Negative("Cancel").
		Value(&save)

	huh.NewForm(
		huh.NewGroup(connNameInput, dbTypeSelect),
		huh.NewGroup(connectionStringSelect),
	).Run()

	switch csType {
	case "full":
		conString = runConnectionStringEnterFull(dbType)
	case "env":
		conString = runConnectionStringFromEnv(dbType)
	case "buildEnv":
		conString = runConnectionStringBuilder(dbType)
	case "build":
		conString = runConnectionStringBuilder(dbType)
	}

	huh.NewForm(
		huh.NewGroup(descInput, confirm),
	).Run()

	if save {
		err := h.connectionService.Create(&models.ConnectionList{
			Name:             conName,
			Description:      desc,
			Type:             dbType,
			ConnectionString: conString,
		})
		if err != nil {
			return err
		}
	}

	fmt.Println("Selected: ", dbType, conName, conString, desc, save)
	return nil
}

func runConnectionStringBuilder(dbType string) string {
	csFormat := ""

	if dbType == "mssql" {
		csFormat = "sqlserver://%s:%s@%s?database=%s"
	} else if dbType == "postgres" {
		csFormat = "postgresql://%s:%s@%s?database=%s"
	}
	// csFormat := "postgres://%s:%s@%s/%s?sslmode=%s"
	var (
		user     string
		password string
		hostname string
		database string
		// port     int
	)

	hostInput := huh.NewInput().
		Title("Enter connection string details").
		Prompt("Host name: ").
		// PlaceholderFunc(func() string {
		// 	if dbType == "mssql" {
		// 		return "sqlserver://user:password@localhost?database=dbname"
		// 	} else if dbType == "postgres" {
		// 		return "postgresql://user:password@localhost?database=dbname"
		// 	}
		// 	return "Enter the connection string"
		// }, nil).
		// Validate(nameIsNotUsed).
		Value(&hostname)

	userInput := huh.NewInput().
		Prompt("User name: ").
		Value(&user)
	pwdInput := huh.NewInput().
		Prompt("Password: ").
		EchoMode(huh.EchoModePassword).
		// Password(true).
		Value(&password)

	databaseInput := huh.NewInput().
		Prompt("Database name: ").
		Value(&database)

	huh.NewForm(
		huh.NewGroup(hostInput, userInput, pwdInput, databaseInput),
	// huh.NewGroup(descInput, confirm),
	).Run()

	return fmt.Sprintf(csFormat, user, password, hostname, database)
}

func runConnectionStringEnterFull(dbType string) string {
	csFormat := ""

	if dbType == "mssql" {
		csFormat = "sqlserver://%s:%s@%s?database=%s"
	} else if dbType == "postgres" {
		csFormat = "postgresql://%s:%s@%s?database=%s"
	}
	// csFormat := "postgres://%s:%s@%s/%s?sslmode=%s"
	var (
		constr string
		// port     int
	)

	hostInput := huh.NewInput().
		Title("Enter connection string").
		Prompt("Host name: ").
		PlaceholderFunc(func() string {
			return fmt.Sprintf(csFormat, "user", "password", "localhost", "dbname")
		}, nil).
		// Validate(nameIsNotUsed).
		Value(&constr)

	huh.NewForm(
		huh.NewGroup(hostInput),
	).Run()

	return constr
}
func runConnectionStringFromEnv(dbType string) string {
	csFormat := "%s:%s"
	var envVar string
	opts := []huh.Option[string]{}
	sort.Strings(os.Environ())

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		valLen := len(pair[1])
		if valLen > 20 {
			valLen = 20
		}
		ho := huh.NewOption(fmt.Sprintf("%-*s [%s...]", 25, pair[0], pair[1][:valLen]), pair[0])

		opts = append(opts, ho)

	}

	envSelect := huh.NewSelect[string]().
		Title("Select an environment variable to use for connection string").
		Options(opts...).
		Value(&envVar)

	huh.NewForm(
		huh.NewGroup(envSelect),
	).Run()

	return fmt.Sprintf(csFormat, "ENV", envVar)
}

func NewCreateMsSQLConnectionCommand(cs modelhelper.ConnectionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mssql",
		Args:  cobra.MinimumNArgs(1),
		Short: "Creates a MS SQL connection",
		Run: func(cmd *cobra.Command, args []string) {
			cons, _ := cs.Connections()
			name := args[0]

			_, exists := cons[name]

			if exists {
				panic("A connection with the name already exists")
			}

			con := &models.MsSqlConnection{}

			cs, _ := cmd.Flags().GetString("connection-string")

			if len(cs) > 0 {
				con.ConnectionString = cs

			} else {

				s, _ := cmd.Flags().GetString("server")
				d, _ := cmd.Flags().GetString("database")
				u, _ := cmd.Flags().GetString("user")
				p, _ := cmd.Flags().GetString("password")

				if len(s) > 0 && len(d) > 0 && len(u) > 0 && len(p) > 0 {
					con.Database = d
					con.Server = s
					con.ConnectionString = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", u, p, s, d)
				}
			}
		},
	}

	cmd.Flags().StringP("connection-string", "c", "", "The full connection string to use")
	cmd.Flags().StringP("server", "s", "", "The server")
	cmd.Flags().StringP("database", "d", "", "The database")
	cmd.Flags().StringP("username", "u", "", "The user")
	cmd.Flags().StringP("password", "p", "", "The password")
	return cmd
}

func NewCreateFileConnectionCommand(cs modelhelper.ConnectionService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Args:  cobra.MinimumNArgs(1),
		Short: "Creates a file connection",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Create a file connection")
		},
	}

	return cmd
}
