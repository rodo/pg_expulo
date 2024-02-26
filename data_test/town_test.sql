--- town_test.sql
---
--- This table is used to test the partial deletion of data
--- on target
---
BEGIN;
SELECT plan(2);

SELECT has_table('public', 'town', 'table town exists');

-- all the wanted rows are here and no more
SELECT results_eq(
    $$ SELECT count(*)::int FROM town WHERE area='South' $$,
    'SELECT 2::int',
    'There is 2 rows in table town');

SELECT * FROM finish();
ROLLBACK;
