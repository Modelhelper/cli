package converter

import (
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"strings"
)

type codeModelConverter struct {
	// app *modelhelper.ModelhelperCli
}

// ToCustomModel implements modelhelper.CodeModelConverter
func (c *codeModelConverter) ToCustomModel(key string, language string, project *models.ProjectConfig, custom any) *models.CustomModel {
	base := c.ToBasicModel(key, language, project)
	mdl := &models.CustomModel{
		RootNamespace:             base.RootNamespace,
		Namespace:                 base.Namespace,
		Postfix:                   base.Postfix,
		Prefix:                    base.Prefix,
		ModuleLevelVariablePrefix: base.ModuleLevelVariablePrefix,
		Inject:                    base.Inject,
		Imports:                   base.Imports,
		Project:                   base.Project,
		Feature:                   base.Feature,
		Developer:                 base.Developer,
		Options:                   base.Options,
		PageHeader:                base.PageHeader,
		Custom:                    custom,
	}

	return mdl
}

// ToNameModel implements modelhelper.CodeModelConverter
func (c *codeModelConverter) ToNameModel(key string, language string, project *models.ProjectConfig, name string) *models.NameModel {
	base := c.ToBasicModel(key, language, project)
	mdl := &models.NameModel{
		RootNamespace:             base.RootNamespace,
		Namespace:                 base.Namespace,
		Postfix:                   base.Postfix,
		Prefix:                    base.Prefix,
		ModuleLevelVariablePrefix: base.ModuleLevelVariablePrefix,
		Inject:                    base.Inject,
		Imports:                   base.Imports,
		Project:                   base.Project,
		Feature:                   base.Feature,
		Developer:                 base.Developer,
		Options:                   base.Options,
		PageHeader:                base.PageHeader,
		Name:                      name,
	}

	return mdl
}

// ToCommitHistoryModel implements modelhelper.CodeModelConverter
func (c *codeModelConverter) ToCommitHistoryModel(key string, language string, project *models.ProjectConfig, commitHistory *models.CommitHistory) *models.CommitModel {
	base := c.ToBasicModel(key, language, project)
	mdl := &models.CommitModel{
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
		Feature:                   base.Feature,
	}
	mdl.Name = commitHistory.Name
	mdl.Features = commitHistory.Messages["feat"]
	mdl.Fixes = commitHistory.Messages["fix"]
	mdl.Refactors = commitHistory.Messages["refactor"]
	mdl.Docs = commitHistory.Messages["docs"]
	mdl.Performance = commitHistory.Messages["perf"]
	mdl.Tests = commitHistory.Messages["tests"]
	mdl.Builds = commitHistory.Messages["builds"]
	mdl.Ci = commitHistory.Messages["ci"]
	mdl.Chores = commitHistory.Messages["chores"]
	mdl.Reverts = commitHistory.Messages["reverts"]

	mdl.HasFeatures = len(mdl.Features) > 0
	mdl.HasRefactors = len(mdl.Refactors) > 0
	mdl.HasFixes = len(mdl.Fixes) > 0

	for _, msg := range commitHistory.Messages {
		for _, commit := range msg {
			if commit.IsBreakingChange {
				mdl.BreakingChanges = append(mdl.BreakingChanges, commit)
			}
		}
	}
	mdl.HasBreakingChanges = len(mdl.BreakingChanges) > 0

	mdl.Authors = commitHistory.Authors
	mdl.HasAuthors = len(mdl.Authors) > 0
	return mdl

}

func (c *codeModelConverter) ToFeatureModel(project *models.ProjectConfig) (*models.FeatureModel, []string) {
	featureSet := models.FeatureModel{}
	imports := []string{}

	if project.Features != nil {

		// if project.Features.Logger != nil {
		// 	featureSet.UseLogger = project.Features.Logger.Use
		// 	imports = append(imports, project.Features.Logger.Imports...)
		// 	featureSet.Logger = models.FeatureOptions{
		// 		Namespace:    project.Features.Logger.Namespace,
		// 		PropertyName: *project.Features.Logger.PropertyName,
		// 		Type:         *project.Features.Logger.Type,
		// 	}
		// }

		// if project.Features.Api != nil {
		// 	featureSet.UseApi = project.Features.Api.Use
		// 	imports = append(imports, project.Features.Api.Imports...)

		// 	featureSet.Api = models.FeatureOptions{
		// 		Namespace:    project.Features.Api.Namespace,
		// 		PropertyName: *project.Features.Api.PropertyName,
		// 		Type:         *project.Features.Api.Type,
		// 	}
		// }

		// if project.Features.Db != nil {

		// 	featureSet.UseDb = project.Features.Db.Use
		// 	featureSet.Db = models.DbFeatureOptions{
		// 		Namespace:    project.Features.Db.Namespace,
		// 		PropertyName: *project.Features.Db.PropertyName,
		// 		Type:         *project.Features.Db.Type,
		// 	}
		// }
	}

	return &featureSet, imports
}

// ToBasicModel implements modelhelper.CodeModelConverter
func (c *codeModelConverter) ToBasicModel(identifier, language string, project *models.ProjectConfig) *models.BasicModel {
	b := models.BasicModel{}
	imports := []string{}

	if project == nil {
		project = emptyProject()

	}

	if len(project.Options) > 0 {
		b.Options = project.Options
	}

	b.Project.Exists = project != nil
	b.Project.Name = project.Name
	b.Project.Owner = project.OwnerName

	b.PageHeader = project.Header

	feat, featImports := c.ToFeatureModel(project)
	imports = append(imports, featImports...)
	b.Feature = *feat

	if len(project.RootNamespace) > 0 {
		b.RootNamespace = project.RootNamespace
	}

	if len(identifier) > 0 {
		val, found := project.Setup[identifier]
		if found {
			b.RootNamespace = project.RootNamespace
			imports = append(imports, val.Imports...)

			b.Inject = []models.InjectSection{}
			for _, injectKey := range val.Inject {
				injItem, foundInj := project.Inject[injectKey]
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
		Name:           "",
		Version:        "",
		DefaultKey:     "",
		Options:        make(map[string]string),
		Language:       "",
		Header:         "",
		Custom:         nil,
		Description:    "",
		Setup:          make(map[string]models.Key),
		Inject:         make(map[string]models.Inject),
		OwnerName:      "",
		Features:       nil,
		CustomFeatures: make(map[string]models.CommonProjectFeature),
		Locations:      make(map[string]string),
		Directory:      "",
		UseHeader:      false,
		RootNamespace:  "",
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
		Feature:                   base.Feature,
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
		Feature:                   base.Feature,
		Name:                      entityBase.Name,
		Schema:                    entityBase.Schema,
		Type:                      entityBase.Type,
		Alias:                     entityBase.Alias,
		Description:               entityBase.Description,
		HasDescription:            len(entityBase.Description) > 0,
		HasPrefix:                 false, //len(entityBase.Prefix) > 0,
		NameWithoutPrefix:         "",
		Columns:                   entityBase.Columns,
		NonPrimaryColumns:         entityBase.NonPrimaryColumns,
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

	for i, column := range from.Columns {
		col := toColumnSection(column, out.Name)

		if i == 0 {
			col.IsFirst = true
		}

		if i == len(from.Columns)-1 {
			col.IsLast = true
		}

		out.Columns = append(out.Columns, col)

		if column.IsPrimaryKey {
			out.PrimaryKeys = append(out.PrimaryKeys, col)
		} else {
			out.NonPrimaryColumns = append(out.NonPrimaryColumns, col)
		}

		if column.IsForeignKey {
			out.ForeignKeys = append(out.ForeignKeys, col)
		}
	}

	for _, cr := range from.ChildRelations {
		child := models.EntityRelationModel{}

		child.OwnerShcema = out.Schema
		child.OwnerName = out.Name
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

		// children = append(children)
		child.NameWithoutPrefix = strings.TrimPrefix(cr.Name, out.Name)
		child.HasPrefix = strings.HasPrefix(cr.Name, out.Name)

		child.HasSynonym = cr.HasSynonym
		if child.HasSynonym {
			child.Synonym = cr.Synonym
		}

		child.ModelName = coalesceString(cr.Synonym, cr.Name)

		if cr.Columns != nil && len(cr.Columns) > 0 {
			for ci, col := range cr.Columns {
				childCol := toColumnSection(col, child.Name)
				if ci == 0 {
					childCol.IsFirst = true
				}

				if ci == len(cr.Columns)-1 {
					childCol.IsLast = true
				}

				child.Columns = append(child.Columns, childCol)
				if col.IsPrimaryKey {
					child.PrimaryKeys = append(child.PrimaryKeys, childCol)
				} else {
					child.NonPrimaryColumns = append(child.NonPrimaryColumns, childCol)
				}

				if col.IsForeignKey {
					child.ForeignKeys = append(child.ForeignKeys, childCol)
				}
			}
		}

		out.Children = append(out.Children, child)
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
