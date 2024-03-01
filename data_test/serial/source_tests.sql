--- test on data source
---
BEGIN;
SELECT plan(5);

SELECT has_table('laires'::name, 'root'::name);
SELECT has_table('laires'::name, 'la'::name);
SELECT has_table('laires'::name, 'lb'::name);
SELECT has_table('laires'::name, 'laa'::name);
SELECT has_table('laires'::name, 'lbb'::name);


SELECT * FROM finish();
ROLLBACK;
