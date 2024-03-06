/*
 * Play this file on source only
 */


TRUNCATE results CASCADE;
TRUNCATE race CASCADE;
TRUNCATE skipper CASCADE;
TRUNCATE boat CASCADE;

INSERT INTO boat (name, length) VALUES ('Pen Duick I', 17.6);
INSERT INTO boat (name, updated_at) VALUES ('Pen Duick II', now());
INSERT INTO boat (name) VALUES ('Pen Duick III');
INSERT INTO boat (name) VALUES ('Pen Duick IV');

-- 10 skippers
-- 1 with email
--
INSERT INTO skipper (name, email) VALUES ('Eric Tabarly', 'eric@tabarly.domain');
INSERT INTO skipper (name,country) VALUES ('Florence Arthaud','France'), ('Catherine Chabaud','France');
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
