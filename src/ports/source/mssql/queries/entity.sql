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