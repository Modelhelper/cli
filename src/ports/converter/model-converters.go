package converter

import (
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"strings"
)

type codeModelConverter struct {
	// app *modelhelper.ModelhelperCli
}

// ToBasicModel implements modelhelper.CodeModelConverter
func (c *codeModelConverter) ToBasicModel(identifier, language string, project *models.ProjectConfig) *models.BasicModel {
	b := models.BasicModel{}
	imports := []string{}

	if project == nil {
		project = emptyProject()
	}

	// inject := map[]
	code, codeFound := project.Code[language]

	if len(project.Options) > 0 {
		b.Options = project.Options
	}
	b.Project.Name = project.Name
	b.Project.Owner = project.OwnerName

	b.PageHeader = project.Header

	if len(identifier) > 0 && codeFound {
		val, found := code.Keys[identifier]
		if found {
			b.RootNamespace = code.RootNamespace
			imports = append(imports, val.Imports...)

			b.Inject = []models.InjectSection{}
			for _, injectKey := range val.Inject {
				injItem, foundInj := code.Inject[injectKey]
				if foundInj {
					b.Inject = append(b.Inject, toInjectSection(injItem, b))
				}
			}

			b.Postfix = val.Postfix
			b.Prefix = val.Prefix
			b.Namespace = val.Namespace

		}
	}

	return &b
}

func emptyProject() *models.ProjectConfig {
	p := &models.ProjectConfig{
		Name:        "",
		Version:     "",
		DefaultKey:  "",
		Options:     make(map[string]string),
		Language:    "",
		Header:      "",
		Custom:      nil,
		Description: "",
		Code:        make(map[string]models.Code),
		OwnerName:   "",
	}

	return p
}

// ToEntityListModel implements modelhelper.CodeModelConverter
func (c *codeModelConverter) ToEntityListModel(identifier, language string, project *models.ProjectConfig, entities *[]models.Entity) *models.EntityListModel {
	entitylist := []models.EntityModel{}

	base := c.ToBasicModel(identifier, language, project)
	if entities != nil && len(*entities) > 0 {
		for _, entity := range *entities {

			entityBase := toEntitySection(&entity)
			entitylist = append(entitylist, entityBase)
		}
	}

	out := models.EntityListModel{
		Namespace:                 base.Namespace,
		Postfix:                   base.Postfix,
		Prefix:                    base.Prefix,
		ModuleLevelVariablePrefix: base.ModuleLevelVariablePrefix,
		Inject:                    base.Inject,
		Imports:                   base.Imports,
		Project:                   base.Project,
		Developer:                 base.Developer,
		Options:                   base.Options,
		PageHeader:                base.PageHeader,
		Entities:                  entitylist,
	}

	return &out
}

// ToEntityModel implements modelhelper.CodeModelConverter
func (c *codeModelConverter) ToEntityModel(key, language string, project *models.ProjectConfig, entity *models.Entity) *models.EntityModel {
	entityBase := models.EntityModel{}

	base := c.ToBasicModel(key, language, project)
	if entity != nil {
		entityBase = toEntitySection(entity)
	}

	out := models.EntityModel{
		RootNamespace:             base.RootNamespace,
		Namespace:                 base.Namespace,
		Postfix:                   base.Postfix,
		Prefix:                    base.Prefix,
		ModuleLevelVariablePrefix: base.ModuleLevelVariablePrefix,
		Inject:                    base.Inject,
		Imports:                   base.Imports,
		Project:                   base.Project,
		Developer:                 base.Developer,
		Options:                   base.Options,
		PageHeader:                base.PageHeader,
		Name:                      entityBase.Name,
		Schema:                    entityBase.Schema,
		Type:                      entityBase.Type,
		Alias:                     entityBase.Alias,
		Description:               entityBase.Description,
		HasDescription:            len(entityBase.Description) > 0,
		HasPrefix:                 false, //len(entityBase.Prefix) > 0,
		NameWithoutPrefix:         "",
		Columns:                   entityBase.Columns,
		Parents:                   entityBase.Parents,
		Children:                  entityBase.Children,
		PrimaryKeys:               entityBase.PrimaryKeys,
		ForeignKeys:               entityBase.ForeignKeys,
		UsedAsColumns:             entityBase.UsedAsColumns,
		UsesIdentityColumn:        entityBase.UsesIdentityColumn,
		HasSynonym:                entityBase.HasSynonym,
		Synonym:                   entityBase.Synonym,
		ModelName:                 entityBase.ModelName,
	}

	out.HasChildren = len(entityBase.Children) > 0
	out.HasParents = len(entityBase.Parents) > 0
	return &out
}

func NewCodeModelConverter() modelhelper.CodeModelConverter {
	return &codeModelConverter{}
}

func toInjectSection(from models.Inject, m interface{}) models.InjectSection {
	// name, _ := codegen.Generate("fileName", from.Name, m)
	code := models.InjectSection{
		Name:         from.Name,
		PropertyName: from.PropertyName,
	}

	return code
}

func toEntitySection(from *models.Entity) models.EntityModel {
	out := models.EntityModel{
		Name:               from.Name,
		Schema:             from.Schema,
		Type:               from.Type,
		Alias:              from.Alias,
		Description:        from.Description,
		HasDescription:     len(from.Description) > 0,
		HasPrefix:          false,
		NameWithoutPrefix:  "",
		UsesIdentityColumn: from.UsesIdentityColumn,
		HasSynonym:         from.HasSynonym,
	}

	out.Synonym = from.Synonym

	out.ModelName = coalesceString(from.Synonym, from.Name)

	for _, column := range from.Columns {

		col := toColumnSection(column, out.Name)

		out.Columns = append(out.Columns, col)

		if column.IsPrimaryKey {
			out.PrimaryKeys = append(out.PrimaryKeys, col)
		}

		if column.IsForeignKey {
			out.ForeignKeys = append(out.ForeignKeys, col)
		}
	}

	for _, cr := range from.ChildRelations {
		child := models.EntityRelationModel{}

		child.Name = cr.Name
		child.Schema = cr.Schema
		child.RelatedColumn = models.EntityColumnProps{
			Name:       cr.ColumnName,
			DataType:   cr.ColumnType,
			IsNullable: cr.ColumnNullable,
		}

		child.OwnerColumn = models.EntityColumnProps{
			Name:       cr.OwnerColumnName,
			DataType:   cr.OwnerColumnType,
			IsNullable: cr.OwnerColumnNullable,
		}

		out.Children = append(out.Children, child)
		// children = append(children)
		child.NameWithoutPrefix = strings.TrimPrefix(cr.Name, out.Name)
		child.HasPrefix = strings.HasPrefix(cr.Name, out.Name)

		child.HasSynonym = cr.HasSynonym
		if child.HasSynonym {
			child.Synonym = cr.Synonym
		}

		child.ModelName = coalesceString(cr.Synonym, cr.Name)
	}

	for _, pr := range from.ParentRelations {
		parent := models.EntityRelationModel{}
		parent.Name = pr.Name
		parent.Schema = pr.Schema

		parent.HasDescription = false

		parent.HasSynonym = pr.HasSynonym
		if parent.HasSynonym {
			parent.Synonym = pr.Synonym
		}

		parent.OwnerColumn = models.EntityColumnProps{
			Name:       pr.ColumnName,
			DataType:   pr.ColumnType,
			IsNullable: pr.ColumnNullable,
		}

		parent.RelatedColumn = models.EntityColumnProps{
			Name:       pr.OwnerColumnName,
			DataType:   pr.OwnerColumnType,
			IsNullable: pr.OwnerColumnNullable,
		}

		parent.ModelName = coalesceString(pr.Synonym, pr.Name)

		parent.NameWithoutPrefix = strings.TrimPrefix(pr.Name, out.Name)
		parent.HasPrefix = strings.HasPrefix(pr.Name, out.Name)

		out.Parents = append(out.Parents, parent)
	}

	return out
}

func coalesceString(name ...string) string {
	output := ""
	for _, n := range name {
		if len(n) > 0 {
			output = n
			break
		}
	}

	return output
}

func toProjectModel(input models.ProjectConfig) models.ProjectSection {
	return models.ProjectSection{
		Owner: input.OwnerName,
		Name:  input.Name,
	}

}

func toColumnSection(from models.Column, entityName string) models.EntityColumnModel {
	col := models.EntityColumnModel{
		Description:       from.Description,
		IsForeignKey:      from.IsForeignKey,
		IsPrimaryKey:      from.IsPrimaryKey,
		IsIdentity:        from.IsIdentity,
		IsNullable:        from.IsNullable,
		IsIgnored:         from.IsIgnored,
		IsDeletedMarker:   from.IsDeletedMarker,
		IsCreatedDate:     from.IsCreatedDate,
		IsCreatedByUser:   from.IsCreatedByUser,
		IsModifiedDate:    from.IsModifiedByUser,
		IsModifiedByUser:  from.IsModifiedByUser,
		HasPrefix:         strings.HasPrefix(from.Name, entityName),
		HasDescription:    len(from.Description) > 0,
		Name:              from.Name,
		NameWithoutPrefix: strings.TrimPrefix(from.Name, entityName),
		Collation:         from.Collation,
		ReferencesColumn:  from.ReferencesColumn,
		ReferencesTable:   from.ReferencesTable,
		DataType:          from.DataType,
		Length:            from.Length,
		Precision:         from.Precision,
		Scale:             from.Scale,
		UseLength:         from.UseLength,
		UsePrecision:      from.UsePrecision,
		ForCreate:         from.ForCreate,
	}

	return col
}
