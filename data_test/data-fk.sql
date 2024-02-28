/*
 * Create tables and data to validate foreign keys
 */
BEGIN;

DROP SCHEMA IF EXISTS sunset CASCADE;
CREATE SCHEMA IF NOT EXISTS sunset;

CREATE TABLE sunset.classification (id int primary key, name text);
CREATE TABLE sunset.fruit (name text);

CREATE TABLE sunset.variety (id int  PRIMARY KEY, name text,
classid int REFERENCES sunset.classification(id));


INSERT INTO sunset.classification (id, name) VALUES (1, 'Vegetables');
INSERT INTO sunset.classification (id, name) VALUES (2, 'Fruits');


INSERT INTO sunset.variety (id, name, classid) VALUES (1, 'Boscop', 1);
INSERT INTO sunset.fruit (name) VALUES ('Apple');
INSERT INTO sunset.fruit (name) VALUES ('Peach');
--
ALTER TABLE sunset.fruit ADD COLUMN varid int REFERENCES sunset.variety(id);
UPDATE sunset.fruit SET varid = 1;
--
--

CREATE TABLE sunset.toplevel (id int PRIMARY KEY, name text);
INSERT INTO sunset.toplevel VALUES (1, 'Level 1'),(2, 'Level 2'),(3, 'Level 3');
ALTER TABLE sunset.classification ADD COLUMN levelid int REFERENCES sunset.toplevel(id);
UPDATE sunset.classification SET levelid = 2;


CREATE TABLE sunset.secondlevel (id int PRIMARY KEY, name text,
topid int REFERENCES sunset.toplevel(id));

--
--

--SET CONSTRAINTS ALL DEFERRED;
ALTER TABLE sunset.fruit ALTER CONSTRAINT fruit_varid_fkey  INITIALLY DEFERRED;
INSERT INTO sunset.fruit (name,varid) VALUES ('Peach', 2);
INSERT INTO sunset.variety (id, name, classid) VALUES (2, 'Boscop', 1);

COMMIT;
ALTER TABLE sunset.fruit ALTER CONSTRAINT fruit_varid_fkey  NOT DEFERRABLE;
