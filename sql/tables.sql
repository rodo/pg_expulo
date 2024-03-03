SELECT
        n.nspname || '.' || c.relname

FROM pg_class c
INNER JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE
        c.relkind = 'r'
AND
        n.nspname NOT in ('pg_catalog', 'pg_toast', 'information_schema')
