;

with columns as
    ( select a.attrelid,
             a.attnum,
             a.attname col_name,
             t.typname col_type,
             a.atttypid::regtype col_type_alt, -- coalesce(c.contype = 'p', false) is_pk, coalesce(c.contype = 'f', false) is_fk,
 coalesce(attnotnull = false, true) is_nullable
     from pg_attribute a -- left join pg_constraint c on c.conrelid = a.attrelid and a.attnum = any(c.conkey)
join pg_type t on t.oid = a.atttypid
     where a.attnum > 0 -- where a.attrelid = 'person.name'::regclass
) -- parents

SELECT tbl_ref.relname AS table_name,
       tbl_ref.relnamespace::regnamespace,
       '' some_type,
          col.col_name me_column_name,
          col.col_type me_column_type,
          col.is_nullable me_is_nullable,
          col_ref.col_name ref_column_name,
          col_ref.col_type ref_column_type,
          col_ref.is_nullable ref_is_nullable,
          conname AS constraint_name,
          coalesce(conrelid=confrelid, false) is_self_join
FROM pg_constraint c
left join pg_class tbl_ref on tbl_ref.oid = c.confrelid
join columns col on col.attrelid = c.conrelid
and col.attnum = any(c.conkey)
join columns col_ref on col_ref.attrelid = c.confrelid
and col_ref.attnum = any(c.confkey)
WHERE contype in ('f')
    and conrelid = '%s'::regclass 
    
    -- confrelid::regclass AS table_name,
 -- conname AS constraint_name,
 -- pg_get_constraintdef(c.oid),
 -- confrelid as ref_table_oid,
 -- confrelid::regclass,