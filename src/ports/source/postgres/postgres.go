package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/utils/text"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	// _ "github.com/jackc/pgx/v5"
)

var (
	//go:embed queries/columns.sql
	selectColumnsQuery string

	//go:embed queries/entities.sql
	selectEntitiesQuery string

	//go:embed queries/parent_relations.sql
	parentRelationsQuery string
	//go:embed queries/child_relations.sql
	childRelationsQuery string
	//go:embed queries/indexes.sql
	indexQuery string
)

type postgresSource struct {
	connectionService modelhelper.ConnectionService
	database          *models.GenericConnection[models.PostgresConnection]
}

func NewPostgresSource(cs modelhelper.ConnectionService, connectionName string) modelhelper.SourceService {
	src := &postgresSource{
		connectionService: cs,
	}

	src.database = loadConnection(cs, connectionName)
	return src
}

// Entities implements modelhelper.SourceService
func (server *postgresSource) Entities(pattern string) (*[]models.Entity, error) {
	filter := ""

	if len(pattern) > 0 {
		pattern = strings.Replace(pattern, "*", "%", -1)
		filter = fmt.Sprintf("And c.relname like '%s'", pattern)
	}

	db, err := server.openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	ctx := context.Background()
	return server.entites(db, ctx, filter)
}

// EntitiesFromColumn implements modelhelper.SourceService
func (server *postgresSource) EntitiesFromColumn(column string) (*[]models.Entity, error) {
	filter := ""

	if len(column) > 0 {
		column = strings.Replace(column, "*", "%", -1)
		filter = fmt.Sprintf("AND o.object_id in (select object_id from sys.columns where name like '%s')", column)
	}

	db, err := server.openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	ctx := context.Background()
	return server.entites(db, ctx, filter)
}

// EntitiesFromNames implements modelhelper.SourceService
func (server *postgresSource) EntitiesFromNames(names []string) (*[]models.Entity, error) {
	ls := []models.Entity{}

	for _, name := range names {
		ent, err := server.Entity(name)
		if err != nil {
			return nil, err
		}

		ls = append(ls, *ent)
	}
	return &ls, nil
}

// Entity implements modelhelper.SourceService
func (server *postgresSource) Entity(name string) (*models.Entity, error) {

	db, err := server.openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	ctx := context.Background()

	columns, err := server.getColumns(db, ctx, name)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf(" and c.oid = '%s'::regclass", name)
	entities, err := server.entites(db, ctx, filter)

	entity := &models.Entity{}
	if err != nil {
		return nil, err
	}

	if len(*entities) == 1 {
		for _, e := range *entities {

			entity.Name = e.Name
			entity.Schema = e.Schema
			entity.Description = e.Description
			entity.ColumnCount = len(*columns)
			entity.Columns = *columns
			entity.RowCount = e.RowCount
			entity.Alias = e.Alias

		}
	}

	syn, ok := server.database.Synonyms[entity.Name]
	entity.HasSynonym = ok
	if ok {
		entity.Synonym = syn //.Name
	}

	idxs, err := server.getIndexes(db, ctx, name)
	if err != nil {
		fmt.Printf("Err when fetching indexes %v", err)
	}
	entity.Indexes = *idxs

	pr, err := server.getParentRelations(db, ctx, name)
	if err != nil {
		fmt.Printf("Err when fetching indexes %v", err)
	}
	entity.ParentRelations = *pr

	cr, err := server.getChildRelations(db, ctx, name)
	if err != nil {
		fmt.Printf("Err when fetching indexes %v", err)
	}
	entity.ChildRelations = *cr
	return entity, nil
}

func loadConnection(cs modelhelper.ConnectionService, conname string) *models.GenericConnection[models.PostgresConnection] {
	c, err := cs.Connection(conname)
	if err != nil {
		return nil
	}

	connection, ok := c.(*models.GenericConnection[models.PostgresConnection])

	if !ok {
		return nil
	}

	return connection

}

func (server *postgresSource) openConnection() (*sql.DB, error) {

	cs := server.database.Connection.ConnectionString
	// fmt.Println("Connect with: " + cs)
	// var err error
	db, err := sql.Open("postgres", cs)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	return db, nil
}

func (server *postgresSource) entites(db *sql.DB, ctx context.Context, filter string) (*[]models.Entity, error) {

	sql := fmt.Sprintf(selectEntitiesQuery, filter)

	// --and type in {entityFilter}
	// {tableFilter}

	stmt, err := db.PrepareContext(ctx, sql)

	if err != nil {
		return nil, err
	}
	// Execute query
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := []models.Entity{}

	var e models.Entity

	for rows.Next() {

		if err := rows.Scan(
			&e.Name,
			&e.Type,
			&e.Schema,
			&e.Alias,
			&e.RowCount,
			&e.Description,
			&e.ColumnCount,
			&e.NullableColumnCount,
			&e.IdentityColumnCount,
			&e.ChildRelationCount,
			&e.ParentRelationCount,
			&e.IsVersioned,
			&e.IsHistory,
			&e.HistoryTable,
		); err != nil {
			return nil, err
		} else {
			e.Alias = strings.ToLower(text.Abbreviate(e.Name))

			syn, ok := server.database.Synonyms[e.Name]
			e.HasSynonym = ok
			if ok {
				e.Synonym = syn //.Name
			}

			list = append(list, e)

		}
	}
	// fmt.Println(sql)
	return &list, nil
}

func (server *postgresSource) getColumns(db *sql.DB, ctx context.Context, entityName string) (*models.ColumnList, error) {

	stmt, err := db.PrepareContext(ctx, fmt.Sprintf(selectColumnsQuery, entityName))

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cl := models.ColumnList{}
	var c models.Column

	for rows.Next() {

		if err := rows.Scan(
			&c.Name,
			&c.Description,
			&c.DataType,
			&c.DbType,
			&c.IsNullable,
			&c.IsIdentity,
			&c.IsPrimaryKey,
			&c.IsForeignKey,
			&c.IsReserved,
			&c.Collation,
			&c.Length,
			&c.UseLength,
			&c.Precision,
			&c.Scale,
			&c.UsePrecision,
			&c.ReferencesColumn,
			&c.ReferencesTable,
			&c.ColumnIndex,
			&c.ForCreate,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, err
			}
		} else {

			if c.UseLength {
				len := strconv.Itoa(c.Length)
				if c.Length < 0 {
					len = "max"
				}
				c.DbType = fmt.Sprintf("%s (%s)", c.DataType, len)

			}
			if c.UsePrecision {
				c.DbType = fmt.Sprintf("%s (%v, %v)", c.DataType, c.Precision, c.Scale)
			}
			cl = append(cl, c)
		}
	}
	return &cl, nil
}
func (server *postgresSource) getIndexes(db *sql.DB, ctx context.Context, entityName string) (*[]models.Index, error) {

	stmt, err := db.PrepareContext(ctx, fmt.Sprintf(indexQuery, entityName))

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Index{}
	var r models.Index

	for rows.Next() {

		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.IsPrimaryKey,
			&r.IsUnique,
			&r.IsClustered,
			// &r.AvgFragmentationPercent,
			// &r.AvgPageSpacePercent,
			// &r.AvgRecordSize,
			// &r.Rows,
		); err != nil {
			if err == sql.ErrNoRows {
				return &[]models.Index{}, nil
			} else {
				return &[]models.Index{}, err
			}
		} else {
			list = append(list, r)

			// _, err := list.WriteString(r)
			// if err != nil {
			// 	return nil, err
			// }

			// return fromJson([]byte(list.String()))
		}
	}

	return &list, nil
}
func (server *postgresSource) getParentRelations(db *sql.DB, ctx context.Context, entityName string) (*[]models.Relation, error) {

	stmt, err := db.PrepareContext(ctx, fmt.Sprintf(parentRelationsQuery, entityName))

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Relation{}
	var r models.Relation

	for rows.Next() {

		if err := rows.Scan(
			// &r.GroupIndex,
			&r.Name,
			&r.Schema,
			&r.Type,
			&r.ColumnName,
			&r.ColumnType,
			&r.ColumnNullable,
			&r.OwnerColumnName,
			&r.OwnerColumnType,
			&r.OwnerColumnNullable,
			&r.ContraintName,
			&r.IsSelfJoin,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, err
			}
		} else {

			syn, ok := server.database.Synonyms[r.Name]
			r.HasSynonym = ok
			if ok {
				r.Synonym = syn //.Name
			}

			list = append(list, r)
		}
	}
	return &list, nil
}
func (server *postgresSource) getChildRelations(db *sql.DB, ctx context.Context, entityName string) (*[]models.Relation, error) {

	stmt, err := db.PrepareContext(ctx, fmt.Sprintf(childRelationsQuery, entityName))

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Relation{}
	var r models.Relation

	for rows.Next() {

		if err := rows.Scan(
			// &r.GroupIndex,
			&r.Name,
			&r.Schema,
			&r.Type,
			&r.OwnerColumnName,
			&r.OwnerColumnType,
			&r.OwnerColumnNullable,
			&r.ColumnName,
			&r.ColumnType,
			&r.ColumnNullable,
			&r.ContraintName,
			&r.IsSelfJoin,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, err
			}
		} else {
			syn, ok := server.database.Synonyms[r.Name]
			r.HasSynonym = ok
			if ok {
				r.Synonym = syn //.Name
			}
			list = append(list, r)
		}
	}
	return &list, nil
}
