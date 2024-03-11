/*
 * Create tables and data to validate foreign keys and serial
 */
SET client_min_messages TO WARNING;

BEGIN;

DROP SCHEMA IF EXISTS books CASCADE;
CREATE SCHEMA IF NOT EXISTS books;

DROP SCHEMA IF EXISTS editors CASCADE;
CREATE SCHEMA IF NOT EXISTS editors;



SET search_path = 'books';

CREATE TABLE editors.editor (id serial PRIMARY KEY, name text );

INSERT INTO editors.editor (name) VALUES ('Livre de poche');

CREATE TABLE author (id serial PRIMARY KEY, name text );

INSERT INTO author (name) VALUES ('Victor Hugo');
INSERT INTO author (name) VALUES ('Jules Verne');

CREATE TABLE books (id serial PRIMARY KEY, title text, author int REFERENCES author(id), editor int REFERENCES editors.editor(id) );

INSERT INTO books (title, author) VALUES ('L''assomoir', 1);
INSERT INTO books (title, author, editor) VALUES ('La Jangada', 2, 1);
INSERT INTO books (title, author) VALUES ('Le Rayon Vert', 2);


CREATE TABLE quote (id serial PRIMARY KEY, title text, book int REFERENCES books(id) );

INSERT INTO quote (title, book) VALUES ('Ça ne promet pas beaucoup de bonheur', 1);

COMMIT ;

GRANT USAGE ON SCHEMA books TO expulo;
REVOKE ALL ON ALL TABLES IN SCHEMA books FROM expulo;
REVOKE ALL ON ALL SEQUENCES IN SCHEMA books FROM expulo;
GRANT SELECT on ALL TABLES IN SCHEMA books TO expulo ;

GRANT USAGE ON SCHEMA editors TO expulo;
REVOKE ALL ON ALL TABLES IN SCHEMA editors FROM expulo;
REVOKE ALL ON ALL SEQUENCES IN SCHEMA editors FROM expulo;
GRANT SELECT on ALL TABLES IN SCHEMA editors TO expulo ;
