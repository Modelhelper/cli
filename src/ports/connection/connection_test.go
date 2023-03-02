package connection

import (
	"modelhelper/cli/modelhelper/models"
	"testing"
)

func TestConnectionList(t *testing.T) {
	cfg := getDefaultConfig()
	consrv := NewConnectionService(cfg)

	clist, err := consrv.Connections()
	if err != nil {
		t.Errorf("Could not list connections")
	}

	expected := 1
	actual := len(clist)

	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}

	con, found := clist["stages"]

	if !found {
		t.Errorf("Expected %v got %v", "stages to exist", "stages does not exists")
	}

	if con.IsDefault != true {
		t.Errorf("Expected %v got %v", "stages to be default", "stages is not default")
	}

}

func Test_MsSQlConnection(t *testing.T) {
	cfg := getDefaultConfig()
	consrv := NewConnectionService(cfg)

	c, err := consrv.Connection("stages")

	if err != nil {
		t.Errorf("Could not list connections")
	}

	con, ok := c.(*models.GenericConnection[models.MsSqlConnection])

	if !ok {
		t.Errorf("Expected %v got %v", "models.MsSqlConnection", "something else")

	}

	if con == nil {
		t.Errorf("Expected %v got %v", "not nil", "nil")

	}

}

func getDefaultConfig() *models.Config {
	cfg := &models.Config{
		DirectoryName:     "C:\\Users\\Hans-PetterEitvet\\.modelhelper",
		DefaultConnection: "stages",
	}
	return cfg
}
