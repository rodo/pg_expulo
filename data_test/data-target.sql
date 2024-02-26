/*
 * Play this file on target only
 */

TRUNCATE town;

-- we sync only rows with area=North
-- data for other area exists on target
INSERT INTO town (name, area) VALUES ('Sidney', 'South');
INSERT INTO town (name, area) VALUES ('Hobart', 'South');
