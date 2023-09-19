package code

import (
	"regexp"
	"testing"
)

func testAdvancedMessage() string {
	return `fix(core): This should be the title

This is the body explaining more of the 
message

closes #V2-1234, #V2-4321
ref #V2-3212
#42312
BREAKING CHANGE 

This is the message
`
}

func TestBreakingChange(t *testing.T) {
	pat := "BREAKING CHANGE"
	rex := regexp.MustCompile(pat)

	is, body := checkForBreakingChange(testAdvancedMessage(), rex)
	bex := "This is the message"
	if is < 0 {
		t.Errorf("Expected %v but got %v", true, is)
	}

	if body != bex {
		t.Errorf("Expected %v but got %v", body, bex)

	}
}
func Test_That_We_Get_References(t *testing.T) {
	pat := `(.*?)?\#([0-9a-zA-Z-\\\.]*)`
	rex := regexp.MustCompile(pat)

	refs := getReferences(testAdvancedMessage(), rex)
	if len(refs) == 0 {
		t.Errorf("Expected %v but got %v", 1, len(refs))
	}

	bex := "V2-1234"
	if refs[0].Id != bex {
		t.Errorf("Expected %v but got %v", bex, refs[0].Id)
	}
}

func Test_Commit_history(t *testing.T) {
	repo := "C:\\dev\\projects\\mh\\cli"
	cs := &codeCommitService{}
	history, err := cs.GetCommitHistory(repo, nil)

	if err != nil {
		t.Log(err)
		t.Errorf("Expected get commit err to be nil")
	}
	if history == nil {
		t.Errorf("Expected history to not be nil")

	}

	if history.Authors != nil && len(history.Authors) != 1 {
		t.Errorf("Expected %v but got %v", 1, len(history.Authors))
	}
}
