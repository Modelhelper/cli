package connection

import (
	"io/ioutil"
	"log"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type connectionListService struct {
	cfg *models.Config
}

// BaseConnection implements modelhelper.ConnectionService
func (c *connectionListService) BaseConnection(name string) (*models.ConnectionList, error) {
	cons, err := c.Connections()

	if err != nil {
		return nil, err
	}

	con, f := cons[name]
	if !f {
		return nil, nil
	}

	return con, nil
}

// Connection implements modelhelper.ConnectionService
func (c *connectionListService) Connection(name string) (any, error) {
	cons, err := c.Connections()

	if err != nil {
		return nil, err
	}

	if cons != nil {
		item, found := cons[name]

		if found {
			switch item.Type {
			case "mssql":
				return loadGenericConnection[models.MsSqlConnection](item.Path)
			case "postgres":
				return loadGenericConnection[models.PostgresConnection](item.Path)
			case "file":
				return loadGenericConnection[models.FileConnection](item.Path)
			}

		}
	}
	return nil, nil
}

// Connections implements modelhelper.ConnectionService
func (c *connectionListService) Connections() (map[string]*models.ConnectionList, error) {
	defcon := c.cfg.DefaultConnection

	fileMap := make(map[string]*models.ConnectionList)

	conDir := filepath.Join(c.cfg.DirectoryName, "connections")
	filepath.Walk(conDir, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(fullPath, "yaml") || strings.HasSuffix(fullPath, "yml")) {
			cl, err := loadConnectionListFromFile(fullPath)
			if err != nil {
				return err
			}

			cl.IsDefault = cl.Name == defcon
			cl.Path = fullPath

			fileMap[cl.Name] = cl
		}

		return nil
	})

	return fileMap, nil

}

func NewConnectionService(cfg *models.Config) modelhelper.ConnectionService {
	return &connectionListService{cfg}
}

func loadConnectionListFromFile(fileName string) (*models.ConnectionList, error) {
	var t *models.ConnectionList

	dat, e := os.ReadFile(fileName)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
		return nil, e
	}

	err := yaml.Unmarshal(dat, &t)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil, err
	}

	if t != nil {
		t.Path = fileName
	}
	return t, nil
}

func loadGenericConnection[T models.GenericConnectionType](fileName string) (*models.GenericConnection[T], error) {
	t := &models.GenericConnection[T]{}

	dat, e := os.ReadFile(fileName)
	if e != nil {
		log.Fatalf("cannot load file: %v", e)
		return nil, e
	}

	err := yaml.Unmarshal(dat, &t)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
		return nil, err
	}

	if t != nil {
		t.Path = fileName
	}
	return t, nil
}
