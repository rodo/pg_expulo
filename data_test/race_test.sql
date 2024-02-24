--- race_test.sql
BEGIN;
SELECT plan(1);

SELECT has_table('public', 'race', 'table race exists');

SELECT * FROM finish();
ROLLBACK;
