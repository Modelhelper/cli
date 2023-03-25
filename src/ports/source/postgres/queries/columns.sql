with keys as (
    select c.conrelid, c.contype, unnest(conkey) col_id
    from pg_constraint c
    WHERE c.contype in ('f', 'p') 
        -- and c.conrelid = 'person.name'::regclass
)
select 
    c.attname, 
    '' as description,
    t.typname,
    t.typname as db_type,
    -- c.atttypid::regtype,
    case when c.attnotnull then False else True end as nullable,
    case when attidentity = 'a' then True else False end as id,
    COALESCE(pk.contype = 'p', False) as pk,
    COALESCE(fk.contype = 'f', False) as fk,
    -- case when i.indrelid is null then False else True end as idx,
    False as is_reserved,
    attcollation::regcollation as collation,
    -- COALESCE(i.indisclustered, False) clustered,
    0 as length,
    False as use_length,
    0 as prec,
    0 as scale,
    False as use_prec,
    '' as references_column,
    '' as references_table,
    c.attnum,
    '' as for_create
from pg_attribute c
-- join pg_description d on d.oid = c.
join pg_type t on t.oid = c.atttypid
left join keys pk on pk.conrelid = c.attrelid and pk.col_id = c.attnum and pk.contype = 'p'
left join keys fk on fk.conrelid = c.attrelid and fk.col_id = c.attnum and fk.contype = 'f'
left join pg_index i on c.attrelid = i.indrelid AND c.attnum = ANY(i.indkey)
where attnum > 0 
    and attrelid = '%s'::regclass
order by attnum