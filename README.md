# expulo PostgreSQL data tool to extract, purge and load.  The extract
part will fetch data with a set of `SELECT` statement execute on the
**source**, the purge action will anonymise data based on ruls defined
in configuration json file. The load action is a set on `INSERT`
statements run on **destination**

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

You can run pgexpulo without any parameter

```code
pg_expluo
```

## Configuration

All configuration is done on `config.json` file which is read in the
directory (will be change soon)

The main concept of pg_expulo is configuration is explicti at `table`
level but implicit on `columns`, that means pg_expulo will only act on
table that are declared in configuration. On column side pg_expulo
will take all columns of the table in account.