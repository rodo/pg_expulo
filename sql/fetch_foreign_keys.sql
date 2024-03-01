WITH tables AS ( -- qry fetch_table_foreign_keys

    SELECT
        t.oid,
        t.relname,
        CONCAT(n.nspname, '.', t.relname) AS tablename,
        a.attname,
        a.attnum

    FROM pg_class AS t
    INNER JOIN pg_attribute a ON a.attrelid = t.oid
    LEFT JOIN pg_namespace AS n ON t.relnamespace = n.oid
    WHERE CONCAT(n.nspname, '.', t.relname) = ANY ($1::text[])
--    WHERE CONCAT(n.nspname, '.', t.relname) = ANY ('{sunset.variety,sunset.fruit,sunset.classification,sunset.toplevel}'::text[])

    AND t.relkind = 'r' AND a.attnum > 0
),

cte AS (
    SELECT DISTINCT
        ts.tablename AS tablename,
        tt.tablename AS table_targ,
        conname,
        ts.attname column_linked,
        tt.attname column_target

    FROM pg_constraint AS c
    INNER JOIN tables AS ts ON (ts.oid = c.conrelid AND ts.attnum = c.conkey[1])
    INNER JOIN tables AS tt ON (tt.oid = c.confrelid AND tt.attnum = c.confkey[1])

    WHERE c.contype = 'f'
    -- TODO see how to manage foreign keys with 2 columns

)
SELECT
tablename,
table_targ,
conname,
tablename || '.' || column_linked,
table_targ || '.' || column_target

FROM cte
