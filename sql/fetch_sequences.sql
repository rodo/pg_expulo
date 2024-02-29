WITH sequences AS (

SELECT schemaname || '.' || sequencename as sequencename,
last_value,

'nextval(''' || schemaname || '.' || sequencename || '''::regclass)' as default_value

FROM pg_catalog.pg_sequences

),

columns AS (

SELECT table_schema || '.' ||  table_name as tablename,
column_name, ordinal_position,
column_default

FROM information_schema.columns
)


SELECT
        tablename,
        column_name,
        sequencename,
        NEXTVAL(sequencename),
        NEXTVAL(sequencename),-- need to repeat to instantiate the Go Struct
        ordinal_position as column_position

FROM sequences as s
INNER JOIN columns c ON c.column_default = s.default_value
