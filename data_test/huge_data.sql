/*
 * Append a lot of data to test with volumetry
 */



INSERT INTO boat (name, length) VALUES ('Pen Duick I', 17.6);
INSERT INTO boat (name, updated_at) VALUES ('Pen Duick II', now());
INSERT INTO boat (name) VALUES ('Pen Duick III');
INSERT INTO boat (name) VALUES ('Pen Duick IV');

-- 10 skippers
-- 1 with email
--

INSERT INTO skipper (name) SELECT 'name' FROM generate_series(1,10e4,1);




-- 55 results
--
--
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

--
-- we sync only rows with area=North
-- data for other area exists on target
INSERT INTO town (name, area) SELECT gen_random_uuid()::text, 'North' FROM generate_series(1,10e4,1);
