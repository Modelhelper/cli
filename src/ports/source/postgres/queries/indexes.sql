select i.indexrelid,
       c.relname,
       i.indisprimary,
       i.indisunique,
       i.indisclustered
from pg_index i
join pg_class c on i.indexrelid = c.oid
where i.indrelid = '%s'::regclass