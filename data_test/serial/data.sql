/*
 * Create tables and data to validate foreign keys and serial
 */
SET client_min_messages TO WARNING;

BEGIN;

SET search_path = 'laires';

INSERT INTO root (name) SELECT generate_series(1, 40, 1)::text;

INSERT INTO la (topid,name) SELECT generate_series(1, 10, 1), generate_series(1, 10, 1);

INSERT INTO lb (topid,name) SELECT generate_series(2, 10, 1), generate_series(1, 10, 1);


COMMIT ;
