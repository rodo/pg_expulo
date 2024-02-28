WITH tables AS ( -- qry fetch_table_foreign_keys

    SELECT
        c.oid,
        c.relname,
        CONCAT(n.nspname, '.', c.relname) AS tablename

    FROM pg_class AS c
    LEFT JOIN pg_namespace AS n ON c.relnamespace = n.oid
    WHERE CONCAT(n.nspname, '.', c.relname) = ANY ($1::text[])
--    WHERE CONCAT(n.nspname, '.', c.relname) = ANY ('{sunset.variety,sunset.fruit,sunset.classification,sunset.toplevel}'::text[])

    AND c.relkind = 'r'
),

cte AS (
    SELECT DISTINCT

        ts.tablename AS tablename,
        -- tt.tablename AS table_targ,
        conname

    FROM pg_constraint AS c
    LEFT JOIN tables AS ts ON ts.oid = c.conrelid
    LEFT JOIN tables AS tt ON tt.oid = c.confrelid
    INNER JOIN pg_trigger AS tr
    ON (tr.tgconstraint = c.oid AND ts.oid = tr.tgrelid AND tt.oid = tr.tgconstrrelid)

    WHERE c.contype = 'f' AND tgenabled = 'O'
)

select * from cte
