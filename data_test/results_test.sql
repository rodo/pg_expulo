--- results_test.sql
BEGIN;
SELECT plan(2);

SELECT has_table('public', 'results', 'table results exists');

SELECT results_eq(
    'SELECT count(*)::int FROM results',
    'SELECT 55::int',
    'There is 55 rows in table results'
);


SELECT * FROM finish();
ROLLBACK;
