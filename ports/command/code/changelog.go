package code

import (
	"fmt"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/ui"
	"modelhelper/cli/utils/path"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func NewChangelogCommand(app *modelhelper.ModelhelperCli) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "changelog",
		Short: "Reads the commit log and shows as a changelog",
		Long:  "",
		Run:   codeChangelogHandler(app),
	}

	cmd.Flags().StringP("repo", "r", "", "Where the repo is located")

	return cmd
}

type authorRenderer struct {
	rows map[string]models.Author
}

func (d *authorRenderer) Rows() [][]string {
	var rows [][]string

	for _, e := range d.rows {
		p := message.NewPrinter(language.English)

		r := []string{
			e.Name,
			p.Sprintf("%d", e.Commits),
			p.Sprintf("%v", e.First.Format("02-Jan-2006")),
			p.Sprintf("%v", e.Last.Format("02-Jan-2006")),
		}

		rows = append(rows, r)
	}

	return rows

}

func (d *authorRenderer) Header() []string {
	return []string{"Name", "Commits", "First", "Last"}
}

func codeChangelogHandler(app *modelhelper.ModelhelperCli) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		repo, _ := cmd.Flags().GetString("repo")
		if len(repo) == 0 {
			repo = path.CurrentDirectory()
		}
		commits, err := app.Code.CommitHistory.GetCommitHistory(repo, nil)

		if err != nil {
			// handle error
			fmt.Println(err)
			return
		}

		commitModel := convertToModel(commits)
		PrintCommits(commitModel)
	}
}

func PrintCommits(commits *models.CommitModel) {
	printCommitArray("Features", commits.Features)
	printCommitArray("Fixes", commits.Fixes)
	printCommitArray("Refactors", commits.Refactors)

	printCommitArray("BREAKING CHANGES", commits.BreakingChanges)

	fmt.Printf("\nAuthors\n\n")
	table := &authorRenderer{
		rows: commits.Authors,
	}

	ui.RenderTable(table)
	// for au, cnt := range commits.Authors {
	// 	fmt.Printf("%-20s: commits: %5d, first: %10v, last: %10v\n", au, cnt.Commits, cnt.First.Format("02-Jan-2006"), cnt.Last.Format("02-Jan-2006"))
	// }
}

func printCommitArray(title string, list []models.Commit) {

	if len(list) == 0 {
		return
	}

	fmt.Printf("%s\n\n", title)
	for _, cmt := range list {
		scope := ""
		if len(cmt.Scope) > 0 {
			scope = fmt.Sprintf("(%s): ", cmt.Scope)
		}
		fmt.Printf("\t* %s%s\n", scope, cmt.Title)
	}
}

func convertToModel(in *models.CommitHistory) *models.CommitModel {
	mdl := &models.CommitModel{}

	mdl.Features = in.Messages["feat"]
	mdl.Refactors = in.Messages["refactor"]
	mdl.Fixes = in.Messages["fixes"]

	for _, msg := range in.Messages {
		for _, commit := range msg {
			if commit.IsBreakingChange {
				mdl.BreakingChanges = append(mdl.BreakingChanges, commit)
			}
		}
	}

	mdl.Authors = in.Authors
	return mdl
}
