/*
 * Create tables and data to validate foreign keys
 */
BEGIN;

DROP SCHEMA IF EXISTS sunset CASCADE;
CREATE SCHEMA IF NOT EXISTS sunset;

CREATE TABLE sunset.root (foocolor int default 0, fooid serial PRIMARY KEY, fooname text);
CREATE TABLE sunset.classification (fooid int primary key, fooname text);
CREATE TABLE sunset.fruit (fooname text);

CREATE TABLE sunset.variety (foocolor int default 0, fooid int  PRIMARY KEY, fooname text,
fooclassid int REFERENCES sunset.classification(fooid));


INSERT INTO sunset.classification (fooid, fooname) VALUES (1, 'Vegetables');
INSERT INTO sunset.classification (fooid, fooname) VALUES (2, 'Fruits');


INSERT INTO sunset.variety (fooid, fooname, fooclassid) VALUES (1, 'Boscop', 1);
INSERT INTO sunset.fruit (fooname) VALUES ('Apple');
INSERT INTO sunset.fruit (fooname) VALUES ('Peach');

-- Need a table with serial to valfooidate the append method at first
--
INSERT INTO sunset.root (fooname) VALUES ('The root level');
--
ALTER TABLE sunset.fruit ADD COLUMN foovarid int REFERENCES sunset.variety(fooid);
UPDATE sunset.fruit SET foovarid = 1;
--
--

CREATE TABLE sunset.toplevel (fooid int PRIMARY KEY, fooname text);
INSERT INTO sunset.toplevel VALUES (1, 'Level 1'),(2, 'Level 2'),(3, 'Level 3');
ALTER TABLE sunset.classification ADD COLUMN levelfooid int REFERENCES sunset.toplevel(fooid);
UPDATE sunset.classification SET levelfooid = 2;


CREATE TABLE sunset.secondlevel (fooid int PRIMARY KEY, fooname text,
topfooid int REFERENCES sunset.toplevel(fooid));

--
--

--SET CONSTRAINTS ALL DEFERRED;
ALTER TABLE sunset.fruit ALTER CONSTRAINT fruit_foovarid_fkey  INITIALLY DEFERRED;
INSERT INTO sunset.fruit (fooname,foovarid) VALUES ('Peach', 2);
INSERT INTO sunset.variety (fooid, fooname, fooclassid) VALUES (2, 'Boscop', 1);

COMMIT;
ALTER TABLE sunset.fruit ALTER CONSTRAINT fruit_foovarid_fkey  NOT DEFERRABLE;
