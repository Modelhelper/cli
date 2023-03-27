package code

import (
	"fmt"
	"io/ioutil"
	"log"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"os"
	"regexp"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type codeCommitService struct {
}

// GetAuthors implements modelhelper.CommitHistoryService
func (*codeCommitService) GetAuthors(repoPath string, options *models.CommitHistoryOptions) (*models.Author, error) {
	panic("unimplemented")
}

// GetTags implements modelhelper.CommitHistoryService
func (*codeCommitService) GetTags(repoPath string, options *models.CommitHistoryOptions) ([]models.GitTag, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Error opening repository: %s", err)
	}

	list := []models.GitTag{}
	var tag *models.GitTag
	tagos, _ := repo.TagObjects()
	err = tagos.ForEach(func(t *object.Tag) error {
		tag.Hash = t.Hash.String()
		tag.When = t.Tagger.When
		tag.Name = t.Name
		tag.Message = t.Message

		list = append(list, *tag)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return list, nil
}

func conventionalCommitTypes() []string {
	return []string{
		"feat",
		"fix",
		"docs",
		"style",
		"refactor",
		"perf",
		"test",
		"build",
		"ci",
		"chore",
		"revert",
	}
}

// GetCommitHistory implements modelhelper.CommitHistoryService
func (cs *codeCommitService) GetCommitHistory(repoPath string, options *models.CommitHistoryOptions) (*models.CommitHistory, error) {

	cc := &models.CommitHistory{
		Repo:     repoPath,
		Messages: make(map[string][]models.Commit),
		Authors:  make(map[string]models.Author),
	}

	repoPath, _ = findClosestGitDir(repoPath)
	// Open the Git repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
		// log.Fatalf("Error opening repository: %s", err)
	}

	typeString := strings.Join(conventionalCommitTypes(), "|")
	regPat := fmt.Sprintf(`(%s)(\(.*\))?:(.*)`, typeString)
	regx := regexp.MustCompile(regPat)

	regxBreaking := regexp.MustCompile("BREAKING CHANGE")
	regxRefs := regexp.MustCompile(`(.*?)?\#([0-9a-zA-Z-\\\.]*)`)

	var tag *object.Tag
	tagos, _ := repo.TagObjects()
	_ = tagos.ForEach(func(t *object.Tag) error {
		tag = t
		return nil
	})
	gitopt := &git.LogOptions{
		Order: git.LogOrderCommitterTime,
	}

	if tag != nil {
		gitopt.Since = &tag.Tagger.When
	}
	// t := time.Now().AddDate(-1, 0, 0)
	// Get the commit history from the recent tag up to HEAD
	commitIter, err := repo.Log(gitopt)
	if err != nil {
		return nil, err

		// log.Fatalf("Error getting commit history: %v", err)
	}

	// fmt.Printf("Commits since %v\n\n", tag.Tagger.When)
	// Print the commit messages
	err = commitIter.ForEach(func(c *object.Commit) error {

		aut, ok := cc.Authors[c.Author.Name]
		if !ok {
			cc.Authors[c.Author.Name] = models.Author{
				Name:    c.Author.Name,
				First:   c.Committer.When,
				Last:    c.Committer.When,
				Commits: 1,
			}
		} else {
			if aut.Last.Before(c.Committer.When) {
				aut.Last = c.Committer.When
			}
			if aut.First.After(c.Committer.When) {
				aut.First = c.Committer.When
			}
			aut.Commits++
			cc.Authors[c.Author.Name] = aut
		}

		cmt := parseMessage(c, regx)
		bcStart, bcBody := checkForBreakingChange(c.Message, regxBreaking)
		refs := getReferences(c.Message, regxRefs)

		cmt.IsBreakingChange = bcStart > 0
		if bcStart > 0 && len(bcBody) > 0 {
			cmt.BreakingChangeBody = bcBody
		}

		cmt.References = refs

		cc.Messages[cmt.Type] = append(cc.Messages[cmt.Type], *cmt)

		return nil
	})
	if err != nil {
		return nil, err
		// log.Fatalf("Error iterating through commits: %s", err)
	}

	return cc, nil
}

func NewCommitHistoryService() modelhelper.CommitHistoryService {
	return &codeCommitService{}
}

func parseMessage(message *object.Commit, regx *regexp.Regexp) *models.Commit {
	msg := &models.Commit{}

	lines := strings.Split(message.Message, "\r\n")

	matches := regx.FindStringSubmatch(lines[0])

	if len(matches) > 0 {

		for idx, match := range matches {

			switch idx {
			case 0:
				msg.Message = match
			case 1:
				msg.Type = strings.TrimSpace(match)
			case 2:
				scope := strings.TrimSpace(match)
				scope = strings.TrimPrefix(scope, "(")
				scope = strings.TrimSuffix(scope, ")")
				msg.Scope = scope
			case 3:
				msg.Title = strings.TrimSpace(match)
			}

		}
	} else {
		msg.Title = lines[0]
	}

	msg.Author = message.Author.Name
	msg.When = message.Author.When
	msg.Hash = message.Hash.String()
	msg.Body = strings.Join(lines[1:], "\n")
	return msg
}

func getReferences(message string, regx *regexp.Regexp) []models.Reference {
	refs := []models.Reference{}
	matches := regx.FindAllStringSubmatch(message, -1)

	if len(matches) > 0 {

		for _, set := range matches {
			ref := models.Reference{}
			for idx, match := range set {

				switch idx {

				case 2:
					ref.Id = strings.TrimSpace(match)

				}
			}
			refs = append(refs, ref)

		}
	}

	return refs
}
func checkForBreakingChange(message string, regx *regexp.Regexp) (int, string) {
	matches := regx.FindStringIndex(message)

	startPos := -1
	body := ""

	if len(matches) > 0 {

		startPos = matches[0]
		pos := matches[0] + len("BREAKING CHANGE")
		body = message[pos:]

	}

	body = strings.TrimPrefix(body, "\n")
	body = strings.TrimSuffix(body, "\n")
	body = strings.TrimSpace(body)
	return startPos, body
}

func findClosestGitDir(startPath string) (string, bool) {

	folders := strings.Split(startPath, string(os.PathSeparator))

	for i := len(folders); i > 2; i-- {
		testPath := strings.Join(folders[:i], string(os.PathSeparator))
		files, err := ioutil.ReadDir(testPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {

			if f.IsDir() && f.Name() == ".git" {
				return testPath, true
			}
		}
	}

	return "", false
}
