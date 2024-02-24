--- results_test.sql
BEGIN;
SELECT plan(1);

SELECT has_table('public', 'results', 'table results exists');

SELECT * FROM finish();
ROLLBACK;
