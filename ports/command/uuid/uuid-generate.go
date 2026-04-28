package uuid

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"strings"

	"github.com/spf13/cobra"

	gguuid "github.com/google/uuid"
)

func NewGenerateUUIDCommand(app *modelhelper.ModelhelperCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Short:   "Generate a new UUID using RFC 4122",
		Run:     uuidGenerateCommandHandler(app),
	}

	cmd.Flags().StringP("version", "v", "7", "The UUID version to generate (1, 3, 4, 5, 6, 7), default is 7")
	cmd.Flags().BoolP("uppercase", "u", false, "Outputs the UUID in uppercase, default is false")
	cmd.Flags().BoolP("braces", "b", false, "Outputs the UUID with braces, default is false")
	cmd.Flags().Bool("hyphens", true, "Outputs the UUID with hyphens, default is true")
	cmd.Flags().Int32P("count", "c", 1, "The number of UUIDs to generate, default is 1")

	return cmd
}

func uuidGenerateCommandHandler(app *modelhelper.ModelhelperCli) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		list := []string{}

		fmt.Println("Generating UUID...")
		version, _ := cmd.Flags().GetString("version")
		uppercase, _ := cmd.Flags().GetBool("uppercase")
		braces, _ := cmd.Flags().GetBool("braces")
		hyphens, _ := cmd.Flags().GetBool("hyphens")
		count, _ := cmd.Flags().GetInt32("count")

		if count < 1 {
			count = 1
		}

		for i := int32(0); i < count; i++ {
			uid, err := generateUUID(version)
			if err != nil {
				fmt.Printf("Error generating UUID: %v\n", err)
			}

			// fmt.Println("Generated UUID:", uid.String())

			// // Components of the uid
			// fmt.Printf("Version: %d\n", uid.Version())
			// fmt.Printf("Variant: %d\n", uid.Variant())
			// fmt.Printf("Timestamp: %d\n", uid.Time())
			// fmt.Printf("Clock Sequence: %d\n", uid.ClockSequence())

			id := formatUUID(uid, uppercase, braces, hyphens)
			list = append(list, id)
		}

		for _, id := range list {
			fmt.Println(id)
		}
	}
}

func formatUUID(uid gguuid.UUID, uppercase, braces, hyphens bool) string {
	id := uid.String()

	if uppercase {
		id = strings.ToUpper(id)
	}

	if !hyphens {
		id = strings.ReplaceAll(id, "-", "")
	}

	if braces {
		id = fmt.Sprintf("{%s}", id)
	}

	return id
}

func generateUUID(version string) (gguuid.UUID, error) {

	switch version {
	case "1":
		return gguuid.New(), nil
	case "3":
		return gguuid.NewSHA1(gguuid.NameSpaceURL, []byte("example.com")), nil
	case "4":
		return gguuid.NewRandom()
	case "5":
		return gguuid.NewSHA1(gguuid.NameSpaceURL, []byte("example.com")), nil
	case "6":
		return gguuid.NewV6()
	case "7":
	default:
		return gguuid.NewV7()
	}

	return gguuid.NewV7()
}
