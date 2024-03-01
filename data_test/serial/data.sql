/*
 * Create tables and data to validate foreign keys and serial
 */
SET client_min_messages TO WARNING;

BEGIN;

SET search_path = 'laires';

INSERT INTO root (name) SELECT generate_series(1, 4000, 1)::text;

INSERT INTO la (topid,name) SELECT generate_series(1, 100, 1), generate_series(1, 100, 1);

INSERT INTO lb (topid,name) SELECT generate_series(2, 100, 2), generate_series(1, 100, 1);


INSERT INTO lbb (lbid,name) SELECT id, name FROM lb ORDER BY RANDOM() LIMIT 10;
INSERT INTO lbb (lbid,name) SELECT id, name FROM lb ORDER BY RANDOM() LIMIT 100;

INSERT INTO laa (laid,name) SELECT id, name FROM la ORDER BY RANDOM() LIMIT 10;
INSERT INTO laa (laid,name) SELECT id, name FROM la ORDER BY RANDOM() LIMIT 100;


COMMIT ;
