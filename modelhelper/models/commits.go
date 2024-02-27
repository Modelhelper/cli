package models

import "time"

type Author struct {
	Commits int
	First   time.Time
	Last    time.Time
	Name    string
}

type Reference struct {
	Name string
	Id   string
}
type Commit struct {
	Message            string
	Type               string
	Scope              string
	Title              string
	When               time.Time
	Author             string
	Hash               string
	Body               string
	IsBreakingChange   bool
	BreakingChangeBody string
	References         []Reference
}
type CommitHistory struct {
	Name     string
	Messages map[string][]Commit
	Authors  map[string]Author
	Tags     []string
	From     time.Time
	To       time.Time
	Repo     string
}

type GitTag struct {
	Name    string
	Message string
	When    time.Time
	Hash    string
}
type CommitModel struct {
	Name                      string
	RootNamespace             string
	Namespace                 string
	Postfix                   string
	Prefix                    string
	ModuleLevelVariablePrefix string
	Inject                    []InjectSection
	Imports                   []string
	Project                   ProjectSection
	Feature                   FeatureModel
	Developer                 DeveloperSection
	Options                   map[string]string
	PageHeader                string
	IsConventialCommits       bool
	HasFeatures               bool
	Features                  []Commit
	HasRefactors              bool
	Refactors                 []Commit
	HasFixes                  bool
	Fixes                     []Commit
	Docs                      []Commit
	Performance               []Commit
	Tests                     []Commit
	Builds                    []Commit
	Ci                        []Commit
	Chores                    []Commit
	Reverts                   []Commit
	HasBreakingChanges        bool
	BreakingChanges           []Commit
	HasAuthors                bool
	Authors                   map[string]Author
}

type CommitHistoryOptions struct {
	StartFrom *string // date | tag | latest tag
	FromTag   *GitTag
	From      *time.Time
	Until     *time.Time
}
