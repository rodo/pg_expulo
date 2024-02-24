--- boat_test.sql
BEGIN;
SELECT plan(2);

SELECT has_table('public', 'boat', 'table boat exists');

SELECT results_eq(
    'SELECT count(*)::int FROM boat',
    'SELECT 4::int',
    'There is 4 rows in table boat');

SELECT * FROM finish();
ROLLBACK;
