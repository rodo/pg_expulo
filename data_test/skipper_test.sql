--- skipper_test.sql
BEGIN;
SELECT plan(2);

SELECT has_table('public', 'skipper', 'table skipper exists');

SELECT results_eq(
    'SELECT count(*)::int FROM skipper',
    'SELECT 1::int',
    'There is 3 rows in table skipper');

SELECT * FROM finish();
ROLLBACK;
