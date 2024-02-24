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

INSERT INTO cheater VALUES (default);

-- we sync only rows with profile=public
--
INSERT INTO race (name, profile) VALUES ('Vendée Globe', 'public');
INSERT INTO race (name, profile) VALUES ('Globe Challenge', 'public');
INSERT INTO race (name, profile) VALUES ('RORC Caribbean 600', 'public');
INSERT INTO race (name, profile) VALUES ('Plastimo Lorient Mini - PLM 6.50', 'public');
INSERT INTO race (name, profile) VALUES ('Les potes et moi', 'private');
