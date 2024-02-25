WITH tbl AS (
    SELECT
        0::bigint AS nb_fk,
        CONCAT(n.nspname, '.', t.relname) AS tablename
    FROM pg_class AS t
    LEFT JOIN pg_namespace AS n ON t.relnamespace = n.oid
    LEFT JOIN pg_constraint AS c ON t.oid = c.confrelid

    WHERE t.relkind = 'r' AND c.oid IS NULL

),

cte AS (
    SELECT
        n.nspname AS schemaname,
        cs.relname AS table_src,
        CONCAT(n.nspname, '.', cd.relname) AS tablename,
        COUNT(cd.relname) OVER (PARTITION BY cd.relname) AS nb_fk

    FROM pg_constraint AS t
    INNER JOIN pg_class AS cs ON t.conrelid = cs.oid
    INNER JOIN pg_class AS cd ON t.confrelid = cd.oid
    LEFT JOIN pg_namespace AS n ON cd.relnamespace = n.oid
    WHERE t.contype = 'f'
)

SELECT DISTINCT
    tablename,
    nb_fk
FROM cte
UNION ALL
SELECT DISTINCT
    tablename,
    nb_fk
FROM tbl
WHERE tablename in ($1)
ORDER BY nb_fk ASC
