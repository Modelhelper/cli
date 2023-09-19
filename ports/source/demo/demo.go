package demo

import (
	"embed"
	"fmt"
	"log"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/utils/text"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed entities/*
var entities embed.FS

type demoSource struct {
	// connectionService modelhelper.ConnectionService
}

func NewDemoSource() modelhelper.SourceService {
	src := &demoSource{
		// connectionService: cs,
	}

	// src.database = loadConnection(cs, connectionName)
	return src
}

func (server *demoSource) Entity(name string) (*models.Entity, error) {

	entityFiles, err := server.Entities("")
	if err != nil {
		fmt.Errorf("%v", err)
	}
	//var entity Entity

	for _, e := range *entityFiles {
		if strings.EqualFold(e.Name, name) {
			return &e, nil
		}
	}

	return nil, nil
}

func (server *demoSource) EntitiesFromNames(names []string) (*[]models.Entity, error) {
	all, err := server.Entities("")

	if err != nil {
		return nil, err
	}

	list := []models.Entity{}

	for _, name := range names {
		for _, entity := range *all {
			if entity.Name == name {
				list = append(list, entity)
			}
		}
	}
	return &list, nil
}

func (server *demoSource) Entities(pattern string) (*[]models.Entity, error) {
	list := []models.Entity{}
	baseList := getDemoEntities()

	for _, fent := range baseList {

		entity := fent.toSourceEntity()
		entity.ParentRelations = fent.getParentRelations(baseList)
		entity.ParentRelationCount = len(entity.ParentRelations)

		entity.ChildRelations = fent.getChildRelations(baseList)
		entity.ChildRelationCount = len(entity.ChildRelations)

		entity.ColumnCount = len(entity.Columns)
		entity.Alias = text.Abbreviate(entity.Name)
		list = append(list, entity)

	}

	return &list, nil

}
func (server *demoSource) EntitiesFromColumn(column string) (*[]models.Entity, error) {
	list := []models.Entity{}
	baseList := getDemoEntities()

	for _, fent := range baseList {

		entity := fent.toSourceEntity()
		entity.ParentRelations = fent.getParentRelations(baseList)
		entity.ParentRelationCount = len(entity.ParentRelations)

		entity.ChildRelations = fent.getChildRelations(baseList)
		entity.ChildRelationCount = len(entity.ChildRelations)

		entity.ColumnCount = len(entity.Columns)
		entity.Alias = text.Abbreviate(entity.Name)
		list = append(list, entity)

	}

	return &list, nil

}

func getDemoEntities() []fileEntity {
	list := []fileEntity{}
	root := "entities"
	files, err := entities.ReadDir(root)

	if err != nil {
		return nil
	}

	for _, file := range files {
		if !file.IsDir() {
			fn := file.Name()
			fullPath := path.Join(root, fn)
			blob, err := entities.ReadFile(fullPath)

			if err != nil {
				return nil
			}

			var fent fileEntity
			err = yaml.Unmarshal(blob, &fent)
			if err != nil {
				log.Fatalf("cannot unmarshal data: %v", err)
			}
			list = append(list, fent)
		}
	}

	return list
}

func (server *demoSource) ConnectionStringPart(part string) string {
	return ""
}
func (server *demoSource) ParseConnectionString() {
}
