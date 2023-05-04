	--declare @entityName nvarchar(100) = 'video'
    ;with Reserved as (
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
	),PrimaryKeys as (

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
	), ColList as (
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
        , ColumnIndex = c.column_id
        
	from sys.columns c         
	left join sys.types st on st.user_type_id = c.user_type_id
   -- left join IgnoredColumns s on s.Name = c.name
	left join Reserved r on r.name = c.name
	left join PrimaryKeys pkc on pkc.ColumnId = c.column_id and pkc.ObjectId = c.object_id -- c.name
	left join ForeignKeyColumns fkc on fkc.ColumnId = c.column_id and c.object_id = fkc.ObjectId
	left join sys.extended_properties ep on c.object_id = ep.major_id and minor_id = c.column_id and ep.name = 'MS_description'
	where object_id = object_id(@entityName)
		and generated_always_type = 0
)
select *
, ForCreate = + '[' + Name + ']'
            + ' ' + DbType
            + '' + case 
                when UseLength = 1 and Length > 0 then '(' + cast(Length as varchar(10)) + ')' 
                when UseLength = 1 and Length < 0 then '(max)' 
                else '' end
            + ' ' + case when IsNullable = 1 then 'NULL' else 'NOT NULL' end
 from ColList 
order by ColumnIndex