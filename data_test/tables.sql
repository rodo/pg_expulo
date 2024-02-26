SET client_min_messages TO WARNING;

DROP TABLE IF EXISTS race;
DROP TABLE IF EXISTS town;
DROP TABLE IF EXISTS results;
DROP TABLE IF EXISTS boat;
DROP TABLE IF EXISTS skipper;
DROP TABLE IF EXISTS cheater;

CREATE TABLE boat (
  id serial primary key,
  name text,
  created_at timestamp with time zone default now(),
  updated_at timestamp without time zone,
  length integer,
  class_code varchar(10),
  architect text
);

CREATE TABLE skipper (
  id serial primary key,
  name text,
  email text,
  age int,
  created_at timestamp with time zone default now(),
  updated_at timestamp without time zone
);

/* We sync only race with profile=public
 *
 */
CREATE TABLE race (
  id serial primary key,
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

/*
 *
 *
 */
DROP SCHEMA IF EXISTS linked CASCADE;
CREATE SCHEMA IF NOT EXISTS linked;


CREATE TABLE linked.root (name text default 'fool' PRIMARY KEY);
CREATE TABLE linked.sheetone (name text default 'fool' PRIMARY KEY REFERENCES linked.root(name));
CREATE TABLE linked.sheettwo (name text default 'fool' PRIMARY KEY REFERENCES linked.root(name));
CREATE TABLE linked.sheetthree (name text default 'fool' PRIMARY KEY REFERENCES linked.sheetone(name));
CREATE TABLE linked.fish (
  aname text default 'fool' REFERENCES linked.sheetone(name),
  bname text default 'fool' REFERENCES linked.sheettwo(name) );
