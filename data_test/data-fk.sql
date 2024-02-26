/*
 * Create tables and data to validate foreign keys
 */


DROP SCHEMA IF EXISTS sunset CASCADE;
CREATE SCHEMA IF NOT EXISTS sunset;

TRUNCATE sunset.fruit CASCADE;
TRUNCATE sunset.variety CASCADE;
TRUNCATE sunset.classification CASCADE;
TRUNCATE sunset.toplevel CASCADE;

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
