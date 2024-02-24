--- cheater_test.sql
---
--- This table is not present in configuration so stay empty
---
BEGIN;
SELECT plan(2);

SELECT has_table('public', 'cheater', 'table cheater exists');

-- all the wanted rows are here and no more
SELECT results_eq(
    'SELECT count(*)::int FROM cheater',
    'SELECT 0::int',
    'There is 0 rows in table cheater');


SELECT * FROM finish();
ROLLBACK;
