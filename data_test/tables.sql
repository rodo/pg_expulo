DROP TABLE IF EXISTS race;
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
  class_code varchar(10)
);


CREATE TABLE skipper (
  id serial primary key,
  name text,
  age int,
  created_at timestamp with time zone default now(),
  updated_at timestamp without time zone,
  price float default 4.5
);

CREATE TABLE race (
  id serial primary key,
  name text,
  year int
);

/* This table is never empty on destination
 *
 */
CREATE TABLE results (
  id serial primary key,
  year int
);

/* This table is not present in configuration so stay empty on destination
 */
CREATE TABLE cheater (name text default 'fool');
