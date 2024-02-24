--- race_test.sql
BEGIN;
SELECT plan(2);

SELECT has_table('public', 'race', 'table race exists');

-- all the wanted rows are here and no more
SELECT results_eq(
    'SELECT count(*)::int FROM race',
    'SELECT 4::int',
    'There is 4 rows in table race');

SELECT * FROM finish();
ROLLBACK;
