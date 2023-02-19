package source

import (
	"embed"
	"fmt"
	"log"
	"modelhelper/cli/modelhelper"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed entities/*
var entities embed.FS

type DemoSource struct{}

func (server *DemoSource) Entity(name string) (*modelhelper.Entity, error) {

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
func (server *DemoSource) Entities(pattern string) (*[]modelhelper.Entity, error) {
	list := []modelhelper.Entity{}
	baseList := getDemoEntities()

	for _, fent := range baseList {

		entity := fent.toSourceEntity()
		entity.ParentRelations = fent.getParentRelations(baseList)
		entity.ParentRelationCount = len(entity.ParentRelations)

		entity.ChildRelations = fent.getChildRelations(baseList)
		entity.ChildRelationCount = len(entity.ChildRelations)

		entity.ColumnCount = len(entity.Columns)
		entity.Alias = Abbreviate(entity.Name)
		list = append(list, entity)

	}

	return &list, nil

}
func (server *DemoSource) EntitiesFromColumn(column string) (*[]modelhelper.Entity, error) {
	list := []modelhelper.Entity{}
	baseList := getDemoEntities()

	for _, fent := range baseList {

		entity := fent.toSourceEntity()
		entity.ParentRelations = fent.getParentRelations(baseList)
		entity.ParentRelationCount = len(entity.ParentRelations)

		entity.ChildRelations = fent.getChildRelations(baseList)
		entity.ChildRelationCount = len(entity.ChildRelations)

		entity.ColumnCount = len(entity.Columns)
		entity.Alias = Abbreviate(entity.Name)
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

func (server *DemoSource) ConnectionStringPart(part string) string {
	return ""
}
func (server *DemoSource) ParseConnectionString() {
}
