# Quickstart

Create a user on source database with apropriate privileges. We
**strongly** encourage you to protect your source database to avoid
any disaster.

On **source**

```code sql
CREATE USER expulo;
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM expulo;
REVOKE ALL ON ALL SEQUENCES IN SCHEMA public FROM expulo;
GRANT SELECT on ALL TABLES IN SCHEMA public TO expulo ;
```

On **target** database you must use the **owner** of tables to be able to do all operations needed by pg_expulo.

Define the environment variables to use the new user on
**source** and the owner of tables and sequences on **target**

```code
export PGSRCHOST=localhost
export PGSRCPORT=5436
export PGSRCUSER=expulo
export PGSRCPASSWORD=****
export PGSRCDATABASE=source

export PGDSTHOST=localhost
export PGDSTPORT=5436
export PGDSTUSER=owner
export PGDSTPASSWORD=***
export PGDSTDATABASE=destination
```

Generate your configuration file.

```code
pg_expulo -generate -config config.json
```

Edit the configuration file to suit your needs.

Try the configuration to validate everything is ok

```code
pg_expulo -config config.json -try
```

Now run pg_expulo to copy your data from source to target

```code
pg_expulo -config config.json
```

Enjoy your data on target database.
