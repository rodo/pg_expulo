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


INSERT INTO results (year) SELECT generate_series(1960, 2014, 1);

INSERT INTO cheater VALUES (DEFAULT);

-- we sync only rows with profile=public
--
INSERT INTO race (name, profile) VALUES ('Vendée Globe', 'public');
INSERT INTO race (name, profile) VALUES ('Globe Challenge', 'public');
INSERT INTO race (name, profile) VALUES ('RORC Caribbean 600', 'public');
INSERT INTO race (name, profile) VALUES (
    'Plastimo Lorient Mini - PLM 6.50', 'public'
);
INSERT INTO race (name, profile) VALUES ('Les potes et moi', 'private');
--
--
--
INSERT INTO linked.root (name) VALUES ('one');
INSERT INTO linked.root (name) VALUES ('two');
INSERT INTO linked.sheetone (name) VALUES ('one');
INSERT INTO linked.sheetone (name) VALUES ('two');
INSERT INTO linked.sheettwo (name) VALUES ('one');
INSERT INTO linked.sheetthree (name) VALUES ('one');
INSERT INTO linked.fish (aname, bname) VALUES ('one', 'one');
--
-- we sync only rows with area=North
-- data for other area exists on target
INSERT INTO town (name, area) VALUES ('Brest', 'North');
