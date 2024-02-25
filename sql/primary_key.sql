SELECT ci.relname AS index_name

FROM pg_index AS i
INNER JOIN pg_class AS ci ON i.indexrelid = ci.oid
INNER JOIN pg_class AS ct ON i.indrelid = ct.oid
INNER JOIN pg_namespace AS n ON ct.relnamespace = n.oid


WHERE
    indisprimary
    AND nspname = $1
    AND ct.relname = $2
