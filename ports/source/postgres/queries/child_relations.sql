;with columns as (

select  
    a.attrelid, a.attnum, a.attname col_name, t.typname col_type, a.atttypid::regtype col_type_alt, 
    -- coalesce(c.contype = 'p', false) is_pk, coalesce(c.contype = 'f', false) is_fk, 
    coalesce(attnotnull = false, true) is_nullable
from pg_attribute a
join pg_type t on t.oid = a.atttypid
where a.attnum > 0
)
SELECT 
  tbl_ref.relname,
  tbl_ref.relnamespace::regnamespace,
    '' some_type,
  col.col_name me_column_name, 
  col.col_type me_column_type, 
  col.is_nullable me_is_nullable, 
  col_ref.col_name ref_column_name, 
  col_ref.col_type ref_column_type, 
  col_ref.is_nullable ref_is_nullable,
  conname AS contraint_name,
  coalesce(conrelid=confrelid, false) is_self_join
  
FROM pg_constraint c
  join pg_class tbl_ref on tbl_ref.oid = c.conrelid
  join columns col_ref on col_ref.attrelid = c.conrelid and col_ref.attnum = any(c.conkey)
  join columns col on col.attrelid = c.confrelid and col.attnum = any(c.confkey)

WHERE  contype in ('f') 
and confrelid = '%s'::regclass
ORDER  BY conrelid::regclass::text, contype DESC;


--   pg_get_constraintdef(c.oid), 
--   confrelid as ref_table_oid, 
--   confrelid::regclass, conindid as class_oid, 