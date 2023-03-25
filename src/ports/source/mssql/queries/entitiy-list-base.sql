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
		where 
			o.name not in ('sysdiagrams') 
			%s
			and o.[type] in ('V', 'U', 'SN', 'P')
		order by s.name, o.[type], o.name		