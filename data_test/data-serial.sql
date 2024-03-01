/*
 * Create tables and data to validate foreign keys and serial
 */
BEGIN;

DROP SCHEMA IF EXISTS laires CASCADE;
CREATE SCHEMA IF NOT EXISTS laires;

SET search_path = 'laires';

CREATE TABLE laires.root (id serial PRIMARY KEY, color int default 0, name text DEFAULT 'none');
CREATE TABLE laires.la   (id serial PRIMARY KEY, topid int REFERENCES root(id), name text );
CREATE TABLE laires.lb   (id serial PRIMARY KEY, topid int REFERENCES root(id), name text );
CREATE TABLE laires.laa  (id serial PRIMARY KEY, laid  int REFERENCES la(id), name text );
CREATE TABLE laires.lbb  (id serial PRIMARY KEY, lbid  int REFERENCES lb(id), name text );

INSERT INTO root (name) SELECT generate_series(1, 40, 1)::text;

INSERT INTO la (topid,name) SELECT generate_series(1, 10, 1), generate_series(1, 10, 1);

INSERT INTO lb (topid,name) SELECT generate_series(2, 10, 1), generate_series(1, 10, 1);


COMMIT ;
