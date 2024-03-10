SET client_min_messages TO WARNING;

DROP TABLE IF EXISTS results;
DROP TABLE IF EXISTS race;
DROP TABLE IF EXISTS town;
DROP TABLE IF EXISTS boat;
DROP TABLE IF EXISTS skipper;
DROP TABLE IF EXISTS cheater;
DROP TABLE IF EXISTS architect;

-- no special config for this table
-- will inherited from default
CREATE TABLE architect (
  id serial PRIMARY KEY,
  name text,
  guid uuid
);


CREATE TABLE boat (
  id serial PRIMARY KEY,
  name text,
  created_at timestamp with time zone default now(),
  updated_at timestamp without time zone,
  length integer,
  class_code varchar(10),
  architect int REFERENCES architect(id)
);

CREATE TABLE skipper (
  id serial primary key,
  name text,
  email text,
  age int,
  created_at timestamp with time zone default now(),
  updated_at timestamp without time zone,
  town text,
  country text
);

/* We sync only race with profile=public
 *
 */
CREATE TABLE race (
  id serial PRIMARY KEY,
  name text,
  year int,
  profile text default 'private',
  rating float default 4.5
);

/* This table is never empty on destination
 *
 */
CREATE TABLE results (
  id serial primary key,
  race_id int REFERENCES race(id),
  year int
);

/* This table is not present in configuration so stay empty on
   destination
*/
CREATE TABLE cheater (name text default 'fool');

/* This table already contains data on destination and is not full
   purged
   Primary Key non integer
*/
CREATE TABLE town (name text PRIMARY KEY, area text default 'North');

--
--
--
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM expulo;
REVOKE ALL ON ALL SEQUENCES IN SCHEMA public FROM expulo;
GRANT SELECT on ALL TABLES IN SCHEMA public TO expulo ;
