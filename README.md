# pg_expulo EXtract, PUrge, LOad


pg_expulo connect to two PostgreSQL instances, typically one in
production and on in staging environement. The extraction is
configured in a json file, data are read on source, data are anonymize
with your own rules, and load on target.

pg_expulo aims to be explicit, easy to use, compliant and
efficient. It reads data only on tables you mentionned in the config
file, but will deal with all columns without the needs to define all
of them. pg_expulo permits you to anonymize all your sensitive
data. It will load them as fast as possible with bulk insert to the
target.

pg_expulo is able to deal with foreign keys and serial values, it will
automatically set the right value.

## [Quickstart](doc/quickstart.md)

## Usage

The two connections strings are build from env variables, you need at least these 10 variables

    export PGSRCHOST=localhost
    export PGSRCPORT=5432
    export PGSRCUSER=rodo
    export PGSRCPASSWORD=*****
    export PGSRCDATABASE=source
    export PGDSTHOST=localhost
    export PGDSTPORT=5432
    export PGDSTUSER=rodo
    export PGDSTPASSWORD=******
    export PGDSTDATABASE=destination

You can run pg_expulo without any parameter

```code
pg_expulo
```

Use a specific config file

```code
pg_expulo --config config/special.json
```

Just give a try to your configuration. Data will be read, inserted on target, but **ROLLBACK** at the end

```code
pg_expulo --config config/special.json --test
```

You can only purge your target with

```code
pg_expulo --config config/special.json --purge
```


## Configuration

All configuration is done on `config.json` file which is read in the
directory.

The main concept of pg_expulo is configuration is explicti at `table`
level but implicit on `columns`, that means pg_expulo will only act on
table that are declared in configuration. On column side pg_expulo
will take all columns of the table in account.
