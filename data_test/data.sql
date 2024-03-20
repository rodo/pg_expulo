/*
 * Play this file on source only
 */
SET client_min_messages TO WARNING;

TRUNCATE results CASCADE;
TRUNCATE race CASCADE;
TRUNCATE skipper CASCADE;
TRUNCATE boat CASCADE;

INSERT INTO architect (id, name) VALUES (1, 'Cabinet VPLP');
INSERT INTO architect (id, name) VALUES (2, 'André Allègre');
INSERT INTO architect (id, name) VALUES (3, 'Eric Tabarly');
ALTER SEQUENCE architect_id_seq RESTART WITH 10 ;

INSERT INTO boat (name, length) VALUES ('Pen Duick I', 17.6);
INSERT INTO boat (name, updated_at) VALUES ('Pen Duick II', now());
INSERT INTO boat (name, architect) VALUES ('Pen Duick III', 3);
INSERT INTO boat (name, architect) VALUES ('Pen Duick IV',2);


-- 10 skippers
-- 1 with email
-- 4 with country set to null
INSERT INTO skipper (name, email) VALUES ('Eric Tabarly', 'eric@tabarly.domain');
INSERT INTO skipper (name, country) VALUES ('Florence Arthaud','France'), ('Catherine Chabaud','France');
INSERT INTO skipper (name) VALUES ('Loïck Peyron');

INSERT INTO skipper (name, updated_at) VALUES ('Guirec Soudée', now()), ('Loïck Peyron', now()),
('Damien Seguin', now()),
('Anne Caseneuve', now()), ('Clarisse Crémer', now()), ('Ellen MacArthur', now());

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
-- 55 results
--
--
INSERT INTO results (race_id, year) SELECT 1, generate_series(1960, 2014, 1);


--
-- we sync only rows with area=North
-- data for other area exists on target
INSERT INTO town (name, area) VALUES ('Brest', 'North');

--
--
INSERT INTO team (id) VALUES (generate_series(1,10,1));
UPDATE team SET name = 'Team ' || id;


INSERT INTO sponsor (name, team) VALUES ('Sponsor 1',1);
INSERT INTO sponsor (name, team) VALUES ('Sponsor 2',2);
INSERT INTO sponsor (name, team) VALUES ('Sponsor 3',3);
INSERT INTO sponsor (name, team) VALUES ('Sponsor 4',4);
