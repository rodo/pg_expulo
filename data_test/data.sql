/*
 * Play this file on source only
 */

TRUNCATE results;
TRUNCATE race;
TRUNCATE skipper;
TRUNCATE boat;

INSERT INTO boat (name, length) VALUES ('Pen Duick I', 17.6);
INSERT INTO boat (name, updated_at) VALUES ('Pen Duick II', now());
INSERT INTO boat (name) VALUES ('Pen Duick III');
INSERT INTO boat (name) VALUES ('Pen Duick IV');

INSERT INTO skipper (name) VALUES ('Eric Tabarly');


INSERT INTO results (year) VALUES (1960);
