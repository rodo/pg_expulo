--- Used to check the -purge option
BEGIN;
SELECT plan(2);


SELECT results_eq('SELECT count(*)::int FROM boat', 'SELECT 0::int', 'There is 0 rows in table boat');
SELECT results_eq('SELECT count(*)::int FROM skipper', 'SELECT 0::int', 'There is 0 rows in table skipper');

SELECT * FROM finish();
ROLLBACK;
