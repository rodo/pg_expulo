WITH sequences AS (

SELECT schemaname || '.' || sequencename as sequencename,
last_value,

'nextval(''' || schemaname || '.' || sequencename || '''::regclass)' as default_value

FROM pg_catalog.pg_sequences

WHERE schemaname NOT IN ('information_schema','pg_catalog','public')

UNION ALL
--
-- weird but sequences in public schema are not prefixed by schema name
-- sure it will be a good bug to solve with databases with another default schema
--
SELECT schemaname || '.' || sequencename as sequencename,
last_value,

'nextval(''' || sequencename || '''::regclass)' as default_value

FROM pg_catalog.pg_sequences

WHERE schemaname IN ('public')

),

columns AS (

SELECT
        table_schema,
        table_name,
        column_name,
        ordinal_position,
        column_default

FROM information_schema.columns
WHERE table_schema NOT IN ('information_schema','pg_catalog')
)

SELECT
        column_name

FROM sequences as s
INNER JOIN columns c ON c.column_default = s.default_value

WHERE table_schema = $1 AND table_name = $2
