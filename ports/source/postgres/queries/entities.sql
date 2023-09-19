select 
    c.relname , 
    'Table' as table_type, 
    n.nspname, 
    '' as alias,
    reltuples::int,
    COALESCE(d.description, ''), 
    0 as ColumnCount,
    0 as NullableCount,
    0 as IdentityCount,
    0 as ChildrenCount,
    0 as ParentCount,
    False as IsVersioned,
    False as IsHistory,
    '' as HistoryTable
    --* 
from pg_class c
join pg_namespace n on n.oid = c.relnamespace
left join pg_description d on d.classoid = c.oid
join information_schema.tables tt on tt.table_schema = n.nspname and table_name = c.relname
where table_schema not in ('pg_catalog', 'information_schema') %s