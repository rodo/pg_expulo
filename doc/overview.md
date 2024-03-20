# Overview

## Usage

Design to be simple at use with powerful features, running `pg_expulo` is enough once it's configured.
Some option can change the default behaviour

### the debug option

Set log level to DEBUG

```code
pg_expulo -debug
```

### the generate option

Generate a configuration file based on information fetched from target database. The configuration fill will contains all the tables present and some defaults generator.

```code
pg_expulo -generate -config config.auto.json
```

### the defaults option

Use a file to define your own default values for column generator


```code
pg_expulo -generate -defaults samples/defaults.json
```


### the try option

You may want to test if your config file works or if your action will be successful, the `-try` option is here to do that, all the reads will be done on **source** and write will be done on **target** but with a final **ROLLBACK** at the end to ensure you to not chnage any data. As is you can fine tune your config file in case of needed

```code
pg_expulo -try
```

### the purge option

If you need to only clean the data on target,the option `-purge` is here to do that

```code
pg_expulo -purge
```

Obviously you can combine both options

```code
pg_expulo -try -purge
```

### the version option

Print version number on STDOUT and exit

```code
pg_expulo -version
```
