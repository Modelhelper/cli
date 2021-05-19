package source

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

type MsSql struct {
	Connection Connection
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

	p, err := server.getParents(e.Schema, e.Name)
	if err != nil {
		return nil, err
	}

	e.ParentRelations = *p

	cr, err := server.getChildren(e.Schema, e.Name)
	if err != nil {
		return nil, err
	}

	e.ChildRelations = *cr

	idx, err := server.getIndexes(e.Schema, e.Name)
	e.Indexes = *idx

	return e, nil
}
func (server *MsSql) Entities(pattern string) (*[]Entity, error) {
	search := ""

	if len(pattern) > 0 {
		pattern = strings.Replace(pattern, "*", "%", -1)
		search = fmt.Sprintf("And o.Name like '%s'", pattern)
	}
	sql := fmt.Sprintf(`
	with rowcnt (object_id, rowcnt) as (
		SELECT p.object_id, SUM(CASE WHEN (p.index_id < 2) AND (a.type = 1) THEN p.rows ELSE 0 END) 
		FROM sys.partitions p 
		INNER JOIN sys.allocation_units a ON p.partition_id = a.container_id
		join sys.objects o on p.object_id = o.object_id and o.type = 'U'
		--where p.object_id = object_id('Add')
		group by p.object_id
	), colCnt(id, cnt, nullcnt, idcnt) as (
		select object_id, cnt = count(*), sum(cast(is_nullable as int)), sum(cast(is_identity as int))--, sum(cast(is_computed as int))
		from sys.columns 
		group by object_id
	), ParentRelCnt(id, cnt) as (
		select 
			id = parent_object_id, cnt = count(*) 
		from sys.foreign_key_columns
		group by parent_object_id
	), ChildrenRelCnt(id, cnt) as (
		select 
			id = referenced_object_id, cnt = count(*) 
		from sys.foreign_key_columns
		group by referenced_object_id
	)
		select 
			o.name
			,type = CASE 
				when o.type = 'U' then 'Table' 
				when o.type = 'V' then 'View' 
				when o.type = 'SN' then 'Synonym'
				when o.type = 'P' then 'Proc'
				end  
			,[Schema] = s.name
			, Alias = Left(o.name, 1)
			, [RowCount] = isnull(rc.RowCnt, 0)
			, Description = isnull(ep.value, '')
			, ColumnCount = isnull(cc.cnt, 0)
			, NullableCount = isnull(cc.nullcnt, 0)
			, IdentityCount = isnull(cc.idcnt, 0)
			, ChildrenCount = isnull(crc.cnt, 0)
			, ParentCount = isnull(prc.cnt, 0)
			, IsVersioned = case when t.temporal_type = 2 then 1 else 0 end
			, IsHistory = case when t.temporal_type = 1 then 1 else 0 end
			, HistoryTable = isnull(object_name(t.history_table_id), '')
		from sys.objects o
		join sys.schemas s on s.schema_id = o.schema_id
		left join sys.tables t on t.object_id = o.object_id
		left join rowcnt rc on rc.object_id = o.object_id    
		left join sys.extended_properties ep on o.object_id = ep.major_id and minor_id = 0 and ep.name = 'MS_description'
		left join colCnt cc on cc.id = o.object_id
		left join ChildrenRelCnt crc on crc.id = o.object_id
		left join ParentRelCnt prc on prc.id = o.object_id
		where o.name not in ('sysdiagrams') %s
		and o.[type] in ('V', 'U', 'SN', 'P')
		order by s.name, o.[type], o.name		
	`, search)

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
			e.Alias = strings.ToLower(Abbreviate(e.Name))
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
	;with rowcnt (object_id, cnt) as (
		SELECT p.object_id, SUM(CASE WHEN (p.index_id < 2) AND (a.type = 1) THEN p.rows ELSE 0 END) 
		FROM sys.partitions p 
		INNER JOIN sys.allocation_units a ON p.partition_id = a.container_id
		join sys.objects o on p.object_id = o.object_id and o.type = 'U'
		where p.object_id = object_id(@entityName)
		group by p.object_id
	)
select 
	o.name
	,type = CASE 
				when o.type = 'U' then 'Table' 
				when o.type = 'V' then 'View' 
				when o.type = 'SN' then 'Synonym'
				when o.type = 'P' then 'Proc'
				end  
	,[Schema] = s.name
    , description =  isnull(ep.value, '')
    -- , RowCount = rowcnt
    , [RowCount] = rcnt.cnt
    , IsVersioned = case when t.temporal_type = 2 then 1 else 0 end
	, IsHistory = case when t.temporal_type = 1 then 1 else 0 end
    , HistoryTable = isnull(object_name(t.history_table_id), '')
from sys.objects o
join sys.schemas s on s.schema_id = o.schema_id
left join sys.extended_properties ep on o.object_id = ep.major_id and minor_id = 0 and ep.name = 'MS_description'
left join sys.tables t on t.object_id = o.object_id
left join rowcnt rcnt on rcnt.object_id = o.object_id
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
		&e.RowCount,
		&e.IsVersioned,
		&e.IsHistory,
		&e.HistoryTable,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	e.Alias = Abbreviate(e.Name)

	return &e, nil
}

func (server *MsSql) getColumns(schema string, entityName string) (*ColumnList, error) {
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

	cl := ColumnList{}
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

func (server *MsSql) getParents(schema string, entityName string) (*[]Relation, error) {
	db, err := server.openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	ctx := context.Background()

	query := `
select 
	GroupIndex = row_number() over (partition by p1.name order by p1.create_date desc),
	p1.name, 
	[Schema] = SCHEMA_NAME(p1.schema_id),
    type = CASE when p1.type = 'U' then 'Table' when p1.type = 'V' then 'View' end,  	
	ForeignColumnName = cp.name, 
	ForeignColumnType = type_name(cp.user_type_id),
    ForeignColumnNullable = cp.is_nullable,

    PrimaryColumnName = cc.name,
    PrimaryColumnType = type_name(cc.user_type_id),
    PrimaryColumnNullable = cc.is_nullable,
	
	ConstraintName = o1.name, 
	IsSelfJoin = cast(case when fkc.parent_object_id = fkc.referenced_object_id then 1 else 0 end as bit )
from sys.foreign_key_columns fkc

join sys.objects o1 on o1.object_id = fkc.constraint_object_id
join sys.objects r1 on r1.object_id = fkc.parent_object_id
join sys.objects p1 on p1.object_id = fkc.referenced_object_id
join sys.columns cc on fkc.parent_column_id = cc.column_id and cc.object_id = fkc.parent_object_id
join sys.columns cp on fkc.referenced_column_id = cp.column_id and cp.object_id = fkc.referenced_object_id
where fkc.parent_object_id = OBJECT_ID(@entityName)
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

	list := []Relation{}
	var r Relation

	for rows.Next() {

		if err := rows.Scan(
			&r.GroupIndex,
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
			list = append(list, r)
		}
	}
	return &list, nil
}

func (server *MsSql) getChildren(schema string, entityName string) (*[]Relation, error) {
	db, err := server.openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	ctx := context.Background()

	query := `
select 	
	GroupIndex = row_number() over (partition by p1.name order by p1.create_date desc),
	p1.name ,
	[Schema] = SCHEMA_NAME(p1.schema_id),
    type = CASE when p1.type = 'U' then 'Table' when p1.type = 'V' then 'View' end,  	
	PrimaryColumnName = cp.name, 
    PrimaryColumnType = type_name(cp.user_type_id),
    PrimaryColumnNullable = cp.is_nullable,
    
	ForeignColumnName = cc.name,
    ForeignColumnType = type_name(cc.user_type_id),
    ForeignColumnNullable = cc.is_nullable,

	ConstraintName = o1.name

    , IsSelfJoin = cast(case when fkc.parent_object_id = fkc.referenced_object_id then 1 else 0 end as bit )
from sys.foreign_key_columns fkc

join sys.objects o1 on o1.object_id = fkc.constraint_object_id
join sys.objects r1 on r1.object_id = fkc.referenced_object_id
join sys.objects p1 on p1.object_id = fkc.parent_object_id
join sys.columns cc on fkc.parent_column_id = cc.column_id and cc.object_id = fkc.parent_object_id
join sys.columns cp on fkc.referenced_column_id = cp.column_id and cp.object_id = fkc.referenced_object_id
where fkc.referenced_object_id = OBJECT_ID(@entityName)
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

	list := []Relation{}
	var r Relation

	for rows.Next() {

		if err := rows.Scan(
			&r.GroupIndex,
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
			list = append(list, r)
		}
	}
	return &list, nil
}

func (server *MsSql) getIndexes(schema string, entityName string) (*[]Index, error) {
	db, err := server.openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	ctx := context.Background()

	query := `
	SELECT 
    a.index_id as 'id'
    , b.name
    , isnull(avg_fragmentation_in_percent, 0) as 'avgFragmentationPercent'
    , b.is_unique as 'isUnique', is_primary_key as 'IsPrimaryKey'
    , isnull(a.avg_page_space_used_in_percent, 0) as 'AvgPageSpacePercent'
    , isnull(a.avg_record_size_in_bytes, 0) as 'AvgRecordSize'
    , isnull(a.record_count, 0) as 'Rows'
    
FROM sys.dm_db_index_physical_stats (DB_ID(@database), OBJECT_ID(@table), NULL, NULL, NULL) AS a  
JOIN sys.indexes AS b ON a.object_id = b.object_id AND a.index_id = b.index_id

--left join sys.columns c on ic.column_id = c.column_id and c.object_id = ic.object_id
--for json path, INCLUDE_NULL_VALUES
;
	`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	database := server.Connection.ConnectionStringPart("database")
	// Execute query
	rows, err := stmt.Query(
		query,
		sql.Named("database", database),
		sql.Named("table", entityName),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []Index{}
	var r Index

	for rows.Next() {

		if err := rows.Scan(
			// &r,
			&r.ID,
			&r.Name,
			&r.AvgFragmentationPercent,
			&r.IsUnique,
			&r.IsPrimaryKey,
			&r.AvgPageSpacePercent,
			&r.AvgRecordSize,
			&r.Rows,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, err
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

func (server *MsSql) openConnection() (*sql.DB, error) {

	cs := server.Connection.ConnectionString
	// fmt.Println("Connect with: " + cs)
	// var err error
	db, err := sql.Open("sqlserver", cs)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	return db, nil
}

func fromJson(blob []byte) (*[]Index, error) {
	list := []Index{}

	if len(blob) > 0 {

		if err := json.Unmarshal(blob, &list); err != nil {
			return nil, err
		}
	}

	return &list, nil
}

type RelationTreeItem struct {
	KeyName           string
	ParentID          int
	ID                int
	RelatedTable      string
	RelatedColumnName string
	TableName         string
	ColumnName        string
}

func (server *MsSql) GetParentRelationTree(schema string, entityName string) (*[]RelationTreeItem, error) {
	db, err := server.openConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	ctx := context.Background()

	query := `
	;with track_parent as (

		-- select lvl = 1, p.name, p.parent_object_id, referenced_object_id, family = p.parent_object_id, path = cast(object_name(referenced_object_id) + ' > ' + object_name(parent_object_id) as nvarchar(1000) )
		select 
			--  direction = -1
			-- , lvl = 1, 
			p.name
			, p.parent_object_id
			, p.referenced_object_id
			, ParentTable = object_name(p.parent_object_id)
			, ParentColumn = cp.name
			, ReferencedTable = object_name(p.referenced_object_id)
			, ReferencedColumn = cc.name
			-- , family = p.referenced_object_id
			-- , path = cast(object_name(p.parent_object_id) + ' > ' + object_name(p.referenced_object_id) as nvarchar(1000) )
		from sys.foreign_keys p
		join sys.foreign_key_columns fkc on p.parent_object_id = fkc.parent_object_id and fkc.referenced_object_id = p.referenced_object_id
		join sys.columns cp on cp.column_id = fkc.parent_column_id and cp.object_id = fkc.parent_object_id
		join sys.columns cc on cc.column_id = fkc.referenced_column_id and cc.object_id = fkc.referenced_object_id
		where p.parent_object_id = object_id(@tablename)
		union all
		select
			--   direction
			-- , lvl = lvl + 1, 
			  fk.name
			  , fk.parent_object_id
			  , fk.referenced_object_id
			, ParentTable = object_name(fk.parent_object_id)
			, ParentColumn = cp.name
			, ReferencedTable = object_name(fk.referenced_object_id)
			, ReferencedColumn = cc.name
			-- , t.family
			-- , path = cast(path  + ' > ' + object_name(fk.referenced_object_id) as nvarchar(1000))
		from sys.foreign_keys fk
		join sys.foreign_key_columns fkc on fk.parent_object_id = fkc.parent_object_id and fkc.referenced_object_id = fk.referenced_object_id
		join sys.columns cp on cp.column_id = fkc.parent_column_id and cp.object_id = fkc.parent_object_id
		join sys.columns cc on cc.column_id = fkc.referenced_column_id and cc.object_id = fkc.referenced_object_id
		--join sys.objects o on o.object_id = fk.parent_object_id
		join track_parent t on t.referenced_object_id = fk.parent_object_id
		
		)
			select 
					-- direction = 0
					-- , lvl = 0
					  KeyName = ''
					, parentId = -1
					, Id = p.object_id
					, ParentTable = ''
					, ParentColumn = ''
					, Name =  p.name
					, ColumnName = ''
					-- , family =p.object_id
					-- , path = p.name
				from sys.tables p
				where p.object_id = object_id(@tablename)
			union all
			select * from track_parent
			-- order by lvl, family
		;
	`

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}
	// Execute query
	rows, err := stmt.Query(query, sql.Named("tablename", entityName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []RelationTreeItem{}
	var r RelationTreeItem

	for rows.Next() {

		if err := rows.Scan(
			&r.KeyName,
			&r.ParentID,
			&r.ID,
			&r.RelatedTable,
			&r.RelatedColumnName,
			&r.TableName,
			&r.ColumnName,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, err
			}
		} else {
			list = append(list, r)
		}
	}
	return &list, nil
}
