SELECT
        n.nspname as schema,
        c.relname as tablename

FROM pg_class c
INNER JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE
        c.relkind = 'r'
AND
        n.nspname NOT IN ('pg_catalog', 'pg_toast', 'information_schema')
AND
        n.nspname IN ($1)

ORDER BY n.nspname, c.relname
