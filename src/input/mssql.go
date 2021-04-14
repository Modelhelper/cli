package input

import (
	"context"
	"database/sql"
	"log"
	"modelhelper/cli/app"

	_ "github.com/denisenkom/go-mssqldb"
)

type MsSql struct {
	Source app.ConfigSource
}

// type Database interface {
// 	GetEntity(entityName string)
// }

// func GetEntity(entityName string, db *sql.DB) (*types.Entity, error) {
// 	d := db.Driver()

// 	fmt.Println(d)
// 	return nil, nil
// }

func (server *MsSql) CanConnect() (bool, error) {
	return false, nil
}

func (server *MsSql) Entity(name string) (*Entity, error) {
	e, err := server.getEntity(name)
	if err != nil {
		return nil, err
	}

	c, err := server.getColumns(e.Schema, e.Name)
	if err != nil {
		return nil, err
	}

	e.Columns = *c

	return e, nil
}
func (server *MsSql) Entities(pattern string) (*[]Entity, error) {
	sql := `
	with rowcnt (object_id, rowcnt) as (
		SELECT p.object_id, SUM(CASE WHEN (p.index_id < 2) AND (a.type = 1) THEN p.rows ELSE 0 END) 
		FROM sys.partitions p 
		INNER JOIN sys.allocation_units a ON p.partition_id = a.container_id
		join sys.objects o on p.object_id = o.object_id and o.type = 'U'
		--where p.object_id = object_id('Add')
		group by p.object_id
		)
		select 
			o.name
			,type = CASE when o.type = 'U' then 'Table' when o.type = 'V' then 'View' end  
			,[Schema] = s.name
			, Alias = Left(o.name, 1)
			, [RowCount] = isnull(rc.RowCnt, 0)
			, Description = isnull(ep.value, '')
		from sys.objects o
		join sys.schemas s on s.schema_id = o.schema_id
		left join rowcnt rc on rc.object_id = o.object_id    
		left join sys.extended_properties ep on o.object_id = ep.major_id and minor_id = 0 and ep.name = 'MS_description'
		where o.name not in ('sysdiagrams') 
		and [type] in ('V', 'U')
		order by s.name, o.[type], o.name		
	`

	// --and type in {entityFilter}
	// {tableFilter}
	db, err := server.openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	ctx := context.Background()

	stmt, err := db.PrepareContext(ctx, sql)

	if err != nil {
		return nil, err
	}
	// Execute query
	rows, err := stmt.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	list := []Entity{}

	var e Entity

	for rows.Next() {

		if err := rows.Scan(
			&e.Name,
			&e.Type,
			&e.Schema,
			&e.Alias,
			&e.RowCount,
			&e.Description,
		); err != nil {
			return nil, err
		} else {
			list = append(list, e)

		}

	}
	// fmt.Println(sql)
	return &list, nil
}

func (server *MsSql) getEntity(entityName string) (*Entity, error) {

	db, err := server.openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	ctx := context.Background()

	query := `
	select 
	o.name
	,type = CASE when o.type = 'U' then 'Table' when o.type = 'V' then 'View' end  
	,[Schema] = s.name
    , description =  isnull(ep.value, '')
from sys.objects o
join sys.schemas s on s.schema_id = o.schema_id
left join sys.extended_properties ep on o.object_id = ep.major_id and minor_id = 0 and ep.name = 'MS_description'
where o.object_id = object_id(@entityName)	
	`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}
	// Execute query
	row := stmt.QueryRow(query, sql.Named("entityName", entityName))

	var e Entity

	if err := row.Scan(
		&e.Name,
		&e.Type,
		&e.Schema,
		&e.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &e, nil
}

func (server *MsSql) getColumns(schema string, entityName string) (*[]Column, error) {
	db, err := server.openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	ctx := context.Background()

	query := `
	with Reserved as (
		select Name = 'database' union
		select Name = 'version' union
		select Name = 'new' union
		select Name = 'tran' union
		select Name = 'add' union
		select Name = 'insert' union
		select Name = 'inner' union
		select Name = 'index' union
		select Name = 'column' union
		select Name = 'commit' union
		select Name = 'return'        
	),PrimaryKeyColumns as (

		SELECT  
			i.name AS IndexName
			, OBJECT_NAME(ic.OBJECT_ID) AS TableName
			, COL_NAME(ic.OBJECT_ID,ic.column_id) AS PrimaryColumnName
			, ColumnId = ic.column_id
			, ObjectId = ic.object_id
		FROM    sys.indexes AS i 
		INNER JOIN sys.index_columns AS ic ON  i.OBJECT_ID = ic.OBJECT_ID AND i.index_id = ic.index_id
		WHERE i.is_primary_key = 1 and i.object_id = object_id(@entityName)
	), ForeignKeyColumns as (
		select 
			  ColumnName = cc.name
			, ColumnId = cc.column_id
			, ObjectId = cc.object_id
			, ReferencedColumn = pcc.name  
			, ReferencedObjectId = pcc.object_id
			, IsSelfJoin = cast(case when fkc.parent_object_id = fkc.referenced_object_id then 1 else 0 end as bit )
		from sys.foreign_key_columns fkc
		join sys.columns cc on fkc.parent_column_id = cc.column_id and cc.object_id = fkc.parent_object_id
	   join sys.columns pcc on fkc.referenced_column_id = pcc.column_id and pcc.object_id = fkc.referenced_object_id
		where fkc.parent_object_id = OBJECT_ID(@entityName)
	)
	select
		  Name = c.name	        
		, Description = isnull(ep.value, '')
		--, ModelName = c.Name
		, DataType = TYPE_NAME(c.user_type_id)
		, DbType = TYPE_NAME(c.user_type_id)
		, IsNullable = c.is_nullable	        
		, IsIdentity = c.is_identity             
		, IsPrimaryKey = cast (case when pkc.PrimaryColumnName is null then 0 else 1 end as bit)
		, IsForeignKey = cast (case when fkc.ColumnName is null then 0 else 1 end as bit)
		--, IsIgnored = case when s.name is null then 0 else 1 end
		, IsReserved = cast (case when r.name is null then 0 else 1 end as bit)
		--, Selected = cast (1 as bit) --case when s.name is null then 1 else 0 end
		, [Collation] = isnull(c.collation_name, '')
		, Length = case 
			when c.user_type_id = 231 and c.max_length > 0 then c.max_length / 2
			when left(c.name, 1) = 'n' and st.max_length = 8000 then c.max_length / 2
			else c.max_length end
		, UseLength = cast(case when st.precision = 0 and c.collation_name is not null then 1 else 0 end as bit)
		, c.Precision
		, c.Scale
		, UsePrecision = cast(case when st.user_type_id in (108,106) then 1 else 0 end as bit)
		, ReferencesColumn = isnull(fkc.ReferencedColumn, '')
		, ReferencesTable = isnull(object_name(fkc.ReferencedObjectId), '')
	from sys.columns c         
	left join sys.types st on st.user_type_id = c.user_type_id
   -- left join IgnoredColumns s on s.Name = c.name
	left join Reserved r on r.name = c.name
	left join PrimaryKeyColumns pkc on pkc.ColumnId = c.column_id and pkc.ObjectId = c.object_id -- c.name
	left join ForeignKeyColumns fkc on fkc.ColumnId = c.column_id and c.object_id = fkc.ObjectId
	left join sys.extended_properties ep on c.object_id = ep.major_id and minor_id = c.column_id and ep.name = 'MS_description'
	where object_id = object_id(@entityName)
	order by c.column_id
	`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}
	// Execute query
	rows, err := stmt.Query(query, sql.Named("entityName", entityName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cl := []Column{}
	var c Column

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
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, err
			}
		} else {
			cl = append(cl, c)
		}
	}
	return &cl, nil
}

func (server *MsSql) openConnection() (*sql.DB, error) {

	cs := server.Source.ConnectionString
	// fmt.Println("Connect with: " + cs)
	// var err error
	db, err := sql.Open("sqlserver", cs)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	return db, nil
}
