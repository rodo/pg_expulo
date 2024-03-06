--- skipper_test.sql
BEGIN;
SELECT plan(5);

-- without this table we won't do anything ;-)
SELECT has_table('public', 'skipper', 'table skipper exists');

-- all the wanted rows are here and no more
SELECT results_eq(
    'SELECT count(*)::int FROM skipper',
    'SELECT 10::int',
    'There is 10 rows in table skipper');

-- the column name is masked
SELECT row_eq(
    $$ SELECT count(*) FROM skipper WHERE name = 'Eric Tabarly' $$,
    ROW(0::bigint),
    'Column name is masked');

-- the column country is set with null preserved
SELECT row_eq(
    $$ SELECT count(*) FROM skipper WHERE country IS NULL $$,
    ROW(8::bigint),
    'Null are preserved on country');

-- min=42 max=42
SELECT row_eq(
    $$ SELECT sum(age) FROM skipper$$,
    ROW(420::bigint),
    'Min/max is correct');


SELECT * FROM finish();
ROLLBACK;
