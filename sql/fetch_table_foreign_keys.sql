WITH tables AS ( -- qry fetch_foreign_keys
    SELECT
        t.oid,
        n.nspname as schemaname,
        t.relname AS tablename,
        a.attname,
        a.attnum

    FROM pg_class AS t
    INNER JOIN pg_attribute a ON a.attrelid = t.oid
    LEFT JOIN pg_namespace AS n ON t.relnamespace = n.oid
--    WHERE CONCAT(n.nspname, '.', t.relname) = ANY ($1::text[])


    AND t.relkind = 'r' AND a.attnum > 0
),

cte AS (
    SELECT DISTINCT
        ts.schemaname,
        ts.tablename AS tablename,
        tt.tablename AS table_targ,
        conname,
        ts.attname column_linked,
        tt.attname column_target

    FROM pg_constraint AS c
    INNER JOIN tables AS ts ON (ts.oid = c.conrelid AND ts.attnum = c.conkey[1])
    INNER JOIN tables AS tt ON (tt.oid = c.confrelid AND tt.attnum = c.confkey[1])

    WHERE c.contype = 'f' AND conname NOT LIKE 'expulo_%'
    -- TODO see how to manage foreign keys with 2 columns

)
SELECT
  schemaname,
  tablename,
  table_targ,
  column_linked,
  column_target

FROM cte

--WHERE cte.schemaname = 'public' AND tablename= 'boat';
WHERE schemaname = $1 AND tablename= $2;
